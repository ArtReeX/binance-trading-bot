package tracking

type Interval string

type IndicatorsStatus int8

type BotStatus uint8

type OrderStatus string

type Direction struct {
	Base                          string
	Quote                         string
	Interval                      Interval
	PercentOfBudgetPerTransaction float64
}

type Bot struct {
	BuyOrderId      int64
	StopLossOrderId int64
	SellOrderId     int64

	Status BotStatus
}
