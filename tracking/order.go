package tracking

import (
	"log"
	"strconv"
	"sync"
	"time"

	bnc "../binance"
)

// createBuyOrder - функция создания ордера на покупку
func createBuyOrder(base string,
	quote string,
	accuracyQuantity uint8,
	accuracyPrice uint8,
	client *bnc.API) (int64, error) {

	// получение доступного баланса для покупки валюты
	balance, err := client.GetBalanceFree(quote)
	if err != nil {
		return 0, err
	}

	// получение текущей цены валюты
	currentPrice, err := client.GetCurrentPrice(base + quote)
	if err != nil {
		return 0, err
	}

	// создание ордера на покупку
	order, err := client.CreateLimitBuyOrder(base+quote,
		balance/currentPrice*0.10,
		currentPrice,
		accuracyQuantity,
		accuracyPrice)
	if err != nil {
		return 0, err
	}

	log.Println("Создан ордер на продажу с направлением", order.Symbol,
		"по цене", order.Price, "и количеством", order.OrigQuantity)

	return order.OrderID, nil
}

// createSellOrder - функция создания ордера на продажу
func createSellOrder(base string,
	quote string,
	accuracyQuantity uint8,
	accuracyPrice uint8,
	client *bnc.API) (int64, error) {

	// получение доступного баланса для покупки валюты
	balance, err := client.GetBalanceFree(quote)
	if err != nil {
		return 0, err
	}

	// получение текущей цены валюты
	currentPrice, err := client.GetCurrentPrice(base + quote)
	if err != nil {
		return 0, err
	}

	// создание ордера на покупку
	order, err := client.CreateLimitSellOrder(base+quote,
		balance/currentPrice*0.10,
		currentPrice,
		accuracyQuantity,
		accuracyPrice)
	if err != nil {
		return 0, err
	}

	log.Println("Создан ордер на продажу с направлением", order.Symbol,
		"по цене", order.Price, "и количеством", order.OrigQuantity)

	return order.OrderID, nil
}

// createStopLossOrders -  проверяет выполненность существующих ордеров на покупку и создаёт для них STOP-LOSS ордер
func createStopLossOrders(pair string,
	notBackedOrdersID []int64,
	stopLossOrdersID []int64,
	accuracyQuantity uint8,
	accuracyPrice uint8,
	client *bnc.API,
	waitGroupCreateStopLossOrders *sync.WaitGroup) {
	defer waitGroupCreateStopLossOrders.Done()

	for {
		// проверка статуса исполнености ордера
		for index, id := range notBackedOrdersID {
			// получение ордера
			order, err := client.GetOrder(pair, id)
			if err != nil {
				log.Println(err)
				continue
			}

			if order.Status == "FILLED" || order.Status == "PARTIALLY_FILLED" {
				// получение количества исполнения ордера
				quantity, err := strconv.ParseFloat(order.ExecutedQuantity, 64)
				if err != nil {
					log.Fatalln(err)
				}
				// получение цены исполнения ордера
				price, err := strconv.ParseFloat(order.Price, 64)
				if err != nil {
					log.Fatalln(err)
				}

				// создание STOP-LOSS ордера
				stopLossOrder, err := client.CreateStopLimitSellOrder(order.Symbol,
					quantity,
					price-(price*0.01),
					price-(price*0.01),
					accuracyQuantity,
					accuracyPrice)
				if err != nil {
					log.Println(err)
					continue
				}

				// удаление индикатора из списка наблюдения
				copy(notBackedOrdersID[index:], notBackedOrdersID[index+1:])
				notBackedOrdersID = notBackedOrdersID[:len(notBackedOrdersID)-1]

				// добавление идентификатора STOP-LOSS ордера
				stopLossOrdersID = append(stopLossOrdersID, stopLossOrder.OrderID)

				log.Println("Удалён выполненный ордер из списка наблюдения с направлением",
					order.Symbol, "по цене", price, "и количеством", quantity)
				log.Println("Добавлен STOP-LOSS ордер в список наблюдения с направлением",
					order.Symbol, "по цене", order.Price, "и количеством", order.ExecutedQuantity)

			} else if order.Status == "EXPIRED" || order.Status == "REJECTED" || order.Status == "CANCELED" {
				// удаление индикатора из списка наблюдения
				copy(notBackedOrdersID[index:], notBackedOrdersID[index+1:])
				notBackedOrdersID = notBackedOrdersID[:len(notBackedOrdersID)-1]

				log.Println("Удалён не выполненный ордер из списка наблюдения с направлением",
					order.Symbol, "по цене", order.Price, "и количеством", order.OrigQuantity)
			}
			time.Sleep(time.Second * 1)
		}
	}
}

// checkStopLossOrders - функция проверяет статус STOP-LOSS ордеров
func checkStopLossOrders(pair string, stopLossOrdersID []int64, client *bnc.API, waitGroupCheckStopLossOrders *sync.WaitGroup) {
	defer waitGroupCheckStopLossOrders.Done()
	for {
		// проверка статуса исполнености ордера
		for index, id := range stopLossOrdersID {
			// получение ордера
			order, err := client.GetOrder(pair, id)
			if err != nil {
				log.Println(err)
				continue
			}

			// удаление текущего отсутствующего ордера
			if order.Status != "NEW" && order.Status != "PARTIALLY_FILLED" {
				copy(stopLossOrdersID[index:], stopLossOrdersID[index+1:])
				stopLossOrdersID = stopLossOrdersID[:len(stopLossOrdersID)-1]

				log.Println("Удалён отсутствующий STOP-LOSS ордер из списка наблюдения с направлением",
					order.Symbol, "по цене", order.Price, "и количеством", order.OrigQuantity)
			}
		}
		time.Sleep(time.Second * 1)
	}
}
