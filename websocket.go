package goftx

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/grishinsana/goftx/models"
)

const (
	wsUrlFormat = "wss://ftx.%s/ws/"

	writeWait         = time.Second * 10
	reconnectCount    = int(10)
	reconnectInterval = time.Second
	streamTimeout     = time.Second * 60
)

type Stream struct {
	apiKey                 string
	secret                 string
	subAccount             string
	mu                     *sync.Mutex
	url                    string
	dialer                 *websocket.Dialer
	wsReconnectionCount    int
	wsReconnectionInterval time.Duration
	wsTimeout              time.Duration
	isDebugMode            bool
	serverTimeDiff         time.Duration
	authorized             bool
}

func (s *Stream) SetStreamTimeout(timeout time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.wsTimeout = timeout
}

func (s *Stream) SetReconnectionCount(count int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.wsReconnectionCount = count
}

func (s *Stream) SetDebugMode(isDebugMode bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.isDebugMode = isDebugMode
}

func (s *Stream) SetReconnectionInterval(interval time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.wsReconnectionInterval = interval
}

func (s *Stream) printf(format string, v ...interface{}) {
	if !s.isDebugMode {
		return
	}
	if len(v) > 0 {
		log.Printf(format+"\n", v)
	} else {
		log.Printf(format + "\n")
	}
}

func (s *Stream) connect(requests ...models.WSRequest) (*websocket.Conn, error) {
	conn, _, err := s.dialer.Dial(s.url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	s.printf("connected to %v", s.url)

	err = s.subscribe(conn, requests)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	conn.SetPongHandler(func(msg string) error {
		s.printf("PONG")
		_ = conn.SetReadDeadline(time.Now().Add(s.wsTimeout))
		return nil
	})

	return conn, nil
}

func (s *Stream) serve(ctx context.Context, requests ...models.WSRequest) (chan interface{}, error) {
	conn, err := s.connect(requests...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	doneC := make(chan struct{})
	eventsC := make(chan interface{}, 1)

	go func() {
		go func() {
			defer close(doneC)

			for {
				message := &models.WsResponse{}
				err = conn.ReadJSON(&message)
				if err != nil {
					s.printf("read msg: %v", err)
					if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
						return
					}
					conn, err = s.reconnect(ctx, requests)
					if err != nil {
						s.printf("reconnect: %+v", err)
						return
					}
					continue
				}

				switch message.Type {
				case models.Subscribed, models.UnSubscribed:
					continue
				}

				var response interface{}
				switch message.Channel {
				case models.TickerChannel:
					response, err = message.MapToTickerResponse()
				case models.TradesChannel:
					response, err = message.MapToTradesResponse()
				case models.OrderBookChannel:
					response, err = message.MapToOrderBookResponse()
				case models.OrdersChannel:
					response, err = message.MapToOrderResponse()
				case models.FillsChannel:
					response, err = message.MapToFillResponse()
				case models.MarketsChannel:
					response = message.Data
				}

				eventsC <- response
			}
		}()

		for {
			select {
			case <-ctx.Done():
				err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					s.printf("write close msg: %v", err)
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
			case <-time.After((s.wsTimeout * 9) / 10):
				s.printf("PING")
				if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeWait)); err != nil {
					s.printf("write ping: %v", err)
				}
			}
		}
	}()

	return eventsC, nil
}

// Credit to https://github.com/go-numb/go-ftx
// nolint:errcheck
func (s *Stream) auth(conn *websocket.Conn) error {
	if s.apiKey == "" {
		return errors.New("credentials is required")
	}

	s.printf("Authenticate websocket connection")
	msec := time.Now().UTC().Add(s.serverTimeDiff).UnixNano() / int64(time.Millisecond)

	mac := hmac.New(sha256.New, []byte(s.secret))
	mac.Write([]byte(fmt.Sprintf("%dwebsocket_login", msec)))
	args := map[string]interface{}{
		"key":  s.apiKey,
		"sign": hex.EncodeToString(mac.Sum(nil)),
		"time": msec,
	}
	if s.subAccount != "" {
		args["subaccount"] = s.subAccount
	}

	return conn.WriteJSON(models.WSRequest{
		Op:   models.Login,
		Args: args,
	})
}

func (s *Stream) reconnect(ctx context.Context, requests []models.WSRequest) (*websocket.Conn, error) {
	started := time.Now()

	for i := 1; i < s.wsReconnectionCount; i++ {
		conn, err := s.connect(requests...)
		if err == nil {
			return conn, nil
		}

		timeout := time.Duration(int64(math.Pow(2, float64(i)))) * time.Second

		select {
		case <-time.After(timeout):
			conn, err := s.connect(requests...)
			if err != nil {
				continue
			}

			return conn, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return nil, errors.Errorf("reconnection failed after %v", time.Since(started))
}

func (s *Stream) subscribe(conn *websocket.Conn, requests []models.WSRequest) error {
	for _, req := range requests {
		if req.IsPrivateChannel() && !s.authorized {
			err := s.auth(conn)
			if err != nil {
				return errors.WithStack(err)
			}
			s.authorized = true
		}

		err := conn.WriteJSON(req)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (s *Stream) SubscribeToFills(ctx context.Context) (chan *models.FillResponse, error) {
	eventsC, err := s.serve(ctx, models.WSRequest{
		Channel: models.FillsChannel,
		Op:      models.Subscribe,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	fillsC := make(chan *models.FillResponse, 1)
	go func() {
		defer close(fillsC)
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-eventsC:
				if !ok {
					return
				}
				fill, ok := event.(*models.FillResponse)
				if !ok {
					return
				}
				fillsC <- fill
			}
		}
	}()

	return fillsC, nil
}

func (s *Stream) SubscribeToOrders(ctx context.Context) (chan *models.OrderResponse, error) {
	eventsC, err := s.serve(ctx, models.WSRequest{
		Channel: models.OrdersChannel,
		Op:      models.Subscribe,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ordersC := make(chan *models.OrderResponse, 1)
	go func() {
		defer close(ordersC)
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-eventsC:
				if !ok {
					return
				}
				order, ok := event.(*models.OrderResponse)
				if !ok {
					return
				}
				ordersC <- order
			}
		}
	}()

	return ordersC, nil
}

// nolint: dupl
func (s *Stream) SubscribeToTickers(ctx context.Context, symbols ...string) (chan *models.TickerResponse, error) {
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

	tickersC := make(chan *models.TickerResponse, 1)
	go func() {
		defer close(tickersC)
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-eventsC:
				if !ok {
					return
				}
				ticker, ok := event.(*models.TickerResponse)
				if !ok {
					return
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
			case event, ok := <-eventsC:
				if !ok {
					return
				}
				data, ok := event.(json.RawMessage)
				if !ok {
					return
				}
				var markets struct {
					Data map[string]*models.Market `json:"data"`
				}
				err = json.Unmarshal(data, &markets)
				if err != nil {
					s.printf("unmarshal markets: %+v", err)
					return
				}
				for _, market := range markets.Data {
					marketsC <- market
				}
			}
		}
	}()

	return marketsC, nil
}

func (s *Stream) SubscribeToTrades(ctx context.Context, symbols ...string) (chan *models.TradeResponse, error) {
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

	tradesC := make(chan *models.TradeResponse, 1)
	go func() {
		defer close(tradesC)
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-eventsC:
				if !ok {
					return
				}
				trades, ok := event.(*models.TradesResponse)
				if !ok {
					return
				}
				for _, trade := range trades.Trades {
					tradesC <- &models.TradeResponse{
						Trade:        trade,
						BaseResponse: trades.BaseResponse,
					}
				}
			}
		}
	}()

	return tradesC, nil
}

// nolint: dupl
func (s *Stream) SubscribeToOrderBooks(ctx context.Context, symbols ...string) (chan *models.OrderBookResponse, error) {
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

	booksC := make(chan *models.OrderBookResponse, 1)
	go func() {
		defer close(booksC)
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-eventsC:
				if !ok {
					return
				}
				book, ok := event.(*models.OrderBookResponse)
				if !ok {
					return
				}
				booksC <- book
			}
		}
	}()

	return booksC, nil
}
