package tracking

import (
	"log"
	"strconv"
	"sync"
	"time"

	bnc "../binance"
)

// DirectionTracking - поток отслеживания направления
func DirectionTracking(direction Direction,
	client *bnc.API,
	waitGroupDirectionTracking *sync.WaitGroup) {
	defer waitGroupDirectionTracking.Done()

	log.Println("Запущено отслеживание по направлению", direction.Base+direction.Quote, "с периодом",
		direction.Interval)

	// инициализация ордеров
	orderInfo := OrderInfo{BuyOrder: &Order{Status: OrderStatusNoCreated},
		StopLossOrder: &Order{Status: OrderStatusNoCreated},
		SellOrder:     &Order{Status: OrderStatusNoCreated}}

	action := make(chan IndicatorsStatus)

	// запуск отслеживания направления индикаторами
	go trackIndicators(direction.Base+direction.Quote, direction.Interval, client, action)

	// определение необходимого действия
	for {
		switch <-action {
		case IndicatorsStatusBuy:
			{
				if orderInfo.BuyOrder.Status == OrderStatusNoCreated &&
					orderInfo.StopLossOrder.Status == OrderStatusNoCreated {
					balance, err := client.GetBalanceFree(direction.Quote)
					if err != nil {
						log.Println(err)
						continue
					}
					currentPrice, err := client.GetCurrentPrice(direction.Base + direction.Quote)
					if err != nil {
						log.Println(err)
						continue
					}
					openOrder, err := client.CreateLimitBuyOrder(direction.Base+direction.Quote,
						balance/currentPrice*(direction.PercentOfBudgetPerTransaction/100),
						currentPrice+currentPrice*0.0005)
					if err != nil {
						log.Println(err)
						continue
					}

					orderInfo.SellOrder.Status = OrderStatusNoCreated
					orderInfo.BuyOrder = &Order{
						Symbol:           openOrder.Symbol,
						OrderID:          openOrder.OrderID,
						ClientOrderID:    openOrder.ClientOrderID,
						Price:            openOrder.Price,
						OrigQuantity:     openOrder.OrigQuantity,
						ExecutedQuantity: openOrder.ExecutedQuantity,
						Status:           OrderStatus(openOrder.Status),
						TimeInForce:      openOrder.TimeInForce,
						Type:             openOrder.Type,
						Side:             openOrder.Side,
						Time:             openOrder.TransactTime,
					}

					// запуск мониторинга за статусом ордера
					go updateOrderStatus(&orderInfo.BuyOrder, client)

					log.Println("Создан ордер", orderInfo.BuyOrder.OrderID, "на покупку с направлением",
						orderInfo.BuyOrder.Symbol, "периодом", direction.Interval, "по цене",
						orderInfo.BuyOrder.Price, "и количеством", orderInfo.BuyOrder.OrigQuantity)

					// создание STOP-LOSS ордера
					go createLinkStopLossOrder(&orderInfo.BuyOrder, &orderInfo.StopLossOrder, &orderInfo.SellOrder,
						client)
				}
			}
		case IndicatorsStatusSell:
			{
				if orderInfo.BuyOrder.Status != OrderStatusNoCreated &&
					orderInfo.StopLossOrder.Status != OrderStatusNoCreated {
					// создание ордера на продажу
					currentPrice, err := client.GetCurrentPrice(orderInfo.BuyOrder.Symbol)
					if err != nil {
						log.Println(err)
						continue
					}
					purchasePrice, err := strconv.ParseFloat(orderInfo.BuyOrder.Price, 64)
					if err != nil {
						log.Println(err)
						continue
					}

					// проверка условий продажи
					if currentPrice+currentPrice*0.001 > purchasePrice {
						// отмена STOP-LOSS ордера
						_, err := client.CancelOrder(orderInfo.StopLossOrder.Symbol,
							orderInfo.StopLossOrder.OrderID)
						if err != nil {
							log.Println(err)
							continue
						}

						log.Println("Отменён STOP-LOSS ордер", orderInfo.StopLossOrder.OrderID,
							"на продажу c направлением", orderInfo.BuyOrder.Symbol, "периодом", direction.Interval,
							"по цене", orderInfo.BuyOrder.Price, "и количеством", orderInfo.BuyOrder.OrigQuantity)

						// создание ордера на продажу
						quantity, err := strconv.ParseFloat(orderInfo.BuyOrder.ExecutedQuantity, 64)
						if err != nil {
							log.Println(err)
							continue
						}
						order, err := client.CreateLimitSellOrder(orderInfo.BuyOrder.Symbol, quantity,
							currentPrice)
						if err != nil {
							log.Println(err)
							continue
						}

						orderInfo.SellOrder = &Order{
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

						log.Println("Создан ордер", orderInfo.SellOrder.OrderID, "на продажу с направлением",
							orderInfo.BuyOrder.Symbol, "периодом", direction.Interval, "по цене",
							orderInfo.SellOrder.Price, "и количеством", orderInfo.BuyOrder.OrigQuantity,
							"выгода составит", currentPrice*quantity-purchasePrice*quantity, direction.Quote)

						// запуск мониторинга за статусом ордера
						go updateOrderStatus(&orderInfo.SellOrder, client)

						orderInfo.BuyOrder.Status = OrderStatusNoCreated
						orderInfo.StopLossOrder.Status = OrderStatusNoCreated
					}
				}
			}
		}
		time.Sleep(time.Second / 10)
	}
}
