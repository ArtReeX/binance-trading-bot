package main

type (
	// Binance - настройки Binance
	Binance struct {
		Key    string  `json:"key"`
		Secret string  `json:"secret"`
		Fee    float64 `json:"fee"`
	}

	// API - настройки API
	API struct {
		Binance Binance `json:"binance"`
	}

	// Direction - настройки торговли
	Direction struct {
		Pair                   string   `json:"pair"`
		Intervals              []string `json:"intervals"`
		PriceForOneTransaction float64  `json:"priceForOneTransaction"`
	}

	// Config - настройки
	Config struct {
		API        API         `json:"api"`
		Directions []Direction `json:"directions"`
	}
)
