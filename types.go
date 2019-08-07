package main

type (
	Binance struct {
		Key    string  `json:"key"`
		Secret string  `json:"secret"`
		Fee    float64 `json:"fee"`
	}

	Api struct {
		Binance Binance `json:"binance"`
	}

	Direction struct {
		Pair                   string   `json:"pair"`
		Intervals              []string `json:"intervals"`
		PriceForOneTransaction float64  `json:"priceForOneTransaction"`
	}

	Config struct {
		Api        Api         `json:"api"`
		Directions []Direction `json:"directions"`
	}
)
