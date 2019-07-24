package tracking

import (
	"log"
	"sync"

	bnc "../binance"
)

// OrderInfo - параметры ордера для текущего направления
type OrderInfo struct {
	sync.Mutex
	NotBackedID int64
	StopLossID  int64
	Price       float64
	Quantity    float64
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

	// создание параметров ордера для текущего направления
	orderInfo := OrderInfo{Price: 0, Quantity: 0, NotBackedID: -1, StopLossID: -1}

	// запуск слежения за STOP-LOSS ордерами
	waitGroupCheckStopLossOrders := new(sync.WaitGroup)
	waitGroupCheckStopLossOrders.Add(1)
	go checkStopLossOrders(base+quote, &orderInfo, client, waitGroupCheckStopLossOrders)

	// запуск слежения за неподкреплёнными (без STOP-LOSS) ордерами
	waitGroupCreateStopLossOrders := new(sync.WaitGroup)
	waitGroupCreateStopLossOrders.Add(1)
	go createStopLossOrders(base+quote,
		&orderInfo,
		accuracyQuantity,
		accuracyPrice,
		client,
		waitGroupCreateStopLossOrders)

	// запуск отслеживания индикаторами
	action := make(chan TypeOrder)
	go trackStochRSI(base+quote, interval, action, client)
	log.Println("Запущено отслеживание по направлению", base+quote, "с периодом", interval)

	// определение необходимого действия
	for {
		switch <-action {
		case TypeOrderBuy:
			{
				if orderInfo.StopLossID == -1 && orderInfo.NotBackedID == -1 {
					// создание ордера на покупку
					orderID, err := createBuyOrder(base,
						quote,
						accuracyQuantity,
						accuracyPrice,
						client)
					if err != nil {
						log.Println(err)
						continue
					}
					orderInfo.NotBackedID = orderID
				}
			}
		case TypeOrderSell:
			{
				if orderInfo.StopLossID != -1 || orderInfo.NotBackedID != -1 {
					// отмена STOP-LOSS ордера
					err := cancelStopLossOrder(base+quote, orderInfo.StopLossID, client)
					if err != nil {
						log.Println(err)
						continue
					}
					// создание ордера на продажу
					_, err = createSellOrder(base,
						quote,
						orderInfo.Quantity,
						accuracyQuantity,
						accuracyPrice,
						client)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
}
