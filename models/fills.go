package models

import (
	"github.com/shopspring/decimal"
)

type GetFillsParams struct {
	Market    *string `json:"market"`
	Limit     *int    `json:"limit"`
	StartTime *int64  `json:"start_time"`
	EndTime   *int64  `json:"end_time"`
	Order     *string `json:"order"`
	OrderID   *int64  `json:"orderId"`
}

type Fill struct {
	Fee           float64         `json:"fee"`
	FeeCurrency   string          `json:"feeCurrency"`
	FeeRate       float64         `json:"feeRate"`
	Future        string          `json:"future"`
	ID            int64           `json:"id"`
	Liquidity     Liquidity       `json:"liquidity"`
	Market        string          `json:"market"`
	BaseCurrency  string          `json:"baseCurrency"`
	QuoteCurrency string          `json:"quoteCurrency"`
	OrderID       int64           `json:"orderId"`
	TradeID       int64           `json:"tradeId"`
	Price         decimal.Decimal `json:"price"`
	Side          Side            `json:"side"`
	Size          decimal.Decimal `json:"size"`
	Time          FTXTime         `json:"time"`
	Type          string          `json:"type"`
}
