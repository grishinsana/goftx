package goftx

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/grishinsana/goftx/models"
)

const (
	apiSubaccounts           = "/subaccounts"
	apiChangeSubaccountName  = "/subaccounts/update_name"
	apiGetSubaccountBalances = "/subaccounts/%s/balances"
	apiTransfer              = "/subaccounts/transfer"
)

type SubAccounts struct {
	client *Client
}

func (s *SubAccounts) GetSubaccounts() ([]*models.SubAccount, error) {
	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiSubaccounts),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := s.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.SubAccount
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s *SubAccounts) CreateSubaccount(nickname string) (*models.SubAccount, error) {
	body, err := json.Marshal(struct {
		Nickname string `json:"nickname"`
	}{Nickname: nickname})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodPost,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiSubaccounts),
		Body:   body,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := s.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result models.SubAccount
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &result, nil
}

func (s *SubAccounts) ChangeSubaccount(nickname, newNickname string) error {
	body, err := json.Marshal(struct {
		Nickname    string `json:"nickname"`
		NewNickname string `json:"newNickname"`
	}{Nickname: nickname, NewNickname: newNickname})
	if err != nil {
		return errors.WithStack(err)
	}

	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodPost,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiChangeSubaccountName),
		Body:   body,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = s.client.do(request)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *SubAccounts) DeleteSubaccount(nickname string) error {
	body, err := json.Marshal(struct {
		Nickname string `json:"nickname"`
	}{Nickname: nickname})
	if err != nil {
		return errors.WithStack(err)
	}

	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodDelete,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiSubaccounts),
		Body:   body,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = s.client.do(request)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *SubAccounts) GetSubaccountBalances(nickname string) ([]*models.Balance, error) {
	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, fmt.Sprintf(apiGetSubaccountBalances, nickname)),
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

func (s *SubAccounts) Transfer(payload *models.TransferPayload) (*models.TransferResponse, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodPost,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiTransfer),
		Body:   body,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := s.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result models.TransferResponse
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &result, nil
}
