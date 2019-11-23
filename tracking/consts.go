package tracking

const (
	// IndicatorsStatusBuy - статус "необходимо покупать"
	IndicatorsStatusBuy IndicatorsStatus = 1
	// IndicatorsStatusNeutral - статус "нет рекомендаций"
	IndicatorsStatusNeutral IndicatorsStatus = 0
	// IndicatorsStatusSell - статус "необходимо продавать"
	IndicatorsStatusSell IndicatorsStatus = -1

	// BotStatusWaitPurchase - статус "ожидает покупки"
	BotStatusWaitPurchase BotStatus = iota
	// BotStatusActivePurchase - статус "происходит покупка"
	BotStatusActivePurchase BotStatus = iota
	// BotStatusWaitSell - статус "ожидает продажи"
	BotStatusWaitSell BotStatus = iota
	// BotStatusActiveSell - статус "происходит продажа"
	BotStatusActiveSell BotStatus = iota
)
