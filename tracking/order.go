package tracking

import (
	"github.com/adshao/go-binance"
	"log"
	"strconv"
	"time"

	bnc "../binance"
)

func getFinalOrder(pair string, id int64, client *bnc.API) *binance.Order {
	for {
		order, err := client.GetOrder(pair, id)
		if err != nil {
			log.Println(err)
			continue
		}
		if OrderStatus(order.Status) != OrderStatusNew && OrderStatus(order.Status) != OrderStatusPartiallyFilled {
			return order
		}
		time.Sleep(time.Second / 5)
	}
}

func trackStopLossOrder(pair string, id *int64, status *BotStatus, newStatus chan BotStatus,
	newId chan<- int64, client *bnc.API) {
	for {
		order, err := client.GetOrder(pair, *id)
		if err != nil {
			log.Println(err)
			continue
		}

		if OrderStatus(order.Status) == OrderStatusCanceled &&
			*status != BotStatusActiveSell &&
			*status != BotStatusWaitPurchase {

			// в случае, если STOP-LOSS ордер был отменён вручную то устанавливаем его снова
			quantity, _ := strconv.ParseFloat(order.OrigQuantity, 64)
			if err != nil {
				log.Println(err)
				continue
			}
			price, _ := strconv.ParseFloat(order.Price, 64)
			if err != nil {
				log.Println(err)
				continue
			}
			stopPrice, _ := strconv.ParseFloat(order.StopPrice, 64)
			if err != nil {
				log.Println(err)
				continue
			}

			createdOrder, err := client.CreateStopLimitSellOrder(order.Symbol, quantity, price, stopPrice)
			if err != nil {
				log.Println(err)
				continue
			}

			log.Println("Добавлен недостающий STOP-LOSS ордер", createdOrder.OrderID, "взамен ордера",
				order.OrderID, "с направлением", createdOrder.Symbol, "по цене",
				createdOrder.Price, "и количеством", createdOrder.OrigQuantity)

			newId <- createdOrder.OrderID
		} else if OrderStatus(order.Status) == OrderStatusFilled {
			// если STOP-LOSS ордер сработал переводим бота в режим покупки
			log.Println("Сработал STOP-LOSS ордер", order.OrderID, "с направлением", order.Symbol, "по цене",
				order.Price, "и количеством", order.OrigQuantity)

			newStatus <- BotStatusWaitPurchase
			return
		} else if OrderStatus(order.Status) == OrderStatusCanceled && *status == BotStatusWaitPurchase {
			// если продажа была выполнена перестаём отслеживать
			return
		} else if OrderStatus(order.Status) != OrderStatusCanceled && *status == BotStatusWaitSell {
			// передвижение STOP-LOSS ордера в точку безубытка
			currentPrice, err := client.GetCurrentPrice(pair)
			if err != nil {
				log.Println(err)
				continue
			}

			stopLossPrice, _ := strconv.ParseFloat(order.Price, 64)
			if err != nil {
				log.Println(err)
				continue
			}

			if currentPrice > stopLossPrice+stopLossPrice*0.0015 {
				quantity, _ := strconv.ParseFloat(order.OrigQuantity, 64)
				if err != nil {
					log.Println(err)
					continue
				}

				// отмена текущего STOP-LOSS ордера
				_, err := client.CancelOrder(order.Symbol, order.OrderID)
				if err != nil {
					log.Println(err)
					continue
				}

				createdOrder, err := client.CreateStopLimitSellOrder(order.Symbol, quantity,
					currentPrice-currentPrice*0.0015,
					currentPrice-currentPrice*0.0012)
				if err != nil {
					log.Println(err)
					continue
				}

				log.Println("Передвинут STOP-LOSS ордер", createdOrder.OrderID, "взамен ордера",
					order.OrderID, "с направлением", createdOrder.Symbol, "по цене",
					createdOrder.Price, "и количеством", createdOrder.OrigQuantity)

				newId <- createdOrder.OrderID
			}
		}
	}
}
