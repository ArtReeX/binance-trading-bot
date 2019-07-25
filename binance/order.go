package binance

import (
	"context"
	"errors"
	"strconv"

	"github.com/adshao/go-binance"
)

// GetOpenOrders - функция получения списка открытых ордеров
func (api *API) GetOpenOrders(pair string) ([]*binance.Order, error) {
	openOrders, err := api.client.NewListOpenOrdersService().Symbol(pair).Do(context.Background())
	if err != nil {
		return nil, errors.New("Не удалось получить открытые ордера: " + err.Error())
	}
	return openOrders, nil
}

// GetOrder - функция получения ордера
func (api *API) GetOrder(pair string, id int64) (*binance.Order, error) {
	order, err := api.client.NewGetOrderService().Symbol(pair).OrderID(id).Do(context.Background())
	if err != nil {
		return nil, errors.New("Не удалось получить ордер: " + err.Error())
	}
	return order, nil
}

// CancelOrder - функция отмены ордера
func (api *API) CancelOrder(pair string, id int64) (*binance.CancelOrderResponse, error) {
	order, err := api.client.NewCancelOrderService().Symbol(pair).OrderID(id).Do(context.Background())
	if err != nil {
		return nil, errors.New("Не удалось отменить ордер: " + err.Error())
	}
	return order, nil
}

// CreateMarketSellOrder - функция создания MARKET ордера на продажу
func (api *API) CreateMarketSellOrder(pair string, quantity float64) (*binance.CreateOrderResponse, error) {
	order, err := api.client.NewCreateOrderService().
		Symbol(pair).
		Side(binance.SideTypeSell).
		Type(binance.OrderTypeMarket).
		Quantity(strconv.FormatFloat(quantity, 'f', 6, 64)).
		NewOrderRespType(binance.NewOrderRespTypeRESULT).
		Do(context.Background())

	if err != nil {
		return nil, errors.New("Не удалось создать MARKET ордер на продажу: " + err.Error())
	}
	return order, nil
}

// CreateMarketBuyOrder - функция создания MARKET ордера на покупку
func (api *API) CreateMarketBuyOrder(pair string, quantity float64) (*binance.CreateOrderResponse, error) {
	order, err := api.client.NewCreateOrderService().
		Symbol(pair).
		Side(binance.SideTypeBuy).
		Type(binance.OrderTypeMarket).
		Quantity(strconv.FormatFloat(quantity, 'f', 6, 64)).
		NewOrderRespType(binance.NewOrderRespTypeRESULT).
		Do(context.Background())

	if err != nil {
		return nil, errors.New("Не удалось создать MARKET ордер на покупку: " + err.Error())
	}
	return order, nil
}

// CreateStopLimitSellOrder - функция создания STOP-LIMIT ордера на продажу
func (api *API) CreateStopLimitSellOrder(pair string,
	quantity float64,
	price float64,
	stopPrice float64) (*binance.CreateOrderResponse, error) {
	order, err := api.client.NewCreateOrderService().
		Symbol(pair).
		Side(binance.SideTypeSell).
		Type(binance.OrderTypeStopLossLimit).
		Quantity(strconv.FormatFloat(quantity, 'f', 6, 64)).
		Price(strconv.FormatFloat(price, 'f', 2, 64)).
		StopPrice(strconv.FormatFloat(price, 'f', 2, 64)).
		TimeInForce(binance.TimeInForceGTC).
		NewOrderRespType(binance.NewOrderRespTypeRESULT).
		Do(context.Background())
	if err != nil {
		return nil, errors.New("Не удалось создать STOP-LIMIT ордер на продажу: " + err.Error())
	}
	return order, nil
}

// CreateStopLimitBuyOrder - функция создания STOP-LIMIT ордера на покупку
func (api *API) CreateStopLimitBuyOrder(
	pair string,
	quantity float64,
	price float64,
	stopPrice float64) (*binance.CreateOrderResponse, error) {
	order, err := api.client.NewCreateOrderService().
		Symbol(pair).
		Side(binance.SideTypeBuy).
		Type(binance.OrderTypeLimit).
		Quantity(strconv.FormatFloat(quantity, 'f', 6, 64)).
		Price(strconv.FormatFloat(price, 'f', 2, 64)).
		StopPrice(strconv.FormatFloat(price, 'f', 2, 64)).
		TimeInForce(binance.TimeInForceGTC).
		NewOrderRespType(binance.NewOrderRespTypeRESULT).
		Do(context.Background())
	if err != nil {
		return nil, errors.New("Не удалось создать STOP-LIMIT ордер на покупку: " + err.Error())
	}
	return order, nil
}

// CreateLimitSellOrder - функция создания LIMIT ордера на продажу
func (api *API) CreateLimitSellOrder(
	pair string,
	quantity float64,
	price float64) (*binance.CreateOrderResponse, error) {
	order, err := api.client.NewCreateOrderService().
		Symbol(pair).
		Side(binance.SideTypeSell).
		Type(binance.OrderTypeLimit).
		Quantity(strconv.FormatFloat(quantity, 'f', 6, 64)).
		Price(strconv.FormatFloat(price, 'f', 2, 64)).
		TimeInForce(binance.TimeInForceGTC).
		NewOrderRespType(binance.NewOrderRespTypeRESULT).
		Do(context.Background())
	if err != nil {
		return nil, errors.New("Не удалось создать LIMIT ордер на продажу: " + err.Error())
	}
	return order, nil
}

// CreateLimitBuyOrder - функция создания LIMIT ордера на продажу
func (api *API) CreateLimitBuyOrder(
	pair string,
	quantity float64,
	price float64) (*binance.CreateOrderResponse, error) {
	order, err := api.client.NewCreateOrderService().
		Symbol(pair).
		Side(binance.SideTypeBuy).
		Type(binance.OrderTypeLimit).
		Quantity(strconv.FormatFloat(quantity, 'f', 6, 64)).
		Price(strconv.FormatFloat(price, 'f', 2, 64)).
		TimeInForce(binance.TimeInForceIOC).
		NewOrderRespType(binance.NewOrderRespTypeRESULT).
		Do(context.Background())
	if err != nil {
		return nil, errors.New("Не удалось создать LIMIT ордер на покупку: " + err.Error())
	}
	return order, nil
}
