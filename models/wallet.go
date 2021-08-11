package models

import "time"

type CreateWithdrawPayload struct {
	Coin    string  `json:"coin"`
	Size    float64 `json:"size"`
	Address string  `json:"address"`
	Tag     string  `json:"tag,omitempty"`
	Method  string  `json:"method,omitempty"`
}

type CreateWithdrawResult struct {
	ID      int64     `json:"id"`
	Coin    string    `json:"coin"`
	Address string    `json:"address"`
	Tag     string    `json:"tag"`
	Fee     float64   `json:"fee"`
	Size    float64   `json:"size"`
	Status  string    `json:"status"`
	TxID    string    `json:"txid"`
	Time    time.Time `json:"time"`
}
