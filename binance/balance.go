package binance

import (
	"context"
	"errors"
	"strconv"

	"github.com/adshao/go-binance"
)

// GetBalances - функция получения баланса аккаунта
func (api *API) GetBalances() ([]binance.Balance, error) {
	res, err := api.client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return nil, errors.New("Не удалось получить баланс аккаунта: " + err.Error())
	}

	// получение исключительно положительных балансов
	var balances []binance.Balance
	for _, value := range res.Balances {
		free, err := strconv.ParseFloat(value.Free, 64)
		if err != nil {
			return nil, errors.New("Не удалось преобразовать строку свободного баланса в дробь")
		}

		locked, err := strconv.ParseFloat(value.Locked, 64)
		if err != nil {
			return nil, errors.New("Не удалось преобразовать строку заблокированного баланса в дробь")
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
				return 0, errors.New("Не удалось преобразовать строку свободного баланса в дробь")
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
				return 0, errors.New("Не удалось преобразовать строку заблокированного баланса в дробь")
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
