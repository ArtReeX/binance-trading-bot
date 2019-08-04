package binance

import (
	"context"
	"errors"
)

func (api *API) GetServerTime() (int64, error) {
	time, err := api.client.NewServerTimeService().Do(context.Background())
	if err != nil {
		return 0, errors.New("не удалось получить время с сервера: " + err.Error())
	}
	return time, nil
}
