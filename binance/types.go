package binance

import "github.com/adshao/go-binance"

type (
	Order struct {
		Symbol                   string
		OrderId                  uint64
		Price                    float64
		OrigQuantity             float64
		ExecutedQuantity         float64
		CummulativeQuoteQuantity float64
		Status                   OrderStatus
		StopPrice                float64
	}

	OrderStatus string

	Api struct {
		Client *binance.Client
		Pairs  map[string]Pair
	}

	Pair struct {
		BaseAsset        string
		QuoteAsset       string
		QuantityAccuracy uint8
		PriceAccuracy    uint8
	}

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

	Depth struct {
		Bids []Bid
		Asks []Ask
	}

	Bid struct {
		Price    float64
		Quantity float64
	}

	Ask struct {
		Price    float64
		Quantity float64
	}
)
