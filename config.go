package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

// GetConfig - получение конфигурации
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
