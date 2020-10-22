package goftx

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/grishinsana/goftx/models"
)

const (
	wsUrl = "wss://ftx.com/ws/"

	websocketTimeout  = time.Second * 60
	pingPeriod        = (websocketTimeout * 9) / 10
	reconnectCount    = int(10)
	reconnectInterval = time.Second
)

type Connection struct {
	conn     *websocket.Conn
	requests []models.WSRequest
}

type Stream struct {
	client                 *Client
	wsReconnectionCount    int
	wsReconnectionInterval time.Duration
}

func (s *Stream) SetReconnectionCount(count int) {
	s.wsReconnectionCount = count
}

func (s *Stream) SetReconnectionInterval(interval time.Duration) {
	s.wsReconnectionInterval = interval
}

func (s *Stream) serve(ctx context.Context, requests ...models.WSRequest) (chan []byte, error) {
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	log.Printf("connected to %v", wsUrl)

	err = subscribe(conn, requests)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	conn.SetPongHandler(func(msg string) error {
		log.Printf("PONG")
		return nil
	})

	doneC := make(chan struct{})
	reconnectC := make(chan struct{})
	eventsC := make(chan []byte, 1)

	go func() {
		defer close(reconnectC)

		go func() {
			defer close(doneC)
			defer conn.Close()

			for {
				message := &models.WsResponse{}
				err = conn.ReadJSON(&message)
				if err != nil {
					log.Printf("read msg: %v\n", err)
					if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
						return
					}
					reconnectC <- struct{}{}
				}

				log.Println("message: ", message.Market, message.Channel, message.Type, message.Code, message.Message)
				log.Printf("data: %s \n", string(message.Data))

				switch message.Type {
				case models.Subscribed:
					continue
				case models.UnSubscribed:
					err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
					if err != nil {
						log.Printf("write close msg: %v\n", err)
						return
					}
				}

				eventsC <- message.Data
			}
		}()

		for {
			select {
			case <-ctx.Done():
				err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Printf("write close msg: %v\n", err)
					return
				}
				select {
				case <-doneC:
					return
				case <-time.After(time.Second):
					return
				}
			case <-doneC:
				return
			case <-reconnectC:
				for i := 0; i < s.wsReconnectionCount; i++ {
					select {
					case <-time.After(s.wsReconnectionInterval):
						conn, _, err = websocket.DefaultDialer.Dial(wsUrl, nil)
						if err == nil {
							log.Printf("reconnected after %v times", i)
							return
						}
					case <-ctx.Done():
						return
					}
				}
			case <-time.After(pingPeriod):
				log.Printf("PING")
				err := conn.WriteControl(websocket.PingMessage, []byte(`{"op": "pong"}`), time.Now().Add(10*time.Second))
				if err != nil && err != websocket.ErrCloseSent {
					log.Printf("write ping: %v", err)
				}
			}
		}
	}()

	return eventsC, nil
}

func subscribe(conn *websocket.Conn, requests []models.WSRequest) error {
	for _, req := range requests {
		err := conn.WriteJSON(req)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (s *Stream) SubscribeToTickers(ctx context.Context, symbols ...string) (chan *models.Ticker, error) {
	if len(symbols) == 0 {
		return nil, errors.New("symbols is missing")
	}

	requests := make([]models.WSRequest, 0, len(symbols))
	for _, symbol := range symbols {
		requests = append(requests, models.WSRequest{
			Channel: models.TickerChannel,
			Market:  symbol,
			Op:      models.Subscribe,
		})
	}

	eventsC, err := s.serve(ctx, requests...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	tickersC := make(chan *models.Ticker, 1)
	go func() {
		defer close(tickersC)
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventsC:
				ticker := &models.Ticker{}
				err = json.Unmarshal(event, &ticker)
				if err != nil {
					log.Printf("unmarshal: %+v\n", err)
					continue
				}
				tickersC <- ticker
			}
		}
	}()

	return tickersC, nil
}

func (s *Stream) SubscribeToMarkets(ctx context.Context) (chan *models.Market, error) {
	eventsC, err := s.serve(ctx, models.WSRequest{
		Channel: models.MarketsChannel,
		Op:      models.Subscribe,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	marketsC := make(chan *models.Market, 1)
	go func() {
		defer close(marketsC)
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventsC:
				var markets struct {
					Data map[string]*models.Market `json:"data"`
				}
				err = json.Unmarshal(event, &markets)
				if err != nil {
					log.Printf("unmarshal: %+v\n", err)
					continue
				}
				for _, market := range markets.Data {
					marketsC <- market
				}
			}
		}
	}()

	return marketsC, nil
}

func (s *Stream) SubscribeToTrades(ctx context.Context, symbols ...string) (chan *models.Trade, error) {
	if len(symbols) == 0 {
		return nil, errors.New("symbols is missing")
	}

	requests := make([]models.WSRequest, 0, len(symbols))
	for _, symbol := range symbols {
		requests = append(requests, models.WSRequest{
			Channel: models.TradesChannel,
			Market:  symbol,
			Op:      models.Subscribe,
		})
	}

	eventsC, err := s.serve(ctx, requests...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	tradesC := make(chan *models.Trade, 1)
	go func() {
		defer close(tradesC)
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventsC:
				trade := &models.Trade{}
				err = json.Unmarshal(event, &trade)
				if err != nil {
					log.Printf("unmarshal: %+v\n", err)
					continue
				}
				tradesC <- trade
			}
		}
	}()

	return tradesC, nil
}

func (s *Stream) SubscribeToOrderBooks(ctx context.Context, symbols ...string) (chan *models.OrderBook, error) {
	if len(symbols) == 0 {
		return nil, errors.New("symbols is missing")
	}

	requests := make([]models.WSRequest, 0, len(symbols))
	for _, symbol := range symbols {
		requests = append(requests, models.WSRequest{
			Channel: models.OrderBookChannel,
			Market:  symbol,
			Op:      models.Subscribe,
		})
	}

	eventsC, err := s.serve(ctx, requests...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	booksC := make(chan *models.OrderBook, 1)
	go func() {
		defer close(booksC)
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventsC:
				orderbook := &models.OrderBook{}
				err = json.Unmarshal(event, &orderbook)
				if err != nil {
					log.Printf("unmarshal: %+v\n", err)
					continue
				}
				booksC <- orderbook
			}
		}
	}()

	return booksC, nil
}
