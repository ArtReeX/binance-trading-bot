package tracking

import (
	"log"
	"strconv"
	"time"

	"github.com/adshao/go-binance"

	bnc "../binance"
)

// checkBuyOrder -  проверяет выполненность существующего ордера на покупку и создаёт для них STOP-LOSS ордер
func checkBuyOrder(buyOrder *binance.Order, stopLossOrder *binance.Order, client *bnc.API) {
	for {
		// проверка статуса исполненности ордера
		if buyOrder != nil &&
			((stopLossOrder == nil && buyOrder.Status == "FILLED") ||
				(stopLossOrder != nil && buyOrder != nil && stopLossOrder.Status != "FILLED")) {
			// получение количества исполнения ордера
			quantity, err := strconv.ParseFloat(buyOrder.ExecutedQuantity, 64)
			if err != nil {
				log.Fatalln(err)
			}
			// получение цены исполнения ордера
			price, err := strconv.ParseFloat(buyOrder.Price, 64)
			if err != nil {
				log.Fatalln(err)
			}

			// создание STOP-LOSS ордера
			createdOrder, err := client.CreateStopLimitSellOrder(buyOrder.Symbol,
				quantity,
				price-(price*0.0095),
				price-(price*0.01))
			if err != nil {
				log.Println(err)
				continue
			}

			// добавление идентификатора STOP-LOSS ордера
			stopLossOrder = &binance.Order{
				Symbol:           createdOrder.Symbol,
				OrderID:          createdOrder.OrderID,
				ClientOrderID:    createdOrder.ClientOrderID,
				Price:            createdOrder.Price,
				OrigQuantity:     createdOrder.OrigQuantity,
				ExecutedQuantity: createdOrder.ExecutedQuantity,
				Status:           createdOrder.Status,
				TimeInForce:      createdOrder.TimeInForce,
				Type:             createdOrder.Type,
				Side:             createdOrder.Side,
				Time:             createdOrder.TransactTime}

			log.Println("Добавлен STOP-LOSS ордер с направлением",
				stopLossOrder.Symbol, "по цене", stopLossOrder.Price, "и количеством", stopLossOrder.OrigQuantity)
		}
		time.Sleep(time.Second * 1)
	}
}

// updateOrderStatus - функция обновляет статус ордера
func updateOrderStatus(renewableOrders []*binance.Order, client *bnc.API) {
	for {
		for _, renewableOrder := range renewableOrders {
			if renewableOrder != nil {
				order, err := client.GetOrder(renewableOrder.Symbol, renewableOrder.OrderID)
				if err != nil {
					log.Println(err)
					continue
				}

				// удаляем если ордер отменён либо обновляем в противном случае
				if order.Status == "CANCELED" || order.Status == "EXPIRED" {
					renewableOrder = nil
				} else if order.Status != renewableOrder.Status {
					renewableOrder = order
					log.Println("Обновлен статус ордера из списка наблюдения с направлением",
						order.Symbol, "по цене", order.Price, "и количеством", order.OrigQuantity)
				}

			}
		}
		time.Sleep(time.Second * 1)
	}
}
