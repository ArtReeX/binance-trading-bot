package main

import (
	"log"
	"time"

	bnc "./binance"
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
	client := bnc.NewClient(config.API.Binance.Key, config.API.Binance.Secret)

	// refactor
	var stopOrder int64
	// бесконечный анализ валюты
	for {
		// получение истории торгов по валюте
		candleHistory, err := client.GetCandleHistory("BTCUSDT", "1m")
		if err != nil {
			log.Println(err)
			continue
		}

		// преобразование данных для StochRSI
		closePrices, err := client.ConvertCandleHistory(candleHistory, bnc.Close)
		if err != nil {
			log.Println(err)
			continue
		}

		// получение данных индикатора StochRSI
		k, d, err := indicators.StochRSI(closePrices, 14, 9, 3, 3)
		if err != nil {
			log.Println(err)
			continue
		}

		// последние две свечи
		kCandleCurrent := k[len(k)-1]
		kCandlePrev := k[len(k)-2]
		dCandleCurrent := d[len(d)-1]
		dCandlePrev := d[len(d)-2]

		// если произошло пересечение быстрой прямой долгую снизу вверх в зоне перепроданности то выполняем покупку
		// если произошло пересечение быстрой прямой долгую сверху вниз в зоне перекупленности то выполняем продажу
		if kCandleCurrent > 80 &&
			dCandleCurrent > 80 &&
			kCandlePrev >= dCandlePrev &&
			kCandleCurrent < dCandleCurrent {

			// продажа
			openOrders, err := client.GetOpenOrders("BTCUSDT")
			if err != nil {
				log.Fatalln(err)
			}
			if len(openOrders) != 0 {
				// отмена открытого STOP-LOSS (LIMIT) ордера
				st, err := client.CancelOrder("BTCUSDT", stopOrder)
				if err != nil {
					log.Fatalln(err)
				}
				log.Println(st)

				// получение доступного свободного баланса для валюты
				balanceFree, err := client.GetBalanceFree("BTC")
				if err != nil {
					log.Fatalln(err)
				}

				// продажа валюты
				res, err := client.CreateMarketCellOrder("BTCUSDT", balanceFree, 6)
				if err != nil {
					log.Fatalln(err)
				}
				log.Println("cell")
				log.Println(res)
			}

		} else if kCandleCurrent < 20 &&
			dCandleCurrent < 20 &&
			kCandlePrev <= dCandlePrev &&
			kCandleCurrent > dCandleCurrent {

			// покупка

			// получение открытых ордеров
			openOrders, err := client.GetOpenOrders("BTCUSDT")
			if err != nil {
				log.Fatalln(err)
			}
			if len(openOrders) == 0 {
				// получение доступного свободного баланса валюты для покупки
				balanceFree, err := client.GetBalanceFree("USDT")
				if err != nil {
					log.Fatalln(err)
				}
				// получение текущей цены валюты
				price, err := client.GetCurrentPrice("BTCUSDT")
				if err != nil {
					log.Fatalln(err)
				}

				// создание ордера для покупки
				res, err := client.CreateMarketBuyOrder("BTCUSDT", (balanceFree/price)*0.1, 6)
				if err != nil {
					log.Fatalln(err)
				}
				log.Println(res)
				// получение доступного заблокированного баланса купленной валюты
				balanceLock, err := client.GetBalanceLocked("BTC")
				if err != nil {
					log.Fatalln(err)
				}
				// установка STOP-LOSS (LIMIT) ордера
				st, err := client.CreateLimitSellOrder("BTCUSDT", balanceLock, price-(price*0.001), 6, 2)
				if err != nil {
					log.Fatalln(err)
				}
				log.Println(st)
				stopOrder = st.OrderID
				log.Println("buy")
			}
		}

		time.Sleep(time.Second)
	}

}
