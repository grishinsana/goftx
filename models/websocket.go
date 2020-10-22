package models

import "encoding/json"

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

type WSRequest struct {
	Channel Channel   `json:"channel"`
	Market  string    `json:"market"`
	Op      Operation `json:"op"`
}

type ResponseType string

const (
	Error        = ResponseType("error")
	Subscribed   = ResponseType("subscribed")
	UnSubscribed = ResponseType("unsubscribed")
	Info         = ResponseType("info")
	Partial      = ResponseType("partial")
	Update       = ResponseType("update")
)

type WsResponse struct {
	Channel Channel         `json:"channel"`
	Market  string          `json:"market"`
	Type    ResponseType    `json:"type"`
	Code    int             `json:"code"`
	Message string          `json:"msg"`
	Data    json.RawMessage `json:"data"`
}
