package main

import (
	"log"
	"strconv"
	"time"

	bnc "./binance"
	indicators "./indicators"
)

// Order -  структура ордера
type Order struct {
	Pair string
	ID   int64
}

func main() {
	log.Println("Запуск бота для торговли на Binance.")

	// получение настроек
	config, err := GetConfig("config.json")
	if err != nil {
		log.Fatalln("Невозможно загрузить файл конфигурации:", err.Error())
	}

	// создание клиента
	client := bnc.NewClient(config.API.Binance.Key, config.API.Binance.Secret)

	// идентификаторы ордеров
	var ordersWithoutStopLoss, stopLossOrders []Order

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

		// открытие STOP-LOSS ордеров на завершённые сделки
		for index, order := range ordersWithoutStopLoss {

			// получение статуса ордера
			order, err := client.GetOrder(order.Pair, order.ID)
			if err != nil {
				log.Fatalln(err)
			}

			// установка STOP-LOSS ордера в случае если обычный ордер был исполнен
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
				// установка STOP-LOSS ордера
				stopOrder, err := client.CreateLimitSellOrder("BTCUSDT", quantity, price-(price*0.001), 6, 2)
				if err != nil {
					log.Fatalln(err)

				}

				log.Println("Открыт STOP-LOSS ордер на продажу по цене", stopOrder.Price)

				// добавление ордера в открытые стоп-ордера
				stopLossOrders = append(stopLossOrders, Order{Pair: stopOrder.Symbol, ID: stopOrder.OrderID})
				// удаление текущего исполненного ордера
				copy(ordersWithoutStopLoss[index:], ordersWithoutStopLoss[index+1:])
				ordersWithoutStopLoss = ordersWithoutStopLoss[:len(ordersWithoutStopLoss)-1]
			}
		}

		// удаление исполненных STOP-LOSS ордеров
		for index, order := range stopLossOrders {

			// получение статуса ордера
			order, err := client.GetOrder(order.Pair, order.ID)
			if err != nil {
				log.Fatalln(err)
			}

			if order.Status == "FILLED" {
				// удаление текущего исполненного ордера
				copy(stopLossOrders[index:], stopLossOrders[index+1:])
				stopLossOrders = stopLossOrders[:len(stopLossOrders)-1]

				log.Println("Удалён исполненный STOP-LOSS ордер по цене", order.Price)
			}
		}

		// если произошло пересечение быстрой прямой долгую снизу вверх в зоне перепроданности то выполняем покупку
		// если произошло пересечение быстрой прямой долгую сверху вниз в зоне перекупленности то выполняем продажу
		if kCandleCurrent > 80 &&
			dCandleCurrent > 80 &&
			kCandlePrev >= dCandlePrev &&
			kCandleCurrent < dCandleCurrent {

			// продажа валюты
			if len(stopLossOrders) != 0 {

				// отмена открытых STOP-LOSS ордеров
				for index, order := range stopLossOrders {
					_, err := client.CancelOrder(order.Pair, order.ID)
					if err != nil {
						log.Fatalln(err)
					}
					// удаление из массива STOP-LOSS ордеров
					copy(stopLossOrders[index:], stopLossOrders[index+1:])
					stopLossOrders = stopLossOrders[:len(stopLossOrders)-1]
				}

				// получение доступного свободного баланса для валюты
				balanceFree, err := client.GetBalanceFree("BTC")
				if err != nil {
					log.Fatalln(err)
				}

				// открытие ордера на продажу
				order, err := client.CreateMarketCellOrder("BTCUSDT", balanceFree, 6)
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
			if len(stopLossOrders) == 0 {
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
				order, err := client.CreateMarketBuyOrder("BTCUSDT", (balanceFree/price)*0.1, 6)
				if err != nil {
					log.Fatalln(err)
				}
				ordersWithoutStopLoss = append(ordersWithoutStopLoss, Order{Pair: order.Symbol, ID: order.OrderID})

				log.Println("Открыт ордер на покупку по цене", order.Price)
			}
		}

		time.Sleep(time.Second)
	}

}
