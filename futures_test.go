package goftx

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grishinsana/goftx/models"
)

func TestFutures_GetFutures(t *testing.T) {
	ftx := New()

	futures, err := ftx.Futures.GetFutures()
	assert.NoError(t, err)
	assert.NotNil(t, futures)
}

func TestFutures_GetFuture(t *testing.T) {
	ftx := New()

	req := require.New(t)

	t.Run("success", func(t *testing.T) {
		expected := &models.Future{
			Name:    "BTC-PERP",
			Enabled: true,
		}

		futures, err := ftx.Futures.GetFuture(expected.Name)
		req.NoError(err)
		req.NotNil(futures)
		req.Equal(expected.Name, futures.Name)
		req.Equal(expected.Enabled, futures.Enabled)
	})

	t.Run("not_found", func(t *testing.T) {
		expected := &models.Future{
			Name:    "incorrect",
			Enabled: true,
		}

		futures, err := ftx.Markets.GetMarketByName(expected.Name)
		req.Error(err)
		req.Nil(futures)
	})
}

func TestFutures_GetFutureStats(t *testing.T) {
	ftx := New()

	futures, err := ftx.Futures.GetFutureStats("BTC-PERP")
	assert.NoError(t, err)
	assert.NotNil(t, futures)
}

func TestFutures_GetFundingRates(t *testing.T) {
	ftx := New()

	req := require.New(t)

	t.Run("success", func(t *testing.T) {
		rate, err := ftx.Futures.GetFundingRates(nil)
		req.NoError(err)
		req.NotNil(rate)
	})

	t.Run("success_with_market", func(t *testing.T) {
		future := "BTC-PERP"
		rates, err := ftx.Futures.GetFundingRates(&models.GetFundingRatesParams{
			Future: &future,
		})
		req.NoError(err)
		req.NotNil(rates)

		for _, rate := range rates {
			req.Equal(rate.Future, future)
		}
	})

	t.Run("success_with_params", func(t *testing.T) {
		rates, err := ftx.Futures.GetFundingRates(&models.GetFundingRatesParams{
			StartTime: PtrInt(int(time.Now().Add(-5 * time.Hour).Unix())),
			EndTime:   PtrInt(int(time.Now().Unix())),
		})
		req.NoError(err)
		req.NotNil(rates)
		req.GreaterOrEqual(rates[0].Time.Unix(), time.Now().Add(-5*time.Hour).Unix())
	})
}
