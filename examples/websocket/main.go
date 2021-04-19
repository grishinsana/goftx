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

	client := goftx.New(
		goftx.WithAuth(os.Getenv("FTX_KEY"), os.Getenv("FTX_SECRET")),
		goftx.WithFTXUS(),
	)

	client.Stream.SetStreamTimeout(60 * time.Second)
	client.Stream.SetDebugMode(true)

	subscribeToTickers(ctx, client)

	// subscribeToMarkets(ctx, client)

	//subscribeToFills(ctx, client)

	//subscribeToOrders(ctx, client)

	subscribeToTrades(ctx, client)

	// subscribeToOrderBooks(ctx, client)

	<-sigs
	cancel()
	time.Sleep(2 * time.Second)
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
	data, err := client.Stream.SubscribeToTrades(ctx, "BTC-PERP")
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

func subscribeToFills(ctx context.Context, client *goftx.Client) {
	data, err := client.Stream.SubscribeToFills(ctx)
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

func subscribeToOrders(ctx context.Context, client *goftx.Client) {
	data, err := client.Stream.SubscribeToOrders(ctx)
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
				log.Printf("------ %+v\n", msg)
			}
		}
	}()
}
