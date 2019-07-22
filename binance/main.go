package binance

import (
	"context"
	"errors"
	"strconv"

	"github.com/adshao/go-binance"
)

// NewClient - функция создания нового клиента для биржи
func NewClient(key string, secret string) *binance.Client {
	return binance.NewClient(key, secret)
}

// GetBalances - функция получения баланса аккаунта
func GetBalances(client *binance.Client) ([]binance.Balance, error) {
	res, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return nil, errors.New("не удаётся получить баланс аккаунта: " + err.Error())
	}

	// получение исключительно положительных балансов
	balances := []binance.Balance{}
	for _, value := range res.Balances {
		free, err := strconv.ParseFloat(value.Free, 64)
		if err != nil {
			return nil, errors.New("не удаётся преобразовать строку свободного баланса в дробь")
		}

		locked, err := strconv.ParseFloat(value.Locked, 64)
		if err != nil {
			return nil, errors.New("не удаётся преобразовать строку заблокированного баланса в дробь")
		}

		if free+locked > 0 {
			balances = append(balances, value)
		}
	}
	return balances, nil
}

// GetBalance - функция получения баланса опередлённой валюты аккаунта
func GetBalance(client *binance.Client, symbol string) (float64, error) {
	balances, err := GetBalances(client)
	if err != nil {
		return 0, err
	}

	for _, balance := range balances {
		if balance.Asset == symbol {
			free, err := strconv.ParseFloat(balance.Free, 64)
			if err != nil {
				return 0, errors.New("не удаётся преобразовать строку баланса в дробь")
			}
			return free, nil
		}
	}

	return 0, nil
}

// GetServerTime - функция получения времени сервера
func GetServerTime(client *binance.Client) (int64, error) {
	time, err := client.NewServerTimeService().Do(context.Background())
	if err != nil {
		return 0, errors.New("не удаётся получить время с сервера: " + err.Error())
	}
	return time, nil
}

// GetCandleHistory - функция получения истории цены для валюты
func GetCandleHistory(client *binance.Client, symbol string, interval string) ([]*binance.Kline, error) {
	priceHistory, err := client.NewKlinesService().Symbol(symbol).Interval(interval).Do(context.Background())
	if err != nil {
		return nil, errors.New("не удаётся получить историю валюты: " + err.Error())
	}

	return priceHistory, nil
}

// GetCurrentPrice - функция получения текущей цены валюты
func GetCurrentPrice(client *binance.Client, symbol string) (float64, error) {
	stat, err := client.NewPriceChangeStatsService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return 0, errors.New("не удаётся получить текущую цену валюты: " + err.Error())
	}
	lastPrice, err := strconv.ParseFloat(stat.LastPrice, 64)
	if err != nil {
		return 0, errors.New("не удаётся преобразовать строку цены валюты в дробь: " + err.Error())
	}
	return lastPrice, nil
}
