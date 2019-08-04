package binance

type Order struct {
	Symbol                   string
	OrderId                  uint64
	Price                    float64
	OrigQuantity             float64
	ExecutedQuantity         float64
	CummulativeQuoteQuantity float64
	Status                   OrderStatus
	StopPrice                float64
	IsWorking                bool
}

type OrderStatus string
