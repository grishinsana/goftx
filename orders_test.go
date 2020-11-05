package goftx

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grishinsana/goftx/models"
)

func TestOrders_GetOpenOrders(t *testing.T) {
	godotenv.Load()

	ftx := New(
		WithAuth(os.Getenv("FTX_KEY"), os.Getenv("FTX_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	require.NoError(t, err)

	market := "ETH/BTC"

	orders, err := ftx.Orders.GetOpenOrders(market)
	assert.NoError(t, err)
	assert.NotNil(t, orders)
}

func TestOrders_GetOrdersHistory(t *testing.T) {
	godotenv.Load()

	ftx := New(
		WithAuth(os.Getenv("FTX_KEY"), os.Getenv("FTX_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	require.NoError(t, err)

	market := "ETH/BTC"

	orders, err := ftx.Orders.GetOrdersHistory(&models.GetOrdersHistoryParams{
		Market:    &market,
		Limit:     nil,
		StartTime: nil,
		EndTime:   nil,
	})
	assert.NoError(t, err)
	assert.NotNil(t, orders)
}

func TestOrders_GetOpenTriggerOrders(t *testing.T) {
	godotenv.Load()

	ftx := New(
		WithAuth(os.Getenv("FTX_KEY"), os.Getenv("FTX_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	require.NoError(t, err)

	market := "ETH/BTC"
	orderType := models.Stop

	orders, err := ftx.Orders.GetOpenTriggerOrders(&models.GetOpenTriggerOrdersParams{
		Market: &market,
		Type:   &orderType,
	})
	assert.NoError(t, err)
	assert.NotNil(t, orders)
}

func TestOrders_GetOrderTriggers(t *testing.T) {
	godotenv.Load()

	ftx := New(
		WithAuth(os.Getenv("FTX_KEY"), os.Getenv("FTX_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	require.NoError(t, err)

	orderID := int64(1111)

	triggers, err := ftx.Orders.GetOrderTriggers(orderID)

	// 400 - Bad Request, orderID doesn't exist
	assert.Error(t, err)
	assert.Nil(t, triggers)
}
