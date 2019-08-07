package binance

import (
	"github.com/adshao/go-binance"
	"strconv"
)

func formatOrder(order binance.Order) Order {
	priceValue, _ := strconv.ParseFloat(order.Price, 64)
	origQuantityValue, _ := strconv.ParseFloat(order.OrigQuantity, 64)
	executedQuantityValue, _ := strconv.ParseFloat(order.ExecutedQuantity, 64)
	cummulativeQuoteQuantityValue, _ := strconv.ParseFloat(order.CummulativeQuoteQuantity, 64)
	stopPriceValue, _ := strconv.ParseFloat(order.StopPrice, 64)

	return Order{
		Symbol:                   order.Symbol,
		OrderId:                  uint64(order.OrderID),
		Price:                    priceValue,
		OrigQuantity:             origQuantityValue,
		ExecutedQuantity:         executedQuantityValue,
		CummulativeQuoteQuantity: cummulativeQuoteQuantityValue,
		Status:                   OrderStatus(order.Status),
		StopPrice:                stopPriceValue,
	}
}
