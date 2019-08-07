package tracking

import (
	"log"
	"sync"
	"time"

	bnc "../binance"
)

// DirectionTracking - поток отслеживания направления
func DirectionTracking(pair string, interval string, priceForOneTransaction float64, fee float64, client *bnc.Api,
	waitGroupDirectionTracking *sync.WaitGroup) {
	defer waitGroupDirectionTracking.Done()

	log.Println("Запущено отслеживание по направлению", pair, "с периодом", interval)

	// инициализация ордеров
	bot := Bot{Status: BotStatusWaitPurchase}

	// определение необходимого действия
	for {
		/*
			// получение истории валюты
			candleHistory, err := client.GetCandleHistory(pair, interval)
			if err != nil {
				log.Println(err)
				continue
			}
		*/

		// получение глубины стакана валюты
		depth, err := client.GetDepth(pair, 50)
		if err != nil {
			log.Println(err)
			continue
		}

		// получение текущей цены
		currentPrice, err := client.GetCurrentPrice(pair)
		if err != nil {
			log.Println(err)
			continue
		}

		switch getIndicatorStatuses( /*candleHistory,*/ depth) {
		case IndicatorsStatusBuy:
			{
				if bot.Status == BotStatusWaitPurchase {
					bot.Status = BotStatusActivePurchase

					balance, err := client.GetBalanceFree(client.Pairs[pair].QuoteAsset)
					if err != nil {
						log.Println(err)
						bot.Status = BotStatusWaitPurchase
						continue
					}

					if balance >= priceForOneTransaction {
						// currentPrice := candleHistoryFormated[len(candleHistoryFormated)-1].Close

						buyOrderId, err := client.CreateLimitBuyOrder(pair,
							priceForOneTransaction/(currentPrice*1.0005), currentPrice*1.0005)
						if err != nil {
							log.Println(err)
							bot.Status = BotStatusWaitPurchase
							continue
						}

						finalBuyOrder, _ := client.GetFinalOrder(pair, buyOrderId)

						if finalBuyOrder.Status == bnc.OrderStatusFilled {
							log.Println("Создан ордер", finalBuyOrder.OrderId, "на покупку с направлением",
								finalBuyOrder.Symbol, "периодом", interval, "по цене",
								finalBuyOrder.Price, "и количеством", finalBuyOrder.OrigQuantity)

							// установка STOP-LOSS ордера
							stopLossOrderId, err := client.CreateStopLimitSellOrder(finalBuyOrder.Symbol,
								finalBuyOrder.OrigQuantity, finalBuyOrder.Price*0.998, finalBuyOrder.Price*0.9985)
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

					// получение параметров ордера на покупку
					buyOrder, err := client.GetOrder(pair, bot.BuyOrderId)
					if err != nil {
						log.Println(err)
						bot.Status = BotStatusWaitSell
						continue
					}

					// проверяем, окупает ли наша сделка комиссию
					if currentPrice > buyOrder.Price*(1+fee) {

						// отмена текущего STOP-LOSS ордера
						err := client.CancelOrder(pair, bot.StopLossOrderId)
						if err != nil {
							log.Println(err)
							bot.Status = BotStatusWaitSell
							continue
						}

						// создание ордера на продажу
						// currentPrice := candleHistoryFormated[len(candleHistoryFormated)-1].Close

						sellOrderId, err := client.CreateLimitSellOrder(buyOrder.Symbol, buyOrder.ExecutedQuantity,
							currentPrice*0.9995)
						if err != nil {
							log.Println(err)
							bot.Status = BotStatusWaitSell
							continue
						}

						finalSellOrder, _ := client.GetFinalOrder(pair, sellOrderId)

						if finalSellOrder.Status == bnc.OrderStatusFilled {
							log.Println("Создан ордер", finalSellOrder.OrderId, "на продажу с направлением",
								finalSellOrder.Symbol, "периодом", interval, "по цене",
								finalSellOrder.Price, "и количеством", finalSellOrder.OrigQuantity, "выгода составила",
								(currentPrice*buyOrder.ExecutedQuantity-
									buyOrder.Price*buyOrder.ExecutedQuantity)*(1-fee), client.Pairs[pair].QuoteAsset)

							bot.Status = BotStatusWaitPurchase
							continue
						} else if finalSellOrder.Status == bnc.OrderStatusExpired {
							// повторная установка STOP-LOSS ордера
							stopLossOrderId, err := client.CreateStopLimitSellOrder(finalSellOrder.Symbol,
								buyOrder.OrigQuantity, buyOrder.Price*0.998, buyOrder.Price*0.9985)
							if err != nil {
								log.Println(err)
								bot.Status = BotStatusWaitSell
								continue
							}

							bot.StopLossOrderId = stopLossOrderId
							bot.Status = BotStatusWaitSell

							continue
						}
					}
				}
			}
		}

		// проверка STOP-LOSS ордера
		if bot.Status == BotStatusWaitSell {
			order, err := client.GetOrder(pair, bot.StopLossOrderId)
			if err != nil {
				log.Println(err)
				continue
			}

			if order.Status == bnc.OrderStatusCanceled {
				// в случае, если STOP-LOSS ордер был отменён вручную то устанавливаем его снова
				createdOrderId, err := client.CreateStopLimitSellOrder(order.Symbol, order.OrigQuantity, order.Price,
					order.StopPrice)
				if err != nil {
					log.Println(err)
					continue
				}

				log.Println("Добавлен недостающий STOP-LOSS ордер", createdOrderId, "взамен ордера", order.OrderId)
				bot.StopLossOrderId = createdOrderId
			} else if order.Status == bnc.OrderStatusFilled {
				// если STOP-LOSS ордер сработал переводим бота в режим покупки
				log.Println("Сработал STOP-LOSS ордер", order.OrderId)
				bot.Status = BotStatusWaitPurchase
			}
		}
		time.Sleep(time.Second)
	}
}
