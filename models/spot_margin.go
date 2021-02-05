package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type BorrowRate struct {
	Coin     string  `json:"coin"`
	Estimate float64 `json:"estimate"`
	Previous float64 `json:"previous"`
}

type LendingRate BorrowRate

type BorrowSummary struct {
	Coin string          `json:"coin"`
	Size decimal.Decimal `json:"size"`
}

type SpotMarginMarketInfo struct {
	Coin         string          `json:"coin"`
	Borrowed     decimal.Decimal `json:"borrowed"`
	Free         decimal.Decimal `json:"free"`
	EsimatedRate float64         `json:"estimatedRate"`
	PreviousRate float64         `json:"previousRate"`
}

type GetSpotMarginMarketInfoResponse struct {
	Base  SpotMarginMarketInfo
	Quote SpotMarginMarketInfo
}

type BorrowHistory struct {
	Coin string          `json:"coin"`
	Cost decimal.Decimal `json:"cost"`
	Rate decimal.Decimal `json:"rate"`
	Size decimal.Decimal `json:"size"`
	Time time.Time       `json:"time"`
}

type LendingHistory BorrowHistory

type LendingOffer struct {
	Coin string          `json:"coin"`
	Rate float64         `json:"rate"`
	Size decimal.Decimal `json:"size"`
}

type LendingInfo struct {
	Coin     string          `json:"coin"`
	Lendable decimal.Decimal `json:"lendable"`
	Locked   decimal.Decimal `json:"locked"`
	MinRate  float64         `json:"minRate"`
	Offered  decimal.Decimal `json:"offered"`
}

type LendingOfferPayload struct {
	Coin string          `json:"coin"`
	Size decimal.Decimal `json:"size"`
	Rate float64         `json:"rate"`
}
