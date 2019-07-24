package tracking

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/adshao/go-binance"

	bnc "../binance"
)

// createBuyOrder - функция создания ордера на покупку
func createBuyOrder(base string,
	quote string,
	accuracyQuantity uint8,
	accuracyPrice uint8,
	client *bnc.API) (*binance.Order, error) {

	// получение доступного баланса для покупки валюты
	balance, err := client.GetBalanceFree(quote)
	if err != nil {
		return nil, err
	}

	// получение текущей цены валюты
	currentPrice, err := client.GetCurrentPrice(base + quote)
	if err != nil {
		return nil, err
	}

	// создание ордера на покупку
	openOrder, err := client.CreateLimitBuyOrder(base+quote,
		balance/currentPrice*0.10,
		currentPrice,
		accuracyQuantity,
		accuracyPrice)
	if err != nil {
		return nil, err
	}

	log.Println("Создан ордер на покупку с направлением", openOrder.Symbol,
		"по цене", openOrder.Price, "и количеством", openOrder.OrigQuantity)

	return &binance.Order{
		Symbol:           openOrder.Symbol,
		OrderID:          openOrder.OrderID,
		ClientOrderID:    openOrder.ClientOrderID,
		Price:            openOrder.Price,
		OrigQuantity:     openOrder.OrigQuantity,
		ExecutedQuantity: openOrder.ExecutedQuantity,
		Status:           openOrder.Status,
		TimeInForce:      openOrder.TimeInForce,
		Type:             openOrder.Type,
		Side:             openOrder.Side,
		Time:             openOrder.TransactTime}, nil
}

// createSellOrder - функция создания ордера на продажу
func createSellOrder(lastBuyOrder *binance.Order,
	accuracyQuantity uint8,
	accuracyPrice uint8,
	client *bnc.API) (*binance.Order, error) {

	// получение старой цены валюты
	purchasePrice, err := strconv.ParseFloat(lastBuyOrder.ExecutedQuantity, 64)
	if err != nil {
		return nil, err
	}

	// получение текущей цены валюты
	currentPrice, err := client.GetCurrentPrice(lastBuyOrder.Price)
	if err != nil {
		return nil, err
	}

	// получение количества
	quantity, err := strconv.ParseFloat(lastBuyOrder.ExecutedQuantity, 64)
	if err != nil {
		return nil, err
	}

	// создание ордера на продажу
	openOrder, err := client.CreateLimitSellOrder(lastBuyOrder.Symbol,
		quantity,
		currentPrice,
		accuracyQuantity,
		accuracyPrice)
	if err != nil {
		return nil, err
	}

	log.Println("Создан ордер на продажу с направлением", openOrder.Symbol,
		"по цене", openOrder.Price, "и количеством", openOrder.OrigQuantity,
		", выгода составит", currentPrice*quantity-purchasePrice*quantity, openOrder.Symbol)

	return &binance.Order{
		Symbol:           openOrder.Symbol,
		OrderID:          openOrder.OrderID,
		ClientOrderID:    openOrder.ClientOrderID,
		Price:            openOrder.Price,
		OrigQuantity:     openOrder.OrigQuantity,
		ExecutedQuantity: openOrder.ExecutedQuantity,
		Status:           openOrder.Status,
		TimeInForce:      openOrder.TimeInForce,
		Type:             openOrder.Type,
		Side:             openOrder.Side,
		Time:             openOrder.TransactTime}, nil
}

// checkLastBuyOrder -  проверяет выполненность существующих ордеров на покупку и создаёт для них STOP-LOSS ордер
func checkLastBuyOrder(pair string,
	orderInfo OrderInfo,
	accuracyQuantity uint8,
	accuracyPrice uint8,
	client *bnc.API,
	waitGroupCheckLastBuyOrder *sync.WaitGroup) {
	defer waitGroupCheckLastBuyOrder.Done()

	for {
		// проверка статуса исполненности ордера
		if orderInfo.LastBuyOrder.Status == "FILLED" || orderInfo.LastBuyOrder.Status == "PARTIALLY_FILLED" {
			// получение количества исполнения ордера
			quantity, err := strconv.ParseFloat(orderInfo.LastBuyOrder.ExecutedQuantity, 64)
			if err != nil {
				log.Fatalln(err)
			}
			// получение цены исполнения ордера
			price, err := strconv.ParseFloat(orderInfo.LastBuyOrder.Price, 64)
			if err != nil {
				log.Fatalln(err)
			}

			// создание STOP-LOSS ордера
			stopLossOrder, err := client.CreateStopLimitSellOrder(orderInfo.LastBuyOrder.Symbol,
				quantity,
				price-(price*0.0095),
				price-(price*0.01),
				accuracyQuantity,
				accuracyPrice)
			if err != nil {
				log.Println(err)
				continue
			}

			// добавление идентификатора STOP-LOSS ордера
			orderInfo.StopLossOrder = &binance.Order{
				Symbol:           stopLossOrder.Symbol,
				OrderID:          stopLossOrder.OrderID,
				ClientOrderID:    stopLossOrder.ClientOrderID,
				Price:            stopLossOrder.Price,
				OrigQuantity:     stopLossOrder.OrigQuantity,
				ExecutedQuantity: stopLossOrder.ExecutedQuantity,
				Status:           stopLossOrder.Status,
				TimeInForce:      stopLossOrder.TimeInForce,
				Type:             stopLossOrder.Type,
				Side:             stopLossOrder.Side,
				Time:             stopLossOrder.TransactTime}

			log.Println("Добавлен STOP-LOSS ордер с направлением",
				orderInfo.StopLossOrder.Symbol, "по цене", orderInfo.StopLossOrder.Price, "и количеством", orderInfo.StopLossOrder.OrigQuantity)

		}
		time.Sleep(time.Second * 1)
	}
}

// updateOrderStatus - функция обновляет статус ордера
func updateOrderStatus(renewableOrder *binance.Order, client *bnc.API, waitGroupUpdateOrderStatus *sync.WaitGroup) {
	defer waitGroupUpdateOrderStatus.Done()
	for {
		if renewableOrder != nil {
			// получение ордера
			order, err := client.GetOrder(renewableOrder.Symbol, renewableOrder.OrderID)
			if err != nil {
				log.Println(err)
				continue
			}
			if order.Status != renewableOrder.Status {
				renewableOrder = order

				log.Println("Обновлен статус ордера из списка наблюдения с направлением",
					order.Symbol, "по цене", order.Price, "и количеством", order.OrigQuantity)
			}
		}
		time.Sleep(time.Second * 1)
	}
}

// cancelStopLossOrder - функция отмены STOP-LOSS ордера
func cancelStopLossOrder(stopLossOrder *binance.Order, client *bnc.API) error {
	_, err := client.CancelOrder(stopLossOrder.Symbol, stopLossOrder.OrderID)
	if err != nil {
		return err
	}
	log.Println("Удалён не выполненный STOP-LOSS ордер из списка наблюдения с направлением",
		stopLossOrder.Symbol, "по цене", stopLossOrder.Price, "и количеством", stopLossOrder.OrigQuantity)
	return nil
}
