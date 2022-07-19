package goftx

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"github.com/grishinsana/goftx/models"
)

func TestClient_GetQuote_GetQuote_GetQuotes(t *testing.T) {
	_ = godotenv.Load()

	ftx := New(
		WithAuth(os.Getenv("FTX_KEY"), os.Getenv("FTX_SECRET")),
	)

	quoteId, err := ftx.CreateQuote(
		&models.CreateQuotePayload{
			FromCoin: "SOL",
			ToCoin:   "USD",
			Size:     decimal.RequireFromString("0.03"),
		},
	)
	require.NoError(t, err)
	require.NotZero(t, quoteId)
	quoteStatus, err := ftx.GetQuote(quoteId)
	require.NoError(t, err)
	expected := &models.QuoteStatus{
		BaseCoin:  "SOL",
		Expired:   false,
		Filled:    false,
		FromCoin:  "SOL",
		QuoteCoin: "USD",
		Side:      "sell",
		ToCoin:    "USD",
	}
	require.Empty(
		t,
		cmp.Diff(
			expected,
			quoteStatus,
			cmpopts.IgnoreFields(
				models.QuoteStatus{},
				"Cost",
				"ID",
				"Price",
				"Proceeds",
			),
		),
	)
}
