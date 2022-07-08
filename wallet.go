package goftx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/grishinsana/goftx/models"
)

const (
	apiBalances = "/wallet/balances"
	apiWithdraw = "/wallet/withdrawals"
)

type Wallet struct {
	client *Client
}

func (s *Wallet) GetBalances() ([]*models.Balance, error) {
	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", s.client.apiURL, apiBalances),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := s.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.Balance
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s *Wallet) Withdraw(ctx context.Context, payload *models.CreateWithdrawPayload) (*models.CreateWithdrawResult, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodPost,
		URL:    fmt.Sprintf("%s%s", s.client.apiURL, apiWithdraw),
		Body:   body,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := s.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result *models.CreateWithdrawResult
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}
