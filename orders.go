package goftx

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/grishinsana/goftx/models"
)

const (
	apiOrders                  = "/orders"
	apiGetOrdersHistory        = "/orders/history"
	apiModifyOrder             = "/orders/%d/modify"
	apiModifyOrderByClientID   = "/orders/by_client_id/%d/modify"
	apiTriggerOrders           = "/conditional_orders"
	apiGetOrderTriggers        = "/conditional_orders/%d/triggers"
	apiGetTriggerOrdersHistory = "/conditional_orders/history"
	apiModifyTriggerOrder      = "/conditional_orders/%d/modify"
)

type Orders struct {
	client *Client
}

func (o *Orders) GetOpenOrders(market string) ([]*models.Order, error) {
	requestParams := Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", o.client.apiURL, apiOrders),
	}
	if market != "" {
		requestParams.Params = map[string]string{
			"market": market,
		}
	}

	request, err := o.client.prepareRequest(requestParams)
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
		URL:    fmt.Sprintf("%s%s", o.client.apiURL, apiGetOrdersHistory),
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
		URL:    fmt.Sprintf("%s%s", o.client.apiURL, apiTriggerOrders),
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
		URL:    fmt.Sprintf("%s%s", o.client.apiURL, fmt.Sprintf(apiGetOrderTriggers, orderID)),
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

func (o *Orders) GetTriggerOrdersHistory(params *models.GetTriggerOrdersHistoryParams) ([]*models.TriggerOrder, error) {
	queryParams, err := PrepareQueryParams(params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", o.client.apiURL, apiGetTriggerOrdersHistory),
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

func (o *Orders) PlaceOrder(payload *models.PlaceOrderPayload) (*models.Order, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodPost,
		URL:    fmt.Sprintf("%s%s", o.client.apiURL, apiOrders),
		Body:   body,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result *models.Order
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (o *Orders) PlaceTriggerOrder(payload *models.PlaceTriggerOrderPayload) (*models.TriggerOrder, error) {
	err := payload.Validate()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodPost,
		URL:    fmt.Sprintf("%s%s", o.client.apiURL, apiTriggerOrders),
		Body:   body,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result *models.TriggerOrder
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (o *Orders) ModifyOrder(payload *models.ModifyOrderPayload, orderID int64) (*models.Order, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodPost,
		URL:    fmt.Sprintf("%s%s", o.client.apiURL, fmt.Sprintf(apiModifyOrder, orderID)),
		Body:   body,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result *models.Order
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (o *Orders) ModifyOrderByClientID(payload *models.ModifyOrderPayload, clientOrderID int64) (*models.Order, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodPost,
		URL:    fmt.Sprintf("%s%s", o.client.apiURL, fmt.Sprintf(apiModifyOrderByClientID, clientOrderID)),
		Body:   body,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result *models.Order
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (o *Orders) ModifyTriggerOrder(payload *models.ModifyTriggerOrderPayload, orderID int64) (*models.TriggerOrder, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodPost,
		URL:    fmt.Sprintf("%s%s", o.client.apiURL, fmt.Sprintf(apiModifyTriggerOrder, orderID)),
		Body:   body,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result *models.TriggerOrder
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (o *Orders) GetOrder(orderID int64) (*models.Order, error) {
	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s/%d", o.client.apiURL, apiOrders, orderID),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result *models.Order
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (o *Orders) GetOrderByClientID(clientOrderID int64) (*models.Order, error) {
	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s/by_client_id/%d", o.client.apiURL, apiOrders, clientOrderID),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response, err := o.client.do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result *models.Order
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (o *Orders) CancelOrder(orderID int64) error {
	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodDelete,
		URL:    fmt.Sprintf("%s%s/%d", o.client.apiURL, apiOrders, orderID),
	})
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = o.client.do(request)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (o *Orders) CancelOrderByClientID(clientOrderID string) error {
	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodDelete,
		URL:    fmt.Sprintf("%s%s/by_client_id/%s", o.client.apiURL, apiOrders, clientOrderID),
	})
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = o.client.do(request)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (o *Orders) CancelOpenTriggerOrder(triggerOrderID int64) error {
	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodDelete,
		URL:    fmt.Sprintf("%s%s/%d", o.client.apiURL, apiTriggerOrders, triggerOrderID),
	})
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = o.client.do(request)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (o *Orders) CancelAllOrders(payload *models.CancelAllOrdersPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return errors.WithStack(err)
	}

	request, err := o.client.prepareRequest(Request{
		Auth:   true,
		Method: http.MethodDelete,
		URL:    fmt.Sprintf("%s%s", o.client.apiURL, apiOrders),
		Body:   body,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = o.client.do(request)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
