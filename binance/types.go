package binance

import "github.com/adshao/go-binance"

type (
	// Order - структура ордера
	Order struct {
		Symbol                   string
		OrderID                  uint64
		Price                    float64
		OrigQuantity             float64
		ExecutedQuantity         float64
		CummulativeQuoteQuantity float64
		Status                   OrderStatus
		StopPrice                float64
	}

	// OrderStatus - статус ордера
	OrderStatus string

	// API - клиент
	API struct {
		Client *binance.Client
		Pairs  map[string]Pair
	}

	// Pair - направление торговли
	Pair struct {
		BaseAsset        string
		QuoteAsset       string
		QuantityAccuracy uint8
		PriceAccuracy    uint8
	}

	// Candle - структура свечи
	Candle struct {
		Open                     float64
		High                     float64
		Low                      float64
		Close                    float64
		Volume                   float64
		QuoteAssetVolume         float64
		TakerBuyBaseAssetVolume  float64
		TakerBuyQuoteAssetVolume float64
	}

	// Depth - структура стакана
	Depth struct {
		Bids []Bid
		Asks []Ask
	}

	// Bid - структура стакана на продажу
	Bid struct {
		Price    float64
		Quantity float64
	}

	// Ask - структура стакана на покупку
	Ask struct {
		Price    float64
		Quantity float64
	}
)
