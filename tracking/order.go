package tracking

import (
	bnc "../binance"
	"log"
	"sync"
)

func trackStopLossOrder(pair string, id *uint64, idMutex *sync.Mutex, status *BotStatus, statusMutex *sync.Mutex,
	client *bnc.API) {
	for {
		order, err := client.GetOrder(pair, *id)
		if err != nil {
			log.Println(err)
			continue
		}

		if order.Status == bnc.OrderStatusCanceled &&
			*status != BotStatusActiveSell &&
			*status != BotStatusWaitPurchase {

			// в случае, если STOP-LOSS ордер был отменён вручную то устанавливаем его снова
			createdOrderId, err := client.CreateStopLimitSellOrder(order.Symbol, order.OrigQuantity, order.Price,
				order.StopPrice)
			if err != nil {
				log.Println(err)
				continue
			}

			log.Println("Добавлен недостающий STOP-LOSS ордер", createdOrderId, "взамен ордера",
				order.OrderId, "с направлением", order.Symbol, "по цене",
				order.Price, "и количеством", order.OrigQuantity)

			idMutex.Lock()
			*id = createdOrderId
			idMutex.Unlock()
		} else if order.Status == bnc.OrderStatusFilled {
			// если STOP-LOSS ордер сработал переводим бота в режим покупки
			log.Println("Сработал STOP-LOSS ордер", order.OrderId, "с направлением", order.Symbol, "по цене",
				order.Price, "и количеством", order.OrigQuantity)

			statusMutex.Lock()
			*status = BotStatusWaitPurchase
			statusMutex.Unlock()
			return
		} else if order.Status == bnc.OrderStatusCanceled && *status == BotStatusWaitPurchase {
			// если продажа была выполнена перестаём отслеживать
			return
		}
	}
}
