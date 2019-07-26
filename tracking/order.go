package tracking

import (
	"log"
	"strconv"
	"time"

	"github.com/adshao/go-binance"

	bnc "../binance"
)

// updateOrderStatus - функция обновления статуса ордера
func updateOrderStatus(renewableOrder **binance.Order, client *bnc.API) {
	for {
		if *renewableOrder != nil {
			order, err := client.GetOrder((*renewableOrder).Symbol, (*renewableOrder).OrderID)
			if err != nil {
				log.Println(err)
				continue
			}
			// если ордер отменён - удаляем, в ином случае - обновляем
			if order.Status == "CANCELED" || order.Status == "EXPIRED" {
				*renewableOrder = nil

				log.Println("Убран ордер", order.OrderID, "из списка наблюдения с направлением",
					order.Symbol, "по цене", order.Price, "и количеством", order.OrigQuantity)
				return
			} else if *renewableOrder != nil && order.Status != (*renewableOrder).Status {
				*renewableOrder = order

				log.Println("Обновлен статус ордера", order.OrderID, "в списке наблюдения с направлением",
					order.Symbol, "по цене", order.Price, "и количеством", order.OrigQuantity)
			}
		} else {
			return
		}
		time.Sleep(time.Second * 1)
	}
}

// createLinkStopLoss - функция создания связующего STOP-LOSS ордера
func createLinkStopLossOrder(buyOrder **binance.Order, stopLossOrder **binance.Order, client *bnc.API) {
	for {
		if *buyOrder != nil && (*buyOrder).Status == "FILLED" && *stopLossOrder == nil {
			// получение количества исполнения ордера
			quantity, err := strconv.ParseFloat((*buyOrder).ExecutedQuantity, 64)
			if err != nil {
				log.Fatalln(err)
				continue
			}
			// получение цены исполнения ордера
			price, err := strconv.ParseFloat((*buyOrder).Price, 64)
			if err != nil {
				log.Fatalln(err)
				continue
			}
			// создание STOP-LOSS ордера
			order, err := client.CreateStopLimitSellOrder((*buyOrder).Symbol,
				quantity,
				price-(price*0.005),
				price-(price*0.0048))
			if err != nil {
				log.Println(err)
				continue
			}

			*stopLossOrder = &binance.Order{
				Symbol:           order.Symbol,
				OrderID:          order.OrderID,
				ClientOrderID:    order.ClientOrderID,
				Price:            order.Price,
				OrigQuantity:     order.OrigQuantity,
				ExecutedQuantity: order.ExecutedQuantity,
				Status:           order.Status,
				TimeInForce:      order.TimeInForce,
				Type:             order.Type,
				Side:             order.Side,
				Time:             order.TransactTime}

			// запуск мониторинга за ордером
			go updateOrderStatus(stopLossOrder, client)

			log.Println("Добавлен STOP-LOSS ордер", (*stopLossOrder).OrderID, "привязанный к ордеру", (*buyOrder).OrderID, "с направлением",
				(*stopLossOrder).Symbol, "по цене", (*stopLossOrder).Price, "и количеством", (*stopLossOrder).OrigQuantity)

		} else if *buyOrder != nil && *stopLossOrder != nil && (*stopLossOrder).Status == "FILLED" {
			purchasePrice, err := strconv.ParseFloat((*buyOrder).Price, 64)
			if err != nil {
				log.Println(err)
				continue
			}
			sellPrice, err := client.GetCurrentPrice((*buyOrder).Symbol)
			if err != nil {
				log.Println(err)
				continue
			}
			quantity, err := strconv.ParseFloat((*buyOrder).ExecutedQuantity, 64)
			if err != nil {
				log.Println(err)
				continue
			}

			log.Println("Сработал STOP-LOSS ордер", (*stopLossOrder).OrderID, "привязанный к ордеру", (*buyOrder).OrderID, "с направлением",
				(*stopLossOrder).Symbol, "по цене", (*stopLossOrder).Price, "и количеством", (*stopLossOrder).OrigQuantity,
				"потеря составила", purchasePrice*quantity-sellPrice*quantity, (*buyOrder).Symbol)

			*buyOrder = nil
			*stopLossOrder = nil

			return
		} else if *buyOrder == nil {
			return
		}
		time.Sleep(time.Second * 1)
	}
}
