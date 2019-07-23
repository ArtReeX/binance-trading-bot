package binance

import (
	"context"
	"errors"
	"strconv"

	"github.com/adshao/go-binance"
)

// API - интерфейс клиента
type API struct {
	client *binance.Client
}

// NewClient - функция создания нового клиента для биржи
func NewClient(key string, secret string) *API {
	return &API{
		client: binance.NewClient(key, secret),
	}
}

// GetBalances - функция получения баланса аккаунта
func (api *API) GetBalances() ([]binance.Balance, error) {
	res, err := api.client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return nil, errors.New("не удалось получить баланс аккаунта: " + err.Error())
	}

	// получение исключительно положительных балансов
	balances := []binance.Balance{}
	for _, value := range res.Balances {
		free, err := strconv.ParseFloat(value.Free, 64)
		if err != nil {
			return nil, errors.New("не удалось преобразовать строку свободного баланса в дробь")
		}

		locked, err := strconv.ParseFloat(value.Locked, 64)
		if err != nil {
			return nil, errors.New("не удалось преобразовать строку заблокированного баланса в дробь")
		}

		if free+locked > 0 {
			balances = append(balances, value)
		}
	}
	return balances, nil
}

// GetBalanceFree - функция получения свободного баланса определённой валюты аккаунта
func (api *API) GetBalanceFree(symbol string) (float64, error) {
	balances, err := api.GetBalances()
	if err != nil {
		return 0, err
	}

	for _, balance := range balances {
		if balance.Asset == symbol {
			free, err := strconv.ParseFloat(balance.Free, 64)
			if err != nil {
				return 0, errors.New("не удалось преобразовать строку свободного баланса в дробь")
			}
			return free, nil
		}
	}

	return 0, nil
}

// GetBalanceLocked - функция получения заблокированного баланса определённой валюты аккаунта
func (api *API) GetBalanceLocked(symbol string) (float64, error) {
	balances, err := api.GetBalances()
	if err != nil {
		return 0, err
	}

	for _, balance := range balances {
		if balance.Asset == symbol {
			locked, err := strconv.ParseFloat(balance.Locked, 64)
			if err != nil {
				return 0, errors.New("не удалось преобразовать строку заблокированного баланса в дробь")
			}
			return locked, nil
		}
	}

	return 0, nil
}

// GetBalanceOverall - функция получения общего баланса определённой валюты аккаунта
func (api *API) GetBalanceOverall(symbol string) (float64, error) {
	balanceFree, err := api.GetBalanceFree(symbol)
	if err != nil {
		return 0, err
	}
	balanceLocked, err := api.GetBalanceLocked(symbol)
	if err != nil {
		return 0, err
	}

	return balanceFree + balanceLocked, nil
}

// GetServerTime - функция получения времени сервера
func (api *API) GetServerTime() (int64, error) {
	time, err := api.client.NewServerTimeService().Do(context.Background())
	if err != nil {
		return 0, errors.New("не удалось получить время с сервера: " + err.Error())
	}
	return time, nil
}

// GetCandleHistory - функция получения истории цены для валюты
func (api *API) GetCandleHistory(symbol string, interval string) ([]*binance.Kline, error) {
	priceHistory, err := api.client.NewKlinesService().Symbol(symbol).Interval(interval).Do(context.Background())
	if err != nil {
		return nil, errors.New("не удалось получить историю валюты: " + err.Error())
	}

	return priceHistory, nil
}

// GetCurrentPrice - функция получения текущей цены валюты
func (api *API) GetCurrentPrice(pair string) (float64, error) {
	stat, err := api.client.NewPriceChangeStatsService().Symbol(pair).Do(context.Background())
	if err != nil {
		return 0, errors.New("не удалось получить текущую цену валюты: " + err.Error())
	}
	currentPrice, err := strconv.ParseFloat(stat.LastPrice, 64)
	if err != nil {
		return 0, errors.New("не удалось преобразовать строку текущей цены валюты в дробь: " + err.Error())
	}
	return currentPrice, nil
}

// GetOpenOrders - функция получения списка открытых ордеров
func (api *API) GetOpenOrders(pair string) ([]*binance.Order, error) {
	openOrders, err := api.client.NewListOpenOrdersService().Symbol(pair).Do(context.Background())
	if err != nil {
		return nil, errors.New("не удалось получить открытые ордера: " + err.Error())
	}
	return openOrders, nil
}

// GetOrder - функция получения ордера
func (api *API) GetOrder(pair string, id int64) (*binance.Order, error) {
	order, err := api.client.NewGetOrderService().Symbol(pair).OrderID(id).Do(context.Background())
	if err != nil {
		return nil, errors.New("не удалось получить ордер: " + err.Error())
	}
	return order, nil
}

// CancelOrder - функция отмены ордера
func (api *API) CancelOrder(pair string, id int64) (*binance.CancelOrderResponse, error) {
	res, err := api.client.NewCancelOrderService().Symbol(pair).OrderID(id).Do(context.Background())
	if err != nil {
		return nil, errors.New("не удалось отменить ордер: " + err.Error())
	}
	return res, nil
}

// CreateMarketCellOrder - функция создания MARKET ордера на продажу
func (api *API) CreateMarketCellOrder(pair string, quantity float64, accuracyQuantity uint8) (*binance.CreateOrderResponse, error) {

	order, err := api.client.NewCreateOrderService().
		Symbol(pair).
		Side(binance.SideTypeSell).
		Type(binance.OrderTypeMarket).
		Quantity(strconv.FormatFloat(quantity, 'f', int(accuracyQuantity), 64)).
		Do(context.Background())

	if err != nil {
		return nil, errors.New("не удалось создать MARKET ордер на продажу: " + err.Error())
	}
	return order, nil
}

// CreateMarketBuyOrder - функция создания MARKET ордера на покупку
func (api *API) CreateMarketBuyOrder(pair string, quantity float64, accuracyQuantity uint8) (*binance.CreateOrderResponse, error) {
	order, err := api.client.NewCreateOrderService().
		Symbol(pair).
		Side(binance.SideTypeBuy).
		Type(binance.OrderTypeMarket).
		Quantity(strconv.FormatFloat(quantity, 'f', int(accuracyQuantity), 64)).
		Do(context.Background())

	if err != nil {
		return nil, errors.New("не удалось создать MARKET ордер на покупку: " + err.Error())
	}
	return order, nil
}

// CreateLimitSellOrder - функция создания LIMIT ордера на продажу
func (api *API) CreateLimitSellOrder(pair string,
	quantity float64,
	price float64,
	accuracyQuantity uint8,
	accuracyPrice uint8) (*binance.CreateOrderResponse, error) {
	order, err := api.client.NewCreateOrderService().
		Symbol(pair).
		Side(binance.SideTypeSell).
		Type(binance.OrderTypeLimit).
		Quantity(strconv.FormatFloat(quantity, 'f', int(accuracyQuantity), 64)).
		Price(strconv.FormatFloat(price, 'f', int(accuracyPrice), 64)).
		TimeInForce(binance.TimeInForceGTC).
		Do(context.Background())

	if err != nil {
		return nil, errors.New("не удалось создать LIMIT ордер на продажу: " + err.Error())
	}
	return order, nil
}

// CreateLimitBuyOrder - функция создания LIMIT ордера на покупку
func (api *API) CreateLimitBuyOrder(
	pair string,
	quantity float64, price float64,
	accuracyQuantity uint8,
	accuracyPrice uint8) (*binance.CreateOrderResponse, error) {
	order, err := api.client.NewCreateOrderService().
		Symbol(pair).
		Side(binance.SideTypeBuy).
		Type(binance.OrderTypeLimit).
		Quantity(strconv.FormatFloat(quantity, 'f', int(accuracyQuantity), 64)).
		Price(strconv.FormatFloat(price, 'f', int(accuracyPrice), 64)).
		TimeInForce(binance.TimeInForceGTC).
		Do(context.Background())

	if err != nil {
		return nil, errors.New("не удалось создать LIMIT ордер на покупку: " + err.Error())
	}
	return order, nil
}
