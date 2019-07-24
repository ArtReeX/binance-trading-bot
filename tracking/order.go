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

	log.Println("Создан ордер на покупку с направлением", order.Symbol,
		"по цене", order.Price, "и количеством", order.OrigQuantity)

	return order.OrderID, nil
}

// createSellOrder - функция создания ордера на продажу
func createSellOrder(base string,
	quote string,
	quantity float64,
	accuracyQuantity uint8,
	accuracyPrice uint8,
	client *bnc.API) (int64, error) {

	// получение текущей цены валюты
	currentPrice, err := client.GetCurrentPrice(base + quote)
	if err != nil {
		return 0, err
	}

	// создание ордера на продажу
	order, err := client.CreateLimitSellOrder(base+quote,
		quantity,
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
	orderInfo *OrderInfo,
	accuracyQuantity uint8,
	accuracyPrice uint8,
	client *bnc.API,
	waitGroupCreateStopLossOrders *sync.WaitGroup) {
	defer waitGroupCreateStopLossOrders.Done()

	for {
		// проверка статуса исполненности ордера
		if orderInfo.NotBackedID != -1 {
			// получение ордера
			order, err := client.GetOrder(pair, orderInfo.NotBackedID)
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

				// удаление индикатора из списка наблюдения и добавление идентификатора STOP-LOSS ордера
				orderInfo.Lock()
				orderInfo.NotBackedID = -1
				orderInfo.StopLossID = stopLossOrder.OrderID
				orderInfo.Quantity = quantity
				orderInfo.Price = price
				orderInfo.Unlock()

				log.Println("Удалён выполненный ордер из списка наблюдения с направлением",
					order.Symbol, "по цене", order.Price, "и количеством", order.ExecutedQuantity)
				log.Println("Добавлен STOP-LOSS ордер в список наблюдения с направлением",
					stopLossOrder.Symbol, "по цене", stopLossOrder.Price, "и количеством", stopLossOrder.OrigQuantity)

			} else if order.Status == "EXPIRED" || order.Status == "REJECTED" || order.Status == "CANCELED" {
				// удаление индикатора из списка наблюдения
				orderInfo.NotBackedID = -1

				log.Println("Удалён не выполненный ордер из списка наблюдения с направлением",
					order.Symbol, "по цене", order.Price, "и количеством", order.OrigQuantity)
			}
		}
		time.Sleep(time.Second * 1)
	}
}

// checkStopLossOrders - функция проверяет статус STOP-LOSS ордера
func checkStopLossOrders(pair string, orderInfo *OrderInfo, client *bnc.API, waitGroupCheckStopLossOrders *sync.WaitGroup) {
	defer waitGroupCheckStopLossOrders.Done()
	for {
		// проверка статуса исполненности ордера
		if orderInfo.StopLossID != -1 {
			// получение ордера
			order, err := client.GetOrder(pair, orderInfo.StopLossID)
			if err != nil {
				log.Println(err)
				continue
			}

			// удаление текущего отсутствующего ордера
			if order.Status != "NEW" && order.Status != "PARTIALLY_FILLED" {
				orderInfo.Lock()
				orderInfo.StopLossID = -1
				orderInfo.Unlock()

				log.Println("Удалён отсутствующий STOP-LOSS ордер из списка наблюдения с направлением",
					order.Symbol, "по цене", order.Price, "и количеством", order.OrigQuantity)
			}
		}

		time.Sleep(time.Second * 1)
	}
}

// cancelStopLossOrder - функция отмены STOP-LOSS ордера
func cancelStopLossOrder(pair string, id int64, client *bnc.API) error {
	cancelOrder, err := client.CancelOrder(pair, id)
	if err != nil {
		return err
	}
	// получение идентификатора оригинального ордера
	originalOrderID, err := strconv.ParseInt(cancelOrder.OrigClientOrderID, 10, 64)
	if err != nil {
		log.Fatalln(err)
	}
	// получение оригинального ордера
	originalOrder, err := client.GetOrder(cancelOrder.Symbol, originalOrderID)
	if err != nil {
		log.Println(err)
	}
	log.Println("Удалён не выполненный ордер из списка наблюдения с направлением",
		originalOrder.Symbol, "по цене", originalOrder.Price, "и количеством", originalOrder.OrigQuantity)

	return nil
}
