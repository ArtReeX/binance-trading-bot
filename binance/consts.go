package binance

const (
	// OrderStatusNew - статус "новый"
	OrderStatusNew OrderStatus = "NEW"

	// OrderStatusPartiallyFilled - статус "частично выполнен"
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"

	// OrderStatusFilled - статус "выполнен"
	OrderStatusFilled OrderStatus = "FILLED"

	// OrderStatusCanceled - статус "отменён"
	OrderStatusCanceled OrderStatus = "CANCELED"

	// OrderStatusPendingCancel - статус "запрошена отмена"
	OrderStatusPendingCancel OrderStatus = "PENDING_CANCEL"

	// OrderStatusRejected - статус "отклонён"
	OrderStatusRejected OrderStatus = "REJECTED"

	// OrderStatusExpired - статус "истёкший"
	OrderStatusExpired OrderStatus = "EXPIRED"
)
