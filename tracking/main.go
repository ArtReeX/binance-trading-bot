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
	bot := Bot{Status: BotStatusWaitPurchase, newStatus: make(chan BotStatus), newStopLossOrderId: make(chan int64)}

	// инициализация канала для получения действия индикаторов
	action := make(chan IndicatorsStatus)

	// запуск отслеживания индикаторами необходимого действия
	go trackIndicators(direction.Base+direction.Quote, direction.Interval, client, action)

	// определение необходимого действия
	for {
		select {
		// инициализация канала для получения нового статуса бота
		case newBotStatus := <-bot.newStatus:
			{
				bot.Status = newBotStatus
			}
		// получение нового идентификатора STOP-LOSS ордера
		case newStopLossId := <-bot.newStopLossOrderId:
			{
				bot.StopLossOrderId = newStopLossId
			}
		// получение действия индикаторов
		case actionIndicators := <-action:
			{
				switch actionIndicators {
				case IndicatorsStatusBuy:
					{
						if bot.Status == BotStatusWaitPurchase {
							bot.Status = BotStatusActivePurchase

							balance, err := client.GetBalanceFree(direction.Quote)
							if err != nil {
								log.Println(err)
								bot.Status = BotStatusWaitPurchase
								continue
							}
							currentPrice, err := client.GetCurrentPrice(direction.Base + direction.Quote)
							if err != nil {
								log.Println(err)
								bot.Status = BotStatusWaitPurchase
								continue
							}
							buyOrder, err := client.CreateLimitBuyOrder(direction.Base+direction.Quote,
								balance/currentPrice*(direction.PercentOfBudgetPerTransaction/100),
								currentPrice+currentPrice*0.0005)
							if err != nil {
								log.Println(err)
								bot.Status = BotStatusWaitPurchase
								continue
							}

							finalBuyOrder := getFinalOrder(buyOrder.Symbol, buyOrder.OrderID, client)

							if OrderStatus(finalBuyOrder.Status) == OrderStatusFilled {
								log.Println("Создан ордер", finalBuyOrder.OrderID, "на покупку с направлением",
									finalBuyOrder.Symbol, "периодом", direction.Interval, "по цене",
									finalBuyOrder.Price, "и количеством", finalBuyOrder.OrigQuantity)

								// установка STOP-LOSS ордера
								quantity, err := strconv.ParseFloat(finalBuyOrder.OrigQuantity, 64)
								if err != nil {
									log.Println(err)
									bot.Status = BotStatusWaitPurchase
									continue
								}
								price, err := strconv.ParseFloat(finalBuyOrder.Price, 64)
								if err != nil {
									log.Println(err)
									bot.Status = BotStatusWaitPurchase
									continue
								}

								stopLossOrder, err := client.CreateStopLimitSellOrder(finalBuyOrder.Symbol, quantity,
									price-price*0.0015, price-price*0.00012)
								if err != nil {
									log.Println(err)
									bot.Status = BotStatusWaitPurchase
									continue
								}

								log.Println("Добавлен STOP-LOSS ордер", stopLossOrder.OrderID, "с направлением",
									stopLossOrder.Symbol, "по цене", stopLossOrder.Price, "и количеством",
									stopLossOrder.OrigQuantity)

								bot.BuyOrderId = buyOrder.OrderID
								bot.StopLossOrderId = stopLossOrder.OrderID
								bot.Status = BotStatusWaitSell

								go trackStopLossOrder(stopLossOrder.Symbol, &bot.StopLossOrderId, &bot.Status,
									bot.newStatus, bot.newStopLossOrderId, client)

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
								bot.Status = BotStatusWaitSell
								continue
							}

							// получение параметров ордера на покупку
							buyOrder, err := client.GetOrder(direction.Base+direction.Quote, bot.BuyOrderId)
							if err != nil {
								log.Println(err)
								bot.Status = BotStatusWaitSell
								continue
							}

							purchasePrice, err := strconv.ParseFloat(buyOrder.Price, 64)
							if err != nil {
								log.Println(err)
								bot.Status = BotStatusWaitSell
								continue
							}

							// создание ордера на продажу
							currentPrice, err := client.GetCurrentPrice(direction.Base + direction.Quote)
							if err != nil {
								log.Println(err)
								bot.Status = BotStatusWaitSell
								continue
							}
							quantity, err := strconv.ParseFloat(buyOrder.ExecutedQuantity, 64)
							if err != nil {
								log.Println(err)
								bot.Status = BotStatusWaitSell
								continue
							}

							sellOrder, err := client.CreateLimitSellOrder(buyOrder.Symbol, quantity,
								currentPrice)
							if err != nil {
								log.Println(err)
								bot.Status = BotStatusWaitSell
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
								quantity, err := strconv.ParseFloat(buyOrder.OrigQuantity, 64)
								if err != nil {
									log.Println(err)
									bot.Status = BotStatusWaitSell
									continue
								}
								price, err := strconv.ParseFloat(buyOrder.Price, 64)
								if err != nil {
									log.Println(err)
									bot.Status = BotStatusWaitSell
									continue
								}

								stopLossOrder, err := client.CreateStopLimitSellOrder(finalSellOrder.Symbol, quantity,
									price-price*0.0015, price-price*0.0012)
								if err != nil {
									log.Println(err)
									bot.Status = BotStatusWaitSell
									continue
								}

								bot.StopLossOrderId = stopLossOrder.OrderID
								bot.Status = BotStatusWaitSell

								go trackStopLossOrder(stopLossOrder.Symbol, &stopLossOrder.OrderID, &bot.Status,
									bot.newStatus, bot.newStopLossOrderId, client)

								continue
							}
						}
					}
				}
			}
		default:
			{
				time.Sleep(time.Second / 10)
			}
		}
	}
}
