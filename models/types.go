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
	FillsChannel     = Channel("fills")
	OrdersChannel    = Channel("orders")
)

type Operation string

const (
	Subscribe   = Operation("subscribe")
	UnSubscribe = Operation("unsubscribe")
	Login       = Operation("login")
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
	TrailingStop = TriggerOrderType("trailingStop")
	TakeProfit   = TriggerOrderType("takeProfit")
)

type FTXTime struct {
	Time time.Time
}

func (f *FTXTime) UnmarshalJSON(data []byte) error {
	var t float64
	err := json.Unmarshal(data, &t)

	// FTX uses ISO format sometimes so we have to detect and handle that differently.
	if err != nil {
		var iso time.Time
		errIso := json.Unmarshal(data, &iso)

		if errIso != nil {
			return err
		}

		f.Time = iso
		return nil
	}

	sec, nsec := math.Modf(t)
	f.Time = time.Unix(int64(sec), int64(nsec))
	return nil
}

func (f FTXTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(f.Time.UnixNano()) / float64(1000000000))
}
