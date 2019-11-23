package tracking

type (
	// IndicatorsStatus - статус индикатора
	IndicatorsStatus int8

	// BotStatus - статус бота
	BotStatus uint8

	// Bot - хранилище показателей бота
	Bot struct {
		// BuyOrderID - идентификатор ордера покупки
		BuyOrderID uint64
		// SellOrderID - идентификатор ордера продажи
		SellOrderID uint64
		// StopLossOrderID - идентификатор ордера STOP-LOSS
		StopLossOrderID uint64

		Status BotStatus
	}
)
