package main

import (
	"log"
	"sync"

	bnc "./binance"
	tracking "./tracking"
)

func main() {
	log.Println("Запуск бота для торговли на Binance.")

	// получение настроек
	config, err := GetConfig("config.json")
	if err != nil {
		log.Fatalln("Невозможно загрузить файл конфигурации:", err.Error())
	}

	// создание клиента
	client := bnc.NewClient(config.API.Binance.Key, config.API.Binance.Secret)

	// создание потоков для отслеживания направлений
	waitGroupDirectionTracking := new(sync.WaitGroup)
	for _, direction := range config.Directions {
		for _, interval := range direction.Intervals {
			waitGroupDirectionTracking.Add(1)
			go tracking.DirectionTracking(tracking.Direction{
				Base:                          direction.Base,
				Quote:                         direction.Quote,
				Interval:                      interval,
				PercentOfBudgetPerTransaction: direction.PercentOfBudgetPerTransaction},
				client,
				waitGroupDirectionTracking)
		}
	}
	waitGroupDirectionTracking.Wait()
}
