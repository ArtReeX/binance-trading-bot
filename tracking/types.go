package tracking

type (
	IndicatorsStatus int8

	BotStatus uint8

	Bot struct {
		BuyOrderId      uint64
		SellOrderId     uint64
		StopLossOrderId uint64

		Status BotStatus
	}
)
