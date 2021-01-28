package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Market struct {
	Name           string          `json:"name"`
	Underlying     string          `json:"underlying"`
	BaseCurrency   string          `json:"baseCurrency"`
	QuoteCurrency  string          `json:"quoteCurrency"`
	Enabled        bool            `json:"enabled"`
	Ask            decimal.Decimal `json:"ask"`
	Bid            decimal.Decimal `json:"bid"`
	Last           decimal.Decimal `json:"last"`
	PostOnly       bool            `json:"postOnly"`
	PriceIncrement decimal.Decimal `json:"priceIncrement"`
	SizeIncrement  decimal.Decimal `json:"sizeIncrement"`
	Restricted     bool            `json:"restricted"`
}

// The bids and asks are formatted like so:
// [[best price, size at price], [next next best price, size at price], ...]
//
// Checksum
// Every message contains a signed 32-bit integer checksum of the orderbook.
// You can run the same checksum on your client orderbook state and compare it to checksum field.
// If they are the same, your client's state is correct.
// If not, you have likely lost or mishandled a packet and should re-subscribe to receive the initial snapshot.
//
// The checksum operates on a string that represents the first 100 orders on the orderbook on either side. The format of the string is:
//
// <best_bid_price>:<best_bid_size>:<best_ask_price>:<best_ask_size>:<second_best_bid_price>:<second_best_ask_price>:...
// For example, if the orderbook was comprised of the following two bids and asks:
//
// bids: [[5000.5, 10], [4995.0, 5]]
// asks: [[5001.0, 6], [5002.0, 7]]
// The string would be '5005.5:10:5001.0:6:4995.0:5:5002.0:7'
//
// If there are more orders on one side of the book than the other, then simply omit the information about orders that don't exist.
//
// For example, if the orderbook had the following bids and asks:
//
// bids: [[5000.5, 10], [4995.0, 5]]
// asks: [[5001.0, 6]]
// The string would be '5005.5:10:5001.0:6:4995.0:5'
//
// The final checksum is the crc32 value of this string.
type OrderBook struct {
	Asks     [][]decimal.Decimal `json:"asks"`
	Bids     [][]decimal.Decimal `json:"bids"`
	Checksum int64               `json:"checksum,omitempty"`
	Time     FTXTime             `json:"time"`
}

type Trade struct {
	ID          int64           `json:"id"`
	Liquidation bool            `json:"liquidation"`
	Price       decimal.Decimal `json:"price"`
	Side        string          `json:"side"`
	Size        decimal.Decimal `json:"size"`
	Time        time.Time       `json:"time"`
}

type HistoricalPrice struct {
	StartTime time.Time       `json:"startTime"`
	Open      decimal.Decimal `json:"open"`
	Close     decimal.Decimal `json:"close"`
	High      decimal.Decimal `json:"high"`
	Low       decimal.Decimal `json:"low"`
	Volume    decimal.Decimal `json:"volume"`
}

type Ticker struct {
	Bid     decimal.Decimal `json:"bid"`
	Ask     decimal.Decimal `json:"ask"`
	BidSize decimal.Decimal `json:"bidSize"`
	AskSize decimal.Decimal `json:"askSize"`
	Last    decimal.Decimal `json:"last"`
	Time    FTXTime         `json:"time"`
}

type Fill struct {
	Fee       decimal.Decimal `json:"fee"`
	FeeRate   decimal.Decimal `json:"feeRate"`
	Future    string          `json:"future"`
	ID        int64           `json:"id"`
	Liquidity string          `json:"liquidity"`
	Market    string          `json:"market"`
	OrderID   int64           `json:"orderId"`
	TradeID   int64           `json:"tradeId"`
	Price     decimal.Decimal `json:"price"`
	Side      string          `json:"side"`
	Size      decimal.Decimal `json:"size"`
	Time      FTXTime         `json:"time"`
	Type      string          `json:"type"`
}

type GetTradesParams struct {
	Limit     *int `json:"limit"`
	StartTime *int `json:"start_time"`
	EndTime   *int `json:"end_time"`
}

type GetHistoricalPricesParams struct {
	Resolution Resolution `json:"resolution"`
	Limit      *int       `json:"limit"`
	StartTime  *int       `json:"start_time"`
	EndTime    *int       `json:"end_time"`
}
