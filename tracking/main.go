package tracking

import (
	"log"
	"strconv"
	"sync"
	"time"

	bnc "../binance"
)

// DirectionTracking - поток отслеживания направления
func DirectionTracking(direction Direction, client *bnc.API, waitGroupDirectionTracking *sync.WaitGroup) {
	defer waitGroupDirectionTracking.Done()

	log.Println("Запущено отслеживание по направлению", direction.Base+direction.Quote, "с периодом",
		direction.Interval)

	// инициализация ордеров
	bot := Bot{Status: BotStatusWaitPurchase}

	// инициализация действия
	action := make(chan IndicatorsStatus)

	// запуск отслеживания индикаторами необходимого действия
	go trackIndicators(direction.Base+direction.Quote, direction.Interval, client, action)

	// определение необходимого действия
	for {
		switch <-action {
		case IndicatorsStatusBuy:
			{
				if bot.Status == BotStatusWaitPurchase {
					bot.Status = BotStatusActivePurchase

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
					buyOrder, err := client.CreateLimitBuyOrder(direction.Base+direction.Quote,
						balance/currentPrice*(direction.PercentOfBudgetPerTransaction/100),
						currentPrice+currentPrice*0.0005)
					if err != nil {
						log.Println(err)
						continue
					}

					finalBuyOrder := getFinalOrder(buyOrder.Symbol, buyOrder.OrderID, client)

					if OrderStatus(finalBuyOrder.Status) == OrderStatusFilled {
						// установка STOP-LOSS ордера
						quantity, _ := strconv.ParseFloat(finalBuyOrder.OrigQuantity, 64)
						price, _ := strconv.ParseFloat(finalBuyOrder.Price, 64)

						stopLossOrder, err := client.CreateStopLimitSellOrder(finalBuyOrder.Symbol, quantity,
							price-price*0.003, price-price*0.0031)
						if err != nil {
							log.Println(err)
							continue
						}

						go trackStopLossOrder(stopLossOrder.Symbol, stopLossOrder.OrderID, &bot.Status, client)

						log.Println("Создан ордер", buyOrder.OrderID, "на покупку с направлением",
							buyOrder.Symbol, "периодом", direction.Interval, "по цене",
							buyOrder.Price, "и количеством", buyOrder.OrigQuantity)

						bot.BuyOrderId = buyOrder.OrderID
						bot.StopLossOrderId = stopLossOrder.OrderID
						bot.Status = BotStatusWaitSell

						continue
					} else if OrderStatus(finalBuyOrder.Status) == OrderStatusExpired {
						bot.Status = BotStatusWaitPurchase
						continue
					}
				}
			}
		case IndicatorsStatusSell:
			{
				if bot.Status == BotStatusWaitSell {
					bot.Status = BotStatusActiveSell

					// отмена текущего STOP-LOSS ордера
					_, err := client.CancelOrder(direction.Base+direction.Quote, bot.StopLossOrderId)
					if err != nil {
						log.Println(err)
						continue
					}

					// получение параметров ордера на покупку
					buyOrder, err := client.GetOrder(direction.Base+direction.Quote, bot.BuyOrderId)

					purchasePrice, err := strconv.ParseFloat(buyOrder.Price, 64)
					if err != nil {
						log.Println(err)
						continue
					}

					// создание ордера на продажу
					currentPrice, err := client.GetCurrentPrice(direction.Base + direction.Quote)
					if err != nil {
						log.Println(err)
						continue
					}
					quantity, err := strconv.ParseFloat(buyOrder.ExecutedQuantity, 64)
					if err != nil {
						log.Println(err)
						continue
					}

					sellOrder, err := client.CreateLimitSellOrder(buyOrder.Symbol, quantity,
						currentPrice)
					if err != nil {
						log.Println(err)
						continue
					}

					finalSellOrder := getFinalOrder(sellOrder.Symbol, sellOrder.OrderID, client)

					if OrderStatus(finalSellOrder.Status) == OrderStatusFilled {
						log.Println("Создан ордер", finalSellOrder.OrderID, "на продажу с направлением",
							finalSellOrder.Symbol, "периодом", direction.Interval, "по цене",
							finalSellOrder.Price, "и количеством", finalSellOrder.OrigQuantity,
							"выгода составит", currentPrice*quantity-purchasePrice*quantity, direction.Quote)

						bot.Status = BotStatusWaitPurchase
						continue
					} else if OrderStatus(finalSellOrder.Status) == OrderStatusExpired {
						// повторная установка STOP-LOSS ордера
						quantity, _ := strconv.ParseFloat(buyOrder.OrigQuantity, 64)
						price, _ := strconv.ParseFloat(buyOrder.Price, 64)

						stopLossOrder, err := client.CreateStopLimitSellOrder(finalSellOrder.Symbol, quantity,
							price-price*0.0045, price-price*0.005)
						if err != nil {
							log.Println(err)
							continue
						}

						bot.Status = BotStatusWaitSell
						go trackStopLossOrder(stopLossOrder.Symbol, stopLossOrder.OrderID, &bot.Status, client)
						continue
					}

				}
			}
		}
		time.Sleep(time.Second / 10)
	}
}
