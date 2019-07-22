package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
)

// Config - структра содержащая параметры бота
type Config struct {
	API struct {
		Binance struct {
			Key    string
			Secret string
		}
	}
}

// GetConfig - функция получения настроек
func GetConfig(path string) (*Config, error) {
	raw, err := ioutil.ReadFile("config.json")
	if err != nil {
		return nil, errors.New("ошибка загрузки конфигурации: " + err.Error())
	}

	config := Config{}

	err = json.Unmarshal(raw, &config)
	if err != nil {
		return nil, errors.New("ошибка парсинга конфигурации: " + err.Error())
	}

	log.Println("Файл конфигурации успешно загружен.")

	return &config, nil
}
