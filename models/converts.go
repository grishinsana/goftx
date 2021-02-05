package models

import "github.com/shopspring/decimal"

type QuoteStatus struct {
	BaseCoin  string          `json:"baseCoin"`
	Cost      decimal.Decimal `json:"cost"`
	Expired   bool            `json:"expired"`
	Filled    bool            `json:"filled"`
	FromCoin  string          `json:"fromCoin"`
	ID        int64           `json:"id"`
	Price     decimal.Decimal `json:"price"`
	Proceeds  decimal.Decimal `json:"proceeds"`
	QuoteCoin string          `json:"quoteCoin"`
	Side      Side            `json:"side"`
	ToCoin    string          `json:"toCoin"`
}

type CreateQuotePayload struct {
	FromCoin string          `json:"fromCoin"`
	ToCoin   string          `json:"toCoin"`
	Size     decimal.Decimal `json:"size"`
}
