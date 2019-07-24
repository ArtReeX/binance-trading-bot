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
			go tracking.DirectionTracking(direction.Base,
				direction.Quote,
				direction.AccuracyQuantity,
				direction.AccuracyPrice,
				interval,
				client,
				waitGroupDirectionTracking)
		}
	}
	waitGroupDirectionTracking.Wait()

	/*
		// идентификаторы ордеров
		ordersWithoutStopLoss, stopLossOrders := make(map[string][]int64), make(map[string][]int64)

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

			// открытие STOP-LIMIT ордеров на завершённые сделки
			for index, id := range ordersWithoutStopLoss["BTCUSDT"] {

				// получение статуса ордера
				order, err := client.GetOrder("BTCUSDT", id)
				if err != nil {
					log.Fatalln(err)
				}

				// установка STOP-LIMIT ордера в случае если обычный ордер был исполнен
				if order.Status == "FILLED" {
					// получение количества купленной валюты
					quantity, err := strconv.ParseFloat(order.OrigQuantity, 64)
					if err != nil {
						log.Fatalln(err)
					}
					// получение цены купленной валюты
					price, err := strconv.ParseFloat(order.Price, 64)
					if err != nil {
						log.Fatalln(err)
					}
					// установка STOP-LIMIT ордера
					stopOrder, err := client.CreateStopLimitSellOrder("BTCUSDT", quantity, price-(price*0.01), price-(price*0.01), 6, 2)
					if err != nil {
						log.Fatalln(err)

					}

					// добавление ордера в открытые стоп-ордера
					stopLossOrders["BTCUSDT"] = append(stopLossOrders["BTCUSDT"], stopOrder.OrderID)
					// удаление текущего исполненного ордера
					copy(ordersWithoutStopLoss["BTCUSDT"][index:], ordersWithoutStopLoss["BTCUSDT"][index+1:])
					ordersWithoutStopLoss["BTCUSDT"] = ordersWithoutStopLoss["BTCUSDT"][:len(ordersWithoutStopLoss["BTCUSDT"])-1]

					log.Println("Открыт STOP-LIMIT ордер на продажу по цене", stopOrder.Price)
				} else if order.Status != "PARTIALLY_FILLED" {
					// удаление текущего исполненного ордера
					copy(ordersWithoutStopLoss["BTCUSDT"][index:], ordersWithoutStopLoss["BTCUSDT"][index+1:])
					ordersWithoutStopLoss["BTCUSDT"] = ordersWithoutStopLoss["BTCUSDT"][:len(ordersWithoutStopLoss["BTCUSDT"])-1]

					log.Println("Удалён не исполненный STOP-LIMIT ордер из списка наблюдения")
				}

			}

			// удаление исполненных STOP-LIMIT ордеров
			for index, id := range stopLossOrders["BTCUSDT"] {

				// получение статуса ордера
				order, err := client.GetOrder("BTCUSDT", id)
				if err != nil {
					log.Fatalln(err)
				}

				if order.Status != "NEW" && order.Status != "PARTIALLY_FILLED" {
					// удаление текущего исполненного ордера
					copy(stopLossOrders["BTCUSDT"][index:], stopLossOrders["BTCUSDT"][index+1:])
					stopLossOrders["BTCUSDT"] = stopLossOrders["BTCUSDT"][:len(stopLossOrders["BTCUSDT"])-1]

					log.Println("Удалён исполненный STOP-LIMIT ордер из списка наблюдения", order.Price)
				}
			}

			// если произошло пересечение быстрой прямой долгую снизу вверх в зоне перепроданности то выполняем покупку
			// если произошло пересечение быстрой прямой долгую сверху вниз в зоне перекупленности то выполняем продажу
			if kCandleCurrent > 80 &&
				dCandleCurrent > 80 &&
				kCandlePrev >= dCandlePrev &&
				kCandleCurrent < dCandleCurrent {

				// продажа валюты
				if len(stopLossOrders["BTCUSDT"]) != 0 {

					// отмена открытых STOP-LIMIT ордеров
					for index, id := range stopLossOrders["BTCUSDT"] {
						_, err := client.CancelOrder("BTCUSDT", id)
						if err != nil {
							log.Fatalln(err)
						}
						// удаление из массива STOP-LIMIT ордеров
						copy(stopLossOrders["BTCUSDT"][index:], stopLossOrders["BTCUSDT"][index+1:])
						stopLossOrders["BTCUSDT"] = stopLossOrders["BTCUSDT"][:len(stopLossOrders["BTCUSDT"])-1]
					}

					// получение доступного свободного баланса для валюты
					balanceFree, err := client.GetBalanceFree("BTC")
					if err != nil {
						log.Fatalln(err)
					}

					// получение текущей цены валюты
					currentPrice, err := client.GetCurrentPrice("BTCUSDT")
					if err != nil {
						log.Fatalln(err)
					}

					// открытие ордера на продажу
					order, err := client.CreateLimitSellOrder("BTCUSDT", balanceFree, currentPrice, 6, 2)
					if err != nil {
						log.Fatalln(err)
					}

					log.Println("Открыт ордер на продажу по цене", order.Price)
				}

			} else if kCandleCurrent < 20 &&
				dCandleCurrent < 20 &&
				kCandlePrev <= dCandlePrev &&
				kCandleCurrent > dCandleCurrent {

				// покупка валюты
				if len(ordersWithoutStopLoss["BTCUSDT"]) == 0 && len(stopLossOrders["BTCUSDT"]) == 0 {
					// получение доступного свободного баланса валюты для покупки
					balanceFree, err := client.GetBalanceFree("USDT")
					if err != nil {
						log.Fatalln(err)
					}
					// получение текущей цены валюты
					currentPrice, err := client.GetCurrentPrice("BTCUSDT")
					if err != nil {
						log.Fatalln(err)
					}
					// создание ордера для покупки
					order, err := client.CreateLimitBuyOrder("BTCUSDT", (balanceFree/currentPrice)*0.1, currentPrice, 6, 2)
					if err != nil {
						log.Fatalln(err)
					}

					// добавление ордера на дальнейшую установку STOP-LIMIT ордера
					ordersWithoutStopLoss["BTCUSDT"] = append(ordersWithoutStopLoss["BTCUSDT"], order.OrderID)

					log.Println("Открыт ордер на покупку по цене", order.Price)
				}
			}

			time.Sleep(time.Second)
		}
	*/
}
