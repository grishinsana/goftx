package goftx

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

const (
	apiUrlFormat = "https://ftx.%s/api"
	apiOtcUrl    = "https://otc.ftx.com/api"

	keyHeaderFormat        = "FTX%s-KEY"
	signHeaderFormat       = "FTX%s-SIGN"
	tsHeaderFormat         = "FTX%s-TS"
	subAccountHeaderFormat = "FTX%s-SUBACCOUNT"
)

type Option func(c *Client)

func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.client = client
	}
}

func WithFTXUS() Option {
	return func(c *Client) {
		c.isFtxUS = true
	}
}

func WithAuth(key, secret string, subAccount ...string) Option {
	return func(c *Client) {
		c.apiKey = key
		c.secret = secret
		c.Stream.apiKey = key
		c.Stream.secret = secret

		if len(subAccount) > 0 {
			c.subAccount = subAccount[0]
		}
	}
}

type Client struct {
	client         *http.Client
	apiKey         string
	secret         string
	subAccount     string
	serverTimeDiff time.Duration
	isFtxUS        bool
	apiURL         string
	SubAccounts
	Markets
	Account
	Stream
	Orders
	Fills
	Converts
	Futures
	SpotMargin
}

func New(opts ...Option) *Client {
	client := &Client{
		client: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(client)
	}

	domain := "com"
	if client.isFtxUS {
		domain = "us"
	}
	client.apiURL = fmt.Sprintf(apiUrlFormat, domain)

	client.SubAccounts = SubAccounts{client: client}
	client.Markets = Markets{client: client}
	client.Account = Account{client: client}
	client.Orders = Orders{client: client}
	client.Fills = Fills{client: client}
	client.Converts = Converts{client: client}
	client.Futures = Futures{client: client}
	client.SpotMargin = SpotMargin{client: client}
	client.Stream = Stream{
		apiKey:                 client.apiKey,
		secret:                 client.secret,
		subAccount:             client.subAccount,
		mu:                     &sync.Mutex{},
		url:                    fmt.Sprintf(wsUrlFormat, domain),
		dialer:                 websocket.DefaultDialer,
		wsReconnectionCount:    reconnectCount,
		wsReconnectionInterval: reconnectInterval,
		wsTimeout:              streamTimeout,
	}

	return client
}

func (c *Client) SetServerTimeDiff() error {
	serverTime, err := c.GetServerTime()
	if err != nil {
		return errors.WithStack(err)
	}
	c.serverTimeDiff = serverTime.Sub(time.Now().UTC())
	c.Stream.serverTimeDiff = c.serverTimeDiff
	return nil
}

type Response struct {
	Success bool            `json:"success"`
	Result  json.RawMessage `json:"result"`
	Error   string          `json:"error,omitempty"`
}

type Request struct {
	Auth    bool
	Method  string
	URL     string
	Headers map[string]string
	Params  map[string]string
	Body    []byte
}

func (c *Client) prepareRequest(request Request) (*http.Request, error) {
	req, err := http.NewRequest(request.Method, request.URL, bytes.NewBuffer(request.Body))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	query := req.URL.Query()
	for k, v := range request.Params {
		query.Add(k, v)
	}
	req.URL.RawQuery = query.Encode()

	if request.Auth {
		nonce := strconv.FormatInt(time.Now().UTC().Add(c.serverTimeDiff).Unix()*1000, 10)
		payload := nonce + req.Method + req.URL.Path
		if req.URL.RawQuery != "" {
			payload += "?" + req.URL.RawQuery
		}
		if len(request.Body) > 0 {
			payload += string(request.Body)
		}

		usPrefix := ""
		if c.isFtxUS {
			usPrefix = "US"
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(fmt.Sprintf(keyHeaderFormat, usPrefix), c.apiKey)
		req.Header.Set(fmt.Sprintf(signHeaderFormat, usPrefix), c.signture(payload))
		req.Header.Set(fmt.Sprintf(tsHeaderFormat, usPrefix), nonce)

		if c.subAccount != "" {
			req.Header.Set(fmt.Sprintf(subAccountHeaderFormat, usPrefix), c.subAccount)
		}
	}

	for k, v := range request.Headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

func (c *Client) do(req *http.Request) ([]byte, error) {
	resp, err := c.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var response Response
	err = json.Unmarshal(res, &response)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if !response.Success {
		return nil, errors.Errorf("Status Code: %d	Error: %v", resp.StatusCode, response.Error)
	}

	return response.Result, nil
}

// nolint:errcheck
func (c *Client) signture(payload string) string {
	mac := hmac.New(sha256.New, []byte(c.secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

func (c *Client) GetServerTime() (time.Time, error) {
	request, err := c.prepareRequest(Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s/time", apiOtcUrl),
	})
	if err != nil {
		return time.Time{}, errors.WithStack(err)
	}

	response, err := c.do(request)
	if err != nil {
		return time.Time{}, errors.WithStack(err)
	}

	var result time.Time
	err = json.Unmarshal(response, &result)
	if err != nil {
		return time.Time{}, errors.WithStack(err)
	}

	return result, nil
}

func (c Client) Ping() error {
	request, err := c.prepareRequest(Request{
		Method: http.MethodGet,
		URL:    c.apiURL,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	response, err := c.do(request)
	if err != nil {
		return errors.WithStack(err)
	}

	var result bool
	err = json.Unmarshal(response, &result)
	if err != nil {
		return errors.WithStack(err)
	}

	if !result {
		return errors.New("service unavailable")
	}

	return nil
}
