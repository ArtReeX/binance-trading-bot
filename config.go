package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

type Binance struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

type API struct {
	Binance Binance `json:"binance"`
}

type Direction struct {
	Base                   string  `json:"base"`
	Quote                  string  `json:"quote"`
	Intervals              string  `json:"intervals"`
	PriceForOneTransaction float64 `json:"priceForOneTransaction"`
}

type Config struct {
	API        API         `json:"api"`
	Directions []Direction `json:"directions"`
}

func GetConfig(path string) (*Config, error) {
	raw, err := ioutil.ReadFile(os.Getenv("HOME") + "/.binance-trading-bot/" + path)
	if err != nil {
		return nil, errors.New("Ошибка загрузки конфигурации: " + err.Error())
	}

	config := new(Config)

	err = json.Unmarshal(raw, &config)
	if err != nil {
		return nil, errors.New("Ошибка парсинга конфигурации: " + err.Error())
	}

	log.Println("Файл конфигурации успешно загружен.")

	return config, nil
}
