package tracking

const (
	IndicatorsStatusBuy     IndicatorsStatus = 1
	IndicatorsStatusNeutral IndicatorsStatus = 0
	IndicatorsStatusSell    IndicatorsStatus = -1
)

const (
	OrderStatusNoCreated       OrderStatus = "NO_CREATED"
	OrderStatusNew             OrderStatus = "NEW"
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"
	OrderStatusFilled          OrderStatus = "FILLED"
	OrderStatusCanceled        OrderStatus = "CANCELED"
	OrderStatusPendingCancel   OrderStatus = "PENDING_CANCEL"
	OrderStatusRejected                    = "REJECTED"
	OrderStatusExpired         OrderStatus = "EXPIRED"
)

const (
	IntervalOneMinute      Interval = "1m"
	IntervalFiveMinutes    Interval = "5m"
	IntervalFifteenMinutes Interval = "15m"
	IntervalThirtyMinutes  Interval = "30m"
	IntervalOneHour        Interval = "1h"
	IntervalTwoHours       Interval = "2h"
	IntervalFourHours      Interval = "4h"
	IntervalSixHours       Interval = "6h"
	IntervalTwelveHours    Interval = "12h"
	IntervalOneDay         Interval = "1d"
	IntervalOneWeek        Interval = "1w"
	IntervalOneMonth       Interval = "1M"
)
