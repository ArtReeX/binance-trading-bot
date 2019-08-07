package binance

import (
	"context"
	"errors"
	"github.com/adshao/go-binance"
	"strconv"
	"strings"
)

func getPairs(client *binance.Client) (map[string]Pair, error) {
	info, err := client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		return nil, errors.New("не удалось получить точность валют" + err.Error())
	}

	pairs := make(map[string]Pair)
	for _, pair := range info.Symbols {
		pairs[pair.Symbol] = Pair{
			BaseAsset:        pair.BaseAsset,
			QuoteAsset:       pair.QuoteAsset,
			QuantityAccuracy: uint8(len(strings.Split(strings.TrimRight(pair.LotSizeFilter().StepSize, "0"), ".")[1])),
			PriceAccuracy:    uint8(len(strings.Split(strings.TrimRight(pair.PriceFilter().TickSize, "0"), ".")[1])),
		}
	}
	return pairs, nil
}

func (api *Api) GetDepth(pair string, limit uint8) (Depth, error) {
	depth, err := api.Client.NewDepthService().Symbol(pair).Limit(int(limit)).Do(context.Background())
	if err != nil {
		return Depth{}, errors.New("не удалось получить глубину валюты" + err.Error())
	}

	formattedDepth := Depth{Bids: make([]Bid, len(depth.Bids)), Asks: make([]Ask, len(depth.Asks))}

	for index, bid := range depth.Bids {
		price, _ := strconv.ParseFloat(bid.Price, 64)
		quantity, _ := strconv.ParseFloat(bid.Quantity, 64)

		formattedDepth.Bids[index] = Bid{
			Price:    price,
			Quantity: quantity,
		}
	}
	for index, ask := range depth.Asks {
		price, _ := strconv.ParseFloat(ask.Price, 64)
		quantity, _ := strconv.ParseFloat(ask.Quantity, 64)

		formattedDepth.Asks[index] = Ask{
			Price:    price,
			Quantity: quantity,
		}
	}

	return formattedDepth, nil
}
