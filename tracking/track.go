package tracking

import (
	"log"
	"sync"

	"github.com/adshao/go-binance"

	bnc "../binance"
)

// OrderInfo - параметры ордера для текущего направления
type OrderInfo struct {
	LastBuyOrder  *binance.Order
	StopLossOrder *binance.Order
}

// DirectionTracking - поток отслеживания направления
func DirectionTracking(base string,
	quote string,
	accuracyQuantity uint8,
	accuracyPrice uint8,
	interval string,
	client *bnc.API,
	waitGroupDirectionTracking *sync.WaitGroup) {

	defer waitGroupDirectionTracking.Done()

	// создание параметров ордера для текущего направления
	orderInfo := &OrderInfo{}

	// запуск слежения за откурытыми и STOP-LOSS ордерами
	waitGroupUpdateOrderStatus := new(sync.WaitGroup)
	waitGroupUpdateOrderStatus.Add(2)
	go updateOrderStatus(orderInfo.LastBuyOrder, client, waitGroupUpdateOrderStatus)
	go updateOrderStatus(orderInfo.StopLossOrder, client, waitGroupUpdateOrderStatus)

	// запуск отслеживания индикаторами
	action := make(chan binance.SideType)
	go trackStochRSI(base+quote, interval, action, client)
	log.Println("Запущено отслеживание по направлению", base+quote, "с периодом", interval)

	// определение необходимого действия
	for {
		switch <-action {
		case binance.SideTypeBuy:
			{
				if (orderInfo.StopLossOrder == nil && orderInfo.LastBuyOrder == nil) ||
					(orderInfo.LastBuyOrder != nil &&
						orderInfo.LastBuyOrder.Status == "EXPIRED") ||
					(orderInfo.StopLossOrder != nil &&
						orderInfo.StopLossOrder.Status == "FILLED") {
					// создание ордера на покупку
					openOrder, err := createBuyOrder(base,
						quote,
						accuracyQuantity,
						accuracyPrice,
						client)
					if err != nil {
						log.Println(err)
						continue
					}
					orderInfo.LastBuyOrder = openOrder
				}
			}
		case binance.SideTypeSell:
			{
				if orderInfo.StopLossOrder != nil && orderInfo.LastBuyOrder != nil &&
					orderInfo.StopLossOrder.Status != "FILLED" && orderInfo.StopLossOrder.Status != "CANCELED" {
					// получение текущей цены валюты
					currentPrice, err := client.GetCurrentPrice(orderInfo.LastBuyOrder.Price)
					if err != nil {
						log.Println(err)
						continue
					}

					// получение цены по которой покупалась валюты
					purchasePrice, err := client.GetCurrentPrice(orderInfo.LastBuyOrder.Price)
					if err != nil {
						log.Println(err)
						continue
					}

					// проверка условий продажи
					if currentPrice > purchasePrice {
						// отмена STOP-LOSS ордера
						err := cancelStopLossOrder(orderInfo.StopLossOrder, client)
						if err != nil {
							log.Println(err)
							continue
						}
						orderInfo.StopLossOrder = nil
						// создание ордера на продажу
						_, err = createSellOrder(orderInfo.LastBuyOrder,
							accuracyQuantity,
							accuracyPrice,
							client)
						if err != nil {
							log.Println(err)
							continue
						}
						orderInfo.LastBuyOrder = nil
					}
				}
			}
		}
	}
}
