package goftx

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/grishinsana/goftx/models"
	"github.com/pkg/errors"
)

const (
	apiGetOpenOrders    = "/orders"
	apiGetOrdersHistory = "/orders/history"
	apiGetTriggerOrders = "/conditional_orders"
	apiGetOrderTriggers = "/conditional_orders/%d/triggers"
)

type Orders struct {
	client *Client
}

func (o *Orders) GetOpenOrders(market string) ([]*models.Order, error) {
	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiGetOpenOrders),
		Params: map[string]string{
			"market": market,
		},
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.Order
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (o *Orders) GetOrdersHistory(params *models.GetOrdersHistoryParams) ([]*models.Order, error) {
	queryParams, err := PrepareQueryParams(params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiGetOrdersHistory),
		Params: queryParams,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.Order
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (o *Orders) GetOpenTriggerOrders(params *models.GetOpenTriggerOrdersParams) ([]*models.TriggerOrder, error) {
	queryParams, err := PrepareQueryParams(params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, apiGetTriggerOrders),
		Params: queryParams,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.TriggerOrder
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (o *Orders) GetOrderTriggers(orderID int64) ([]*models.Trigger, error) {
	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", apiUrl, fmt.Sprintf(apiGetOrderTriggers, orderID)),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result []*models.Trigger
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}
