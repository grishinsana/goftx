package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID            int64           `json:"id"`
	Market        string          `json:"market"`
	Type          OrderType       `json:"type"`
	Side          Side            `json:"side"`
	Price         decimal.Decimal `json:"price"`
	Size          decimal.Decimal `json:"size"`
	FilledSize    decimal.Decimal `json:"filledSize"`
	RemainingSize decimal.Decimal `json:"remainingSize"`
	AvgFillPrice  decimal.Decimal `json:"avgFillPrice"`
	Status        Status          `json:"status"`
	CreatedAt     time.Time       `json:"createdAt"`
	ReduceOnly    bool            `json:"reduceOnly"`
	Ioc           bool            `json:"ioc"`
	PostOnly      bool            `json:"postOnly"`
	Future        string          `json:"future"`
	ClientID      string          `json:"clientId"`
}

type GetOrdersHistoryParams struct {
	Market    *string `json:"market"`
	Limit     *int    `json:"limit"`
	StartTime *int    `json:"start_time"`
	EndTime   *int    `json:"end_time"`
}

type TriggerOrder struct {
	ID               int64            `json:"id"`
	OrderID          int64            `json:"orderId"`
	Market           string           `json:"market"`
	CreatedAt        time.Time        `json:"createdAt"`
	Error            string           `json:"error"`
	Future           string           `json:"future"`
	OrderPrice       decimal.Decimal  `json:"orderPrice"`
	ReduceOnly       bool             `json:"reduceOnly"`
	Side             Side             `json:"side"`
	Size             decimal.Decimal  `json:"size"`
	Status           Status           `json:"status"`
	TrailStart       decimal.Decimal  `json:"trailStart"`
	TrailValue       decimal.Decimal  `json:"trailValue"`
	TriggerPrice     decimal.Decimal  `json:"triggerPrice"`
	TriggeredAt      time.Time        `json:"triggeredAt"`
	Type             TriggerOrderType `json:"type"`
	OrderType        OrderType        `json:"orderType"`
	FilledSize       decimal.Decimal  `json:"filledSize"`
	AvgFillPrice     decimal.Decimal  `json:"avgFillPrice"`
	OrderStatus      string           `json:"orderStatus"`
	RetryUntilFilled bool             `json:"retryUntilFilled"`
}

type GetOpenTriggerOrdersParams struct {
	Market *string           `json:"market"`
	Type   *TriggerOrderType `json:"type"`
}

type Trigger struct {
	Error      string    `json:"error"`
	FilledSize float64   `json:"filledSize"`
	OrderSize  float64   `json:"orderSize"`
	OrderID    int64     `json:"orderId"`
	Time       time.Time `json:"time"`
}
