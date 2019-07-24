package tracking

import (
	"log"
	"sync"

	bnc "../binance"
)

// StopLossOrder - структура STOP-LOSS ордера
type StopLossOrder struct {
	Pair string
	ID   int64
}

// TypeOrder -  тип оредера (ПОКУПКА/ПРОДАЖА)
type TypeOrder string

const (
	// TypeOrderBuy -  ордер на покупку
	TypeOrderBuy TypeOrder = "BUY"
	// TypeOrderSell - ордер на продажу
	TypeOrderSell TypeOrder = "SELL"
)

// DirectionTracking - поток отслеживания направления
func DirectionTracking(base string,
	quote string,
	accuracyQuantity uint8,
	accuracyPrice uint8,
	interval string,
	client *bnc.API,
	waitGroupDirectionTracking *sync.WaitGroup) {

	defer waitGroupDirectionTracking.Done()

	// запуск слежения за STOP-LOSS ордерами
	var stopLossOrdersID []int64
	waitGroupCheckStopLossOrders := new(sync.WaitGroup)
	waitGroupCheckStopLossOrders.Add(1)
	go checkStopLossOrders(base+quote, stopLossOrdersID, client, waitGroupCheckStopLossOrders)

	// запуск слежения за неподкреплёнными (без STOP-LOSS) ордерами
	var notBackedOrdersID []int64
	waitGroupCreateStopLossOrders := new(sync.WaitGroup)
	waitGroupCreateStopLossOrders.Add(1)
	go createStopLossOrders(base+quote,
		notBackedOrdersID,
		stopLossOrdersID,
		accuracyQuantity,
		accuracyPrice,
		client,
		waitGroupCreateStopLossOrders)

	// запуск отслеживания индикаторами
	action := make(chan TypeOrder)
	go trackStochRSI(base+quote, interval, action, client)
	log.Println("Запущено отслеживание по направлению", base, "/", quote, "/", interval)

	// определение необходимого действия
	for {
		switch <-action {
		case TypeOrderBuy:
			{
				if len(stopLossOrdersID) == 0 && len(notBackedOrdersID) == 0 {
					orderID, err := createBuyOrder(base,
						quote,
						accuracyQuantity,
						accuracyPrice,
						client)
					if err != nil {
						log.Println(err)
					} else {
						stopLossOrdersID = append(stopLossOrdersID, orderID)
						log.Println("Создан ордер по направлению", base+quote, "/", interval)
					}
				}
			}
		case TypeOrderSell:
			{
				if len(stopLossOrdersID) == 0 && len(notBackedOrdersID) == 0 {
					orderID, err := createBuyOrder(base,
						quote,
						accuracyQuantity,
						accuracyPrice,
						client)
					if err != nil {
						log.Println(err)
					} else {
						stopLossOrdersID = append(stopLossOrdersID, orderID)
						log.Println("Создан ордер по направлению", base+quote, "/", interval)
					}
				}
				log.Println("Продажа по направлению", base+quote, "/", interval)
			}
		}
	}
}
