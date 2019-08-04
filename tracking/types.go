package tracking

import "sync"

type IndicatorsStatus int8

type BotStatus uint8

type Direction struct {
	Base                   string
	Quote                  string
	Interval               string
	PriceForOneTransaction float64
}

type Bot struct {
	BuyOrderId           uint64
	SellOrderId          uint64
	StopLossOrderId      uint64
	StopLossOrderIdMutex sync.Mutex

	Status      BotStatus
	StatusMutex sync.Mutex
}

type Candle struct {
	Open                     float64
	High                     float64
	Low                      float64
	Close                    float64
	Volume                   float64
	QuoteAssetVolume         float64
	TakerBuyBaseAssetVolume  float64
	TakerBuyQuoteAssetVolume float64
}
