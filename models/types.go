package models

import (
	"encoding/json"
	"math"
	"time"
)

type Resolution int

const (
	Sec15    = 15
	Minute   = 60
	Minute5  = 300
	Minute15 = 900
	Hour     = 3600
	Hour4    = 14400
	Day      = 86400
)

type Channel string

const (
	OrderBookChannel = Channel("orderbook")
	TradesChannel    = Channel("trades")
	TickerChannel    = Channel("ticker")
	MarketsChannel   = Channel("markets")
)

type Operation string

const (
	Subscribe   = Operation("subscribe")
	UnSubscribe = Operation("unsubscribe")
)

type ResponseType string

const (
	Error        = ResponseType("error")
	Subscribed   = ResponseType("subscribed")
	UnSubscribed = ResponseType("unsubscribed")
	Info         = ResponseType("info")
	Partial      = ResponseType("partial")
	Update       = ResponseType("update")
)

type TransferStatus string

const Complete = TransferStatus("complete")

type OrderType string

const (
	LimitOrder  = OrderType("limit")
	MarketOrder = OrderType("market")
)

type Side string

const (
	Sell = Side("sell")
	Buy  = Side("buy")
)

type Status string

const (
	New    = Status("new")
	Open   = Status("open")
	Closed = Status("closed")
)

type TriggerOrderType string

const (
	Stop         = TriggerOrderType("stop")
	TrailingStop = TriggerOrderType("trailing_stop")
	TakeProfit   = TriggerOrderType("take_profit")
)

type FTXTime struct {
	Time time.Time
}

func (f *FTXTime) UnmarshalJSON(data []byte) error {
	var t float64
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	sec, nsec := math.Modf(t)
	f.Time = time.Unix(int64(sec), int64(nsec))
	return nil
}

func (f FTXTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(f.Time.UnixNano()) / float64(1000000000))
}
