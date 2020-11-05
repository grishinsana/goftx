# goftx
FTX exchange golang library

### Install
```shell script
go get github.com/grishinsana/goftx
```

### Usage

> See examples directory and test case for more examples

#### REST
```go
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/grishinsana/goftx"
)

func main() {
	client := goftx.New(
		goftx.WithAuth("API-KEY", "API-SECRET"),
		goftx.WithHTTPClient(&http.Client{
			Timeout: 5 * time.Second,
		}),
	)

	info, err := client.Account.GetAccountInformation()
	if err != nil {
		panic(err)
	}
	fmt.Println(info)
}
```

#### WebSocket
```go
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

	<-sigs
	cancel()
	time.Sleep(time.Second)
}
```

### No Logged In Error
"Not logged in" errors usually come from a wrong signatures. FTX released an article on how to authenticate https://blog.ftx.com/blog/api-authentication/

If you have unauthorized error to private methods, then you need to use SetServerTimeDiff()
```go
ftx := New()
ftx.SetServerTimeDiff()
```
