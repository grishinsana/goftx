package goftx

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/grishinsana/goftx/models"
)

const apiQuotes = "/otc/quotes"

type Converts struct {
	client *Client
}

func (c *Converts) CreateQuote(payload *models.CreateQuotePayload) (int64, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	request, err := c.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodPost,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiQuotes),
		Body:   body,
	})
	if err != nil {
		return 0, errors.WithStack(err)
	}

	response, err := c.client.do(request)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	var result struct {
		QuoteId int64 `json:"quoteId"`
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return result.QuoteId, nil
}

func (c *Converts) GetQuotes(quoteID int64, market *string) ([]*models.QuoteStatus, error) {
	queryParams := make(map[string]string)
	if market != nil {
		queryParams["market"] = *market
	}

	request, err := c.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s/%d", apiUrl, apiQuotes, quoteID),
		Params: queryParams,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := c.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.QuoteStatus
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (c *Converts) AcceptQuote(quoteID int64) error {
	request, err := c.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodPost,
		URL:    fmt.Sprintf("%s%s/%d/accept", apiUrl, apiQuotes, quoteID),
	})
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = c.client.do(request)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
