package tracking

type Interval string

type IndicatorsStatus int8

type Order struct {
	Symbol           string
	OrderID          int64
	ClientOrderID    string
	Price            string
	OrigQuantity     string
	ExecutedQuantity string
	Status           OrderStatus
	TimeInForce      string
	Type             string
	Side             string
	StopPrice        string
	IcebergQuantity  string
	Time             int64
}

type OrderStatus string

type Direction struct {
	Base                          string
	Quote                         string
	Interval                      Interval
	PercentOfBudgetPerTransaction float64
}

type OrderInfo struct {
	BuyOrder      *Order
	StopLossOrder *Order
	SellOrder     *Order
}
