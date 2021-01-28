package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/grishinsana/goftx"
)

func main() {
	client := goftx.New(
		goftx.WithAuth("API-KEY", "API-SECRET", "API-SUBACCOUNT"),
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
