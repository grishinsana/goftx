package models

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type BaseResponse struct {
	Type   ResponseType
	Symbol string
}

type TickerResponse struct {
	Ticker
	BaseResponse
}

type TradesResponse struct {
	Trades []Trade
	BaseResponse
}

type TradeResponse struct {
	Trade
	BaseResponse
}

type OrderBookResponse struct {
	OrderBook
	BaseResponse
}

type WSRequest struct {
	Channel Channel   `json:"channel"`
	Market  string    `json:"market"`
	Op      Operation `json:"op"`
}

type WsResponse struct {
	Channel Channel         `json:"channel"`
	Market  string          `json:"market"`
	Type    ResponseType    `json:"type"`
	Code    int             `json:"code"`
	Message string          `json:"msg"`
	Data    json.RawMessage `json:"data"`
}

func (wr *WsResponse) MapToTradesResponse() (*TradesResponse, error) {
	var trades []Trade
	err := json.Unmarshal(wr.Data, &trades)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &TradesResponse{
		Trades: trades,
		BaseResponse: BaseResponse{
			Type:   wr.Type,
			Symbol: wr.Market,
		},
	}, nil
}

func (wr *WsResponse) MapToTickerResponse() (*TickerResponse, error) {
	ticker := Ticker{}
	err := json.Unmarshal(wr.Data, &ticker)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &TickerResponse{
		Ticker: ticker,
		BaseResponse: BaseResponse{
			Type:   wr.Type,
			Symbol: wr.Market,
		},
	}, nil
}

func (wr *WsResponse) MapToOrderBookResponse() (*OrderBookResponse, error) {
	book := OrderBook{}
	err := json.Unmarshal(wr.Data, &book)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &OrderBookResponse{
		OrderBook: book,
		BaseResponse: BaseResponse{
			Type:   wr.Type,
			Symbol: wr.Market,
		},
	}, nil
}
