package main

import (
	"log"
	"time"

	binance "./binance"
	indicators "./indicators"
)

func main() {
	log.Println("Запуск бота для торговли на Binance.")

	// получение настроек
	config, err := GetConfig("config.json")
	if err != nil {
		log.Fatalln("Невозможно загрузить файл конфигурации:", err.Error())
	}

	// создание клиента
	var client = binance.NewClient(config.API.Binance.Key, config.API.Binance.Secret)

	// trash
	for {
		// получение истории торгов по валюте
		candleHistory, err := binance.GetCandleHistory(client, "BTCUSDT", "1m")
		if err != nil {
			log.Println(err)
			continue
		}

		closePrices, err := binance.ConvertCandleHistory(candleHistory, binance.Close)
		if err != nil {
			log.Println(err)
			continue
		}

		k, d, err := indicators.StochRSI(closePrices, 14, 9, 3, 3)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println(k[len(k)-1], d[len(d)-1])

		// последние две свечи
		kCandleCurrent := k[len(k)-1]
		kCandlePrev := k[len(k)-2]
		dCandleCurrent := d[len(d)-1]
		dCandlePrev := d[len(d)-2]

		// если произошло пересечение быстрой прямой долгую снизу вверх в зоне перепроданности то выполняем покупку
		// если произошло пересечение быстрой прямой долгую сверху вниз в зоне перекупленности то выполняем продажу
		if kCandleCurrent > 80 &&
			dCandleCurrent > 80 &&
			kCandlePrev > dCandlePrev &&
			kCandleCurrent < dCandleCurrent {
			log.Println("cell")
		} else if kCandleCurrent < 20 &&
			dCandleCurrent < 20 &&
			kCandlePrev < dCandlePrev &&
			kCandleCurrent > dCandleCurrent {
			log.Println("buy")
		}
		time.Sleep(time.Second)
	}

}
