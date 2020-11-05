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
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

const (
	apiUrl    = "https://ftx.com/api"
	apiOtcUrl = "https://otc.ftx.com/api"

	keyHeader  = "FTX-KEY"
	signHeader = "FTX-SIGN"
	tsHeader   = "FTX-TS"
)

type Option func(c *Client)

func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.client = client
	}
}

func WithAuth(key, secret string) Option {
	return func(c *Client) {
		c.apiKey = key
		c.secret = secret
	}
}

type Client struct {
	client         *http.Client
	apiKey         string
	secret         string
	serverTimeDiff time.Duration
	SubAccounts
	Markets
	Account
	Stream
	Orders
}

func New(opts ...Option) *Client {
	client := &Client{
		client: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(client)
	}

	client.SubAccounts = SubAccounts{client: client}
	client.Markets = Markets{client: client}
	client.Account = Account{client: client}
	client.Orders = Orders{client: client}
	client.Stream = Stream{
		mu:                     &sync.Mutex{},
		url:                    wsUrl,
		dialer:                 websocket.DefaultDialer,
		wsReconnectionCount:    reconnectCount,
		wsReconnectionInterval: reconnectInterval,
	}

	return client
}

func (c *Client) SetServerTimeDiff() error {
	serverTime, err := c.GetServerTime()
	if err != nil {
		return errors.WithStack(err)
	}
	c.serverTimeDiff = serverTime.Sub(time.Now().UTC())
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

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(keyHeader, c.apiKey)
		req.Header.Set(signHeader, c.signture(payload))
		req.Header.Set(tsHeader, nonce)
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

func (c *Client) prepareQueryParams(params interface{}) map[string]string {
	result := make(map[string]string)

	val := reflect.ValueOf(params).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		result[tag.Get("json")] = valueField.String()
	}

	return result
}

func (c *Client) signture(payload string) string {
	mac := hmac.New(sha256.New, []byte(c.secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

func (c *Client) GetServerTime() (*time.Time, error) {
	request, err := c.prepareRequest(Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s/time", apiOtcUrl),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := c.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result time.Time
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &result, nil
}
