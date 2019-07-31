package tracking

import (
	"log"
	"strconv"
	"time"

	bnc "../binance"
)

// updateOrderStatus - функция обновления статуса ордера
func updateOrderStatus(renewableOrder **Order, client *bnc.API) {
	for {
		if (*renewableOrder).Status != OrderStatusNoCreated {
			order, err := client.GetOrder((*renewableOrder).Symbol, (*renewableOrder).OrderID)
			if err != nil {
				log.Println(err)
				continue
			}
			if OrderStatus(order.Status) == OrderStatusCanceled || OrderStatus(order.Status) == OrderStatusExpired {
				(*renewableOrder).Status = OrderStatusNoCreated
				log.Println("Убран ордер", order.OrderID, "из списка наблюдения с направлением",
					order.Symbol, "по цене", order.Price, "и количеством", order.OrigQuantity)
				return
			} else if OrderStatus(order.Status) != (*renewableOrder).Status {
				*renewableOrder = &Order{
					Symbol:           order.Symbol,
					OrderID:          order.OrderID,
					ClientOrderID:    order.ClientOrderID,
					Price:            order.Price,
					OrigQuantity:     order.OrigQuantity,
					ExecutedQuantity: order.ExecutedQuantity,
					Status:           OrderStatus(order.Status),
					TimeInForce:      order.TimeInForce,
					Type:             order.Type,
					Side:             order.Side,
					Time:             order.Time}

				log.Println("Обновлен статус ордера", order.OrderID, "в списке наблюдения с направлением",
					order.Symbol, "по цене", order.Price, "и количеством", order.OrigQuantity)
			}
		}
		time.Sleep(time.Second * 1)
	}
}

// createLinkStopLoss - функция создания связующего STOP-LOSS ордера
func createLinkStopLossOrder(buyOrder **Order, stopLossOrder **Order, sellOrder **Order, client *bnc.API) {
	for {
		if (*buyOrder).Status != OrderStatusFilled &&
			(*stopLossOrder).Status == OrderStatusNoCreated &&
			(*sellOrder).Status == OrderStatusNoCreated {
			// получение количества исполнения ордера
			quantity, err := strconv.ParseFloat((*buyOrder).ExecutedQuantity, 64)
			if err != nil {
				log.Println(err)
				continue
			}
			// получение цены исполнения ордера
			price, err := strconv.ParseFloat((*buyOrder).Price, 64)
			if err != nil {
				log.Println(err)
				continue
			}
			// создание STOP-LOSS ордера
			order, err := client.CreateStopLimitSellOrder((*buyOrder).Symbol,
				quantity,
				price-(price*0.004),
				price-(price*0.0038))
			if err != nil {
				log.Println(err)
				continue
			}

			*stopLossOrder = &Order{
				Symbol:           order.Symbol,
				OrderID:          order.OrderID,
				ClientOrderID:    order.ClientOrderID,
				Price:            order.Price,
				OrigQuantity:     order.OrigQuantity,
				ExecutedQuantity: order.ExecutedQuantity,
				Status:           OrderStatus(order.Status),
				TimeInForce:      order.TimeInForce,
				Type:             order.Type,
				Side:             order.Side,
				Time:             order.TransactTime}

			// запуск мониторинга за статусом ордера
			go updateOrderStatus(stopLossOrder, client)

			log.Println("Добавлен STOP-LOSS ордер", (*stopLossOrder).OrderID, "привязанный к ордеру",
				(*buyOrder).OrderID, "с направлением", (*stopLossOrder).Symbol, "по цене",
				(*stopLossOrder).Price, "и количеством", (*stopLossOrder).OrigQuantity)

		} else if (*buyOrder).Status != OrderStatusNoCreated &&
			(*stopLossOrder).Status == OrderStatusFilled {
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

			log.Println("Сработал STOP-LOSS ордер", (*stopLossOrder).OrderID, "привязанный к ордеру",
				(*buyOrder).OrderID, "с направлением", (*stopLossOrder).Symbol, "по цене",
				(*stopLossOrder).Price, "и количеством", (*stopLossOrder).OrigQuantity, "потеря составила",
				purchasePrice*quantity-sellPrice*quantity, (*buyOrder).Symbol)

			(*buyOrder).Status = OrderStatusNoCreated
			(*stopLossOrder).Status = OrderStatusNoCreated

			return
		} else if (*buyOrder).Status == OrderStatusNoCreated {
			return
		}

		time.Sleep(time.Second * 1)
	}
}
