package main

import (
	"log"
	"sync"

	bnc "./binance"
	"./tracking"
)

func main() {
	log.Println("Запуск бота для торговли на Binance.")

	// получение настроек
	config, err := GetConfig("config.json")
	if err != nil {
		log.Fatalln("Невозможно загрузить файл конфигурации:", err.Error())
	}

	// создание клиента
	client, err := bnc.NewClient(config.Api.Binance.Key, config.Api.Binance.Secret)
	if err != nil {
		log.Fatalln("Невозможно создать клиент: " + err.Error())
	}

	// создание потоков для отслеживания направлений
	waitGroupDirectionTracking := new(sync.WaitGroup)
	for _, direction := range config.Directions {
		for _, interval := range direction.Intervals {
			waitGroupDirectionTracking.Add(1)
			go tracking.DirectionTracking(direction.Pair, interval, direction.PriceForOneTransaction,
				config.Api.Binance.Fee, client, waitGroupDirectionTracking)
		}
	}
	waitGroupDirectionTracking.Wait()
}
