package goftx

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpotMargin(t *testing.T) {
	_ = godotenv.Load()

	ftx := New(
		WithAuth(os.Getenv("FTX_KEY"), os.Getenv("FTX_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	require.NoError(t, err)

	t.Run("GetBorrowRates", func(t *testing.T) {
		result, err := ftx.SpotMargin.GetBorrowRates()
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("GetLendingRates", func(t *testing.T) {
		result, err := ftx.SpotMargin.GetLendingRates()
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("GetDailyBorrowedAmounts", func(t *testing.T) {
		result, err := ftx.SpotMargin.GetDailyBorrowedAmounts()
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("GetMarketInfo", func(t *testing.T) {
		result, err := ftx.SpotMargin.GetMarketInfo("ETH/BTC")
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("GetBorrowHistory", func(t *testing.T) {
		result, err := ftx.SpotMargin.GetBorrowHistory()
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("GetLendingHistory", func(t *testing.T) {
		result, err := ftx.SpotMargin.GetLendingHistory()
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("GetLendingOffers", func(t *testing.T) {
		result, err := ftx.SpotMargin.GetLendingOffers()
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("GetLendingInfo", func(t *testing.T) {
		result, err := ftx.SpotMargin.GetLendingInfo()
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}
