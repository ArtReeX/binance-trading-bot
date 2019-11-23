package binance

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/adshao/go-binance"
)

// GetOpenOrders - получение списка открытых отдеров
func (api *API) GetOpenOrders(pair string) ([]Order, error) {
	openOrders, err := api.Client.NewListOpenOrdersService().Symbol(pair).Do(context.Background())
	if err != nil {
		return nil, errors.New("не удалось получить открытые ордера: " + err.Error())
	}

	formattedOpenOrders := make([]Order, len(openOrders))
	for index, order := range openOrders {
		formattedOpenOrders[index] = formatOrder(*order)
	}
	return formattedOpenOrders, nil
}

// GetOrder - получение ордера
func (api *API) GetOrder(pair string, id uint64) (Order, error) {
	order, err := api.Client.NewGetOrderService().Symbol(pair).OrderID(int64(id)).Do(context.Background())
	if err != nil {
		return Order{}, errors.New("не удалось получить ордер: " + err.Error())
	}

	return formatOrder(*order), nil
}

// CancelOrder - отмена ордера
func (api *API) CancelOrder(pair string, id uint64) error {
	_, err := api.Client.NewCancelOrderService().Symbol(pair).OrderID(int64(id)).Do(context.Background())
	if err != nil {
		return errors.New("не удалось отменить ордер: " + err.Error())
	}
	return nil
}

// CreateMarketSellOrder - создание ордера на продажу по рыночной цене
func (api *API) CreateMarketSellOrder(pair string, quantity float64) (uint64, error) {
	order, err := api.Client.NewCreateOrderService().
		Symbol(pair).
		Side(binance.SideTypeSell).
		Type(binance.OrderTypeMarket).
		Quantity(strconv.FormatFloat(quantity, 'f', int(api.Pairs[pair].QuantityAccuracy), 64)).
		NewOrderRespType(binance.NewOrderRespTypeRESULT).
		Do(context.Background())

	if err != nil {
		return 0, errors.New("не удалось создать MARKET ордер на продажу: " + err.Error())
	}
	return uint64(order.OrderID), nil
}

// CreateMarketBuyOrder - создание ордера на покупку по рыночной цене
func (api *API) CreateMarketBuyOrder(pair string, quantity float64) (uint64, error) {
	order, err := api.Client.NewCreateOrderService().
		Symbol(pair).
		Side(binance.SideTypeBuy).
		Type(binance.OrderTypeMarket).
		Quantity(strconv.FormatFloat(quantity, 'f', int(api.Pairs[pair].QuantityAccuracy), 64)).
		NewOrderRespType(binance.NewOrderRespTypeRESULT).
		Do(context.Background())

	if err != nil {
		return 0, errors.New("не удалось создать MARKET ордер на покупку: " + err.Error())
	}
	return uint64(order.OrderID), nil
}

// CreateStopLimitSellOrder - создание STOP-LOSS ордера на продажу
func (api *API) CreateStopLimitSellOrder(pair string, quantity float64, price float64,
	stopPrice float64) (uint64, error) {
	order, err := api.Client.NewCreateOrderService().
		Symbol(pair).
		Side(binance.SideTypeSell).
		Type(binance.OrderTypeStopLossLimit).
		Quantity(strconv.FormatFloat(quantity, 'f', int(api.Pairs[pair].QuantityAccuracy), 64)).
		Price(strconv.FormatFloat(price, 'f', int(api.Pairs[pair].PriceAccuracy), 64)).
		StopPrice(strconv.FormatFloat(stopPrice, 'f', int(api.Pairs[pair].PriceAccuracy), 64)).
		TimeInForce(binance.TimeInForceTypeGTC).
		NewOrderRespType(binance.NewOrderRespTypeRESULT).
		Do(context.Background())
	if err != nil {
		return 0, errors.New("не удалось создать STOP-LIMIT ордер на продажу: " + err.Error())
	}
	return uint64(order.OrderID), nil
}

// CreateStopLimitBuyOrder - создание STOP-LOSS ордера на покупку
func (api *API) CreateStopLimitBuyOrder(pair string, quantity float64, price float64,
	stopPrice float64) (uint64, error) {
	order, err := api.Client.NewCreateOrderService().
		Symbol(pair).
		Side(binance.SideTypeBuy).
		Type(binance.OrderTypeLimit).
		Quantity(strconv.FormatFloat(quantity, 'f', int(api.Pairs[pair].QuantityAccuracy), 64)).
		Price(strconv.FormatFloat(price, 'f', int(api.Pairs[pair].PriceAccuracy), 64)).
		StopPrice(strconv.FormatFloat(stopPrice, 'f', int(api.Pairs[pair].PriceAccuracy), 64)).
		TimeInForce(binance.TimeInForceTypeGTC).
		NewOrderRespType(binance.NewOrderRespTypeRESULT).
		Do(context.Background())
	if err != nil {
		return 0, errors.New("не удалось создать STOP-LIMIT ордер на покупку: " + err.Error())
	}
	return uint64(order.OrderID), nil
}

// CreateLimitSellOrder - создание ордера на продажу
func (api *API) CreateLimitSellOrder(pair string, quantity float64, price float64) (uint64, error) {
	order, err := api.Client.NewCreateOrderService().
		Symbol(pair).
		Side(binance.SideTypeSell).
		Type(binance.OrderTypeLimit).
		Quantity(strconv.FormatFloat(quantity, 'f', int(api.Pairs[pair].QuantityAccuracy), 64)).
		Price(strconv.FormatFloat(price, 'f', int(api.Pairs[pair].PriceAccuracy), 64)).
		TimeInForce(binance.TimeInForceTypeFOK).
		NewOrderRespType(binance.NewOrderRespTypeRESULT).
		Do(context.Background())
	if err != nil {
		return 0, errors.New("не удалось создать LIMIT ордер на продажу: " + err.Error())
	}
	return uint64(order.OrderID), nil
}

// CreateLimitBuyOrder - создание ордера на покупку
func (api *API) CreateLimitBuyOrder(pair string, quantity float64, price float64) (uint64, error) {
	order, err := api.Client.NewCreateOrderService().
		Symbol(pair).
		Side(binance.SideTypeBuy).
		Type(binance.OrderTypeLimit).
		Quantity(strconv.FormatFloat(quantity, 'f', int(api.Pairs[pair].QuantityAccuracy), 64)).
		Price(strconv.FormatFloat(price, 'f', int(api.Pairs[pair].PriceAccuracy), 64)).
		TimeInForce(binance.TimeInForceTypeFOK).
		NewOrderRespType(binance.NewOrderRespTypeRESULT).
		Do(context.Background())
	if err != nil {
		return 0, errors.New("не удалось создать LIMIT ордер на покупку: " + err.Error())
	}
	return uint64(order.OrderID), nil
}

// GetFinalOrder - получение ордера
func (api *API) GetFinalOrder(pair string, id uint64) (Order, error) {
	for {
		order, err := api.GetOrder(pair, id)
		if err != nil {
			log.Println("не удалось получить конечный статус ордера:" + err.Error())
			continue
		}
		if order.Status != OrderStatusNew && order.Status != OrderStatusPartiallyFilled {
			return order, nil
		}
	}
}
