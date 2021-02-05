package goftx

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/grishinsana/goftx/models"
)

const (
	apiBorrowRates    = "/spot_margin/borrow_rates"
	apiLendingRates   = "/spot_margin/lending_rates"
	apiBorrowSummary  = "/spot_margin/borrow_summary"
	apiMarketInfo     = "/spot_margin/market_info"
	apiBorrowHistory  = "/spot_margin/borrow_history"
	apiLendingHistory = "/spot_margin/lending_history"
	apiLendingOffers  = "/spot_margin/offers"
	apiLendingInfo    = "/spot_margin/lending_info"
)

type SpotMargin struct {
	client *Client
}

func (s *SpotMargin) GetBorrowRates() ([]*models.BorrowRate, error) {
	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiBorrowRates),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := s.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.BorrowRate
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s *SpotMargin) GetLendingRates() ([]*models.LendingRate, error) {
	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiLendingRates),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := s.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.LendingRate
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s *SpotMargin) GetDailyBorrowedAmounts() ([]*models.BorrowSummary, error) {
	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiBorrowSummary),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := s.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.BorrowSummary
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s *SpotMargin) GetMarketInfo(market string) ([]*models.GetSpotMarginMarketInfoResponse, error) {
	queryParams := map[string]string{
		"market": market,
	}

	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiMarketInfo),
		Params: queryParams,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := s.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.GetSpotMarginMarketInfoResponse
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s *SpotMargin) GetBorrowHistory() ([]*models.BorrowHistory, error) {
	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiBorrowHistory),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := s.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.BorrowHistory
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s *SpotMargin) GetLendingHistory() ([]*models.LendingHistory, error) {
	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiLendingHistory),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := s.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.LendingHistory
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s *SpotMargin) GetLendingOffers() ([]*models.LendingOffer, error) {
	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiLendingOffers),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := s.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.LendingOffer
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s *SpotMargin) GetLendingInfo() ([]*models.LendingInfo, error) {
	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiLendingInfo),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := s.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.LendingInfo
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s *SpotMargin) SubmitLendingOffer(payload *models.LendingOfferPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return errors.WithStack(err)
	}

	request, err := s.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodPost,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiLendingOffers),
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
