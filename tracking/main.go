package tracking

import (
	"log"
	"strconv"
	"sync"

	"github.com/adshao/go-binance"

	bnc "../binance"
)

// Direction - структура направления
type Direction struct {
	Base     string
	Quote    string
	Interval string
}

// OrderInfo - параметры ордера для текущего направления
type OrderInfo struct {
	BuyOrder      *binance.Order
	StopLossOrder *binance.Order
}

// DirectionTracking - поток отслеживания направления
func DirectionTracking(direction Direction,
	client *bnc.API,
	waitGroupDirectionTracking *sync.WaitGroup) {
	defer waitGroupDirectionTracking.Done()

	// создание параметров ордера для текущего направления
	orderInfo := &OrderInfo{}

	// запуск отслеживания индикаторами
	action := make(chan binance.SideType)
	go trackStochRSI(direction.Base+direction.Quote, direction.Interval, action, client)
	log.Println("Запущено отслеживание по направлению", direction.Base+direction.Quote, "с периодом", direction.Interval)

	// определение необходимого действия
	for {
		switch <-action {
		case binance.SideTypeBuy:
			{
				// создание ордера на покупку
				if orderInfo.BuyOrder == nil && orderInfo.StopLossOrder == nil {
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
					openOrder, err := client.CreateLimitBuyOrder(direction.Base+direction.Quote, balance/currentPrice*0.1, currentPrice)
					if err != nil {
						log.Println(err)
						continue
					}

					orderInfo.BuyOrder = &binance.Order{
						Symbol:           openOrder.Symbol,
						OrderID:          openOrder.OrderID,
						ClientOrderID:    openOrder.ClientOrderID,
						Price:            openOrder.Price,
						OrigQuantity:     openOrder.OrigQuantity,
						ExecutedQuantity: openOrder.ExecutedQuantity,
						Status:           openOrder.Status,
						TimeInForce:      openOrder.TimeInForce,
						Type:             openOrder.Type,
						Side:             openOrder.Side,
						Time:             openOrder.TransactTime}

					// запуск мониторинга за статусом
					go updateOrderStatus(orderInfo.BuyOrder, client)

					log.Println("Создан ордер на покупку с направлением", orderInfo.BuyOrder.Symbol, ", периодом", direction.Interval,
						"по цене", orderInfo.BuyOrder.Price, "и количеством", orderInfo.BuyOrder.OrigQuantity)

					// создание STOP-LOSS ордера
					go createLinkStopLossOrder(orderInfo.BuyOrder, orderInfo.StopLossOrder, client)
				}
			}
		case binance.SideTypeSell:
			{
				// создание ордера на продажу
				if orderInfo.BuyOrder != nil && orderInfo.StopLossOrder != nil {
					currentPrice, err := client.GetCurrentPrice(orderInfo.BuyOrder.Symbol)
					if err != nil {
						log.Println(err)
						continue
					}
					purchasePrice, err := strconv.ParseFloat(orderInfo.BuyOrder.ExecutedQuantity, 64)
					if err != nil {
						log.Println(err)
						continue
					}

					// проверка условий продажи
					if currentPrice > purchasePrice {
						// отмена STOP-LOSS ордера
						_, err := client.CancelOrder(orderInfo.StopLossOrder.Symbol, orderInfo.StopLossOrder.OrderID)
						if err != nil {
							log.Println(err)
							continue
						}

						log.Println("Отменён STOP-LOSS ордер на продажу с направлением", orderInfo.BuyOrder.Symbol, ", периодом",
							direction.Interval, "по цене", orderInfo.BuyOrder.Price, "и количеством", orderInfo.BuyOrder.OrigQuantity)

						// создание ордера на продажу
						quantity, err := strconv.ParseFloat(orderInfo.BuyOrder.ExecutedQuantity, 64)
						if err != nil {
							log.Println(err)
							continue
						}
						_, err = client.CreateMarketSellOrder(orderInfo.BuyOrder.Symbol, quantity)
						if err != nil {
							log.Println(err)
							continue
						}

						orderInfo.BuyOrder = nil
						orderInfo.StopLossOrder = nil

						log.Println("Создан ордер на продажу с направлением", orderInfo.BuyOrder.Symbol, ", периодом",
							direction.Interval, "по цене", orderInfo.BuyOrder.Price, "и количеством",
							orderInfo.BuyOrder.OrigQuantity, ", выгода составит", currentPrice*quantity-purchasePrice*quantity, direction.Quote)
					}
				}
			}
		}
	}
}