package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"./tracking"
)

type Direction struct {
	Base                          string
	Quote                         string
	Intervals                     []tracking.Interval
	PercentOfBudgetPerTransaction float64
}

type Config struct {
	API struct {
		Binance struct {
			Key    string
			Secret string
		}
	}
	Directions []Direction
}

func GetConfig(path string) (*Config, error) {
	raw, err := ioutil.ReadFile(filepath.Dir(os.Args[0]) + path)
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
