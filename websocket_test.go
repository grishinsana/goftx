package goftx

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/grishinsana/goftx/models"
)

func TestStream_SubscribeToTickers(t *testing.T) {
	ftx := New()

	symbol := "ETH/BTC"

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	data, err := ftx.Stream.SubscribeToTickers(ctx, symbol)
	require.NoError(t, err)

	go func() {
		<-ctx.Done()
		<-time.After(time.Second)
		if _, ok := <-data; !ok {
			t.Fail()
		}
	}()

	count := 0
	for msg := range data {
		require.Equal(t, symbol, msg.Symbol)
		require.Equal(t, models.Update, msg.Type)
		require.True(t, msg.Last.IsPositive())
		require.True(t, msg.Ask.IsPositive())
		require.True(t, msg.Bid.IsPositive())
		require.True(t, msg.AskSize.IsPositive())
		require.True(t, msg.BidSize.IsPositive())
		require.True(t, msg.Bid.LessThanOrEqual(msg.Ask))
		count++
	}
	require.True(t, count > 0)
}

func TestStream_SubscribeToMarkets(t *testing.T) {
	ftx := New()

	symbol := "ETH/BTC"

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	data, err := ftx.Stream.SubscribeToMarkets(ctx)
	require.NoError(t, err)

	go func() {
		<-ctx.Done()
		<-time.After(time.Second)
		if _, ok := <-data; !ok {
			t.Fail()
		}
	}()

	count := 0
	for msg := range data {
		if msg.Name != symbol {
			continue
		}
		require.Equal(t, symbol, msg.Name)
		require.Equal(t, true, msg.Enabled)
		require.Equal(t, "BTC", msg.QuoteCurrency)
		require.Equal(t, "ETH", msg.BaseCurrency)
		count++
	}
	require.True(t, count > 0)
}

func TestStream_SubscribeToTrades(t *testing.T) {
	ftx := New()

	symbol := "BTC-PERP"

	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	data, err := ftx.Stream.SubscribeToTrades(ctx, symbol)
	require.NoError(t, err)

	go func() {
		<-ctx.Done()
		<-time.After(time.Second)
		if _, ok := <-data; !ok {
			t.Fail()
		}
	}()

	lastID := int64(0)
	for msg := range data {
		require.Equal(t, symbol, msg.Symbol)
		require.True(t, msg.Price.IsPositive())
		require.True(t, msg.Size.IsPositive())
		require.True(t, msg.ID > lastID)
		lastID = msg.ID
	}
}

func TestStream_SubscribeToOrderBooks(t *testing.T) {
	ftx := New()

	symbol := "ETH/BTC"

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	data, err := ftx.Stream.SubscribeToOrderBooks(ctx, symbol)
	require.NoError(t, err)

	go func() {
		<-ctx.Done()
		<-time.After(time.Second)
		if _, ok := <-data; !ok {
			t.Fail()
		}
	}()

	count := 0
	for msg := range data {
		require.Equal(t, symbol, msg.Symbol)
		require.True(t, msg.Type == models.Update || msg.Type == models.Partial)
		require.True(t, len(msg.Bids) > 0 || len(msg.Asks) > 0)
		count++
	}
	require.True(t, count > 0)
}
