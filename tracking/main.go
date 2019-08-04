package tracking

import (
	"log"
	"sync"

	bnc "../binance"
)

// DirectionTracking - поток отслеживания направления
func DirectionTracking(direction Direction, client *bnc.API, waitGroupDirectionTracking *sync.WaitGroup) {
	defer waitGroupDirectionTracking.Done()

	log.Println("Запущено отслеживание по направлению", direction.Base+direction.Quote, "с периодом",
		direction.Interval)

	// инициализация ордеров
	bot := Bot{Status: BotStatusWaitPurchase}

	// определение необходимого действия
	for {
		// получение истории валюты
		candleHistory, err := client.GetCandleHistory(direction.Base+direction.Quote, direction.Interval)
		if err != nil {
			log.Println(err)
			continue
		}
		candleHistoryFormated := FormatBinanceCandles(candleHistory)

		switch getIndicatorStatuses(candleHistoryFormated) {
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

					if balance >= direction.PriceForOneTransaction {
						currentPrice := candleHistoryFormated[len(candleHistoryFormated)-1].Close

						buyOrderId, err := client.CreateLimitBuyOrder(direction.Base+direction.Quote,
							direction.PriceForOneTransaction/currentPrice,
							currentPrice)
						if err != nil {
							log.Println(err)
							bot.Status = BotStatusWaitPurchase
							continue
						}

						finalBuyOrder, _ := client.GetFinalOrder(direction.Base+direction.Quote, buyOrderId)

						if finalBuyOrder.Status == bnc.OrderStatusFilled {
							log.Println("Создан ордер", finalBuyOrder.OrderId, "на покупку с направлением",
								finalBuyOrder.Symbol, "периодом", direction.Interval, "по цене",
								finalBuyOrder.Price, "и количеством", finalBuyOrder.OrigQuantity)

							// установка STOP-LOSS ордера
							stopLossOrderId, err := client.CreateStopLimitSellOrder(finalBuyOrder.Symbol,
								finalBuyOrder.OrigQuantity, finalBuyOrder.Price-finalBuyOrder.Price*0.002,
								finalBuyOrder.Price-finalBuyOrder.Price*0.0015)
							if err != nil {
								log.Println(err)
								bot.Status = BotStatusWaitPurchase
								continue
							}

							log.Println("Добавлен STOP-LOSS ордер", stopLossOrderId, "привязанный к ордеру",
								buyOrderId)

							bot.BuyOrderId = buyOrderId
							bot.StopLossOrderId = stopLossOrderId
							bot.Status = BotStatusWaitSell

							go trackStopLossOrder(finalBuyOrder.Symbol, &bot.StopLossOrderId, &bot.StopLossOrderIdMutex,
								&bot.Status, &bot.StatusMutex, client)

							continue
						} else if finalBuyOrder.Status == bnc.OrderStatusExpired {
							bot.Status = BotStatusWaitPurchase
							continue
						}
					}
				}
			}
		case IndicatorsStatusSell:
			{
				if bot.Status == BotStatusWaitSell {
					bot.Status = BotStatusActiveSell

					// отмена текущего STOP-LOSS ордера
					err := client.CancelOrder(direction.Base+direction.Quote, bot.StopLossOrderId)
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

					// создание ордера на продажу
					currentPrice := candleHistoryFormated[len(candleHistoryFormated)-1].Close

					sellOrderId, err := client.CreateLimitSellOrder(buyOrder.Symbol, buyOrder.ExecutedQuantity,
						currentPrice)
					if err != nil {
						log.Println(err)
						bot.Status = BotStatusWaitSell
						continue
					}

					finalSellOrder, _ := client.GetFinalOrder(direction.Base+direction.Quote, sellOrderId)

					if finalSellOrder.Status == bnc.OrderStatusFilled {
						log.Println("Создан ордер", finalSellOrder.OrderId, "на продажу с направлением",
							finalSellOrder.Symbol, "периодом", direction.Interval, "по цене",
							finalSellOrder.Price, "и количеством", finalSellOrder.OrigQuantity, "выгода составит",
							currentPrice*buyOrder.ExecutedQuantity-buyOrder.Price*buyOrder.ExecutedQuantity,
							direction.Quote)

						bot.Status = BotStatusWaitPurchase
						continue
					} else if finalSellOrder.Status == bnc.OrderStatusExpired {
						// повторная установка STOP-LOSS ордера
						stopLossOrderId, err := client.CreateStopLimitSellOrder(finalSellOrder.Symbol,
							buyOrder.OrigQuantity,
							buyOrder.Price-buyOrder.Price*0.002, buyOrder.Price-buyOrder.Price*0.0015)
						if err != nil {
							log.Println(err)
							bot.Status = BotStatusWaitSell
							continue
						}

						bot.StopLossOrderId = stopLossOrderId
						bot.Status = BotStatusWaitSell

						go trackStopLossOrder(buyOrder.Symbol, &bot.StopLossOrderId, &bot.StopLossOrderIdMutex,
							&bot.Status, &bot.StatusMutex, client)

						continue
					}
				}
			}
		default:
			{

			}
		}
	}
}
