package binance

import (
	"context"
	"errors"
	"strconv"
)

func (api *Api) GetBalanceFree(symbol string) (float64, error) {
	balances, err := api.Client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return 0, errors.New("не удалось получить баланс аккаунта: " + err.Error())
	}

	for _, balance := range balances.Balances {
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

func (api *Api) GetBalanceLocked(symbol string) (float64, error) {
	balances, err := api.Client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return 0, errors.New("не удалось получить баланс аккаунта: " + err.Error())
	}

	for _, balance := range balances.Balances {
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

func (api *Api) GetBalanceOverall(symbol string) (float64, error) {
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
