package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grishinsana/goftx"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	client := goftx.New()

	// subscribeToTickers(ctx, client)

	// subscribeToMarkets(ctx, client)

	// subscribeToTrades(ctx, client)

	subscribeToOrderBooks(ctx, client)

	<-sigs
	cancel()
	time.Sleep(time.Second)
}

func subscribeToTickers(ctx context.Context, client *goftx.Client) {
	data, err := client.Stream.SubscribeToTickers(ctx, "ETH/BTC")
	if err != nil {
		log.Fatalf("%+v", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-data:
				if !ok {
					return
				}
				log.Printf("%+v\n", msg)
			}
		}
	}()
}

func subscribeToMarkets(ctx context.Context, client *goftx.Client) {
	data, err := client.Stream.SubscribeToMarkets(ctx)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-data:
				if !ok {
					return
				}
				log.Printf("%+v\n", msg)
			}
		}
	}()
}

func subscribeToTrades(ctx context.Context, client *goftx.Client) {
	data, err := client.Stream.SubscribeToTrades(ctx, "BTC/USDT")
	if err != nil {
		log.Fatalf("%+v", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-data:
				if !ok {
					return
				}
				log.Printf("%+v\n", msg)
			}
		}
	}()
}

func subscribeToOrderBooks(ctx context.Context, client *goftx.Client) {
	data, err := client.Stream.SubscribeToOrderBooks(ctx, "BTC/USDT")
	if err != nil {
		log.Fatalf("%+v", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-data:
				if !ok {
					return
				}
				log.Printf("%+v\n", msg)
			}
		}
	}()
}
