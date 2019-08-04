package tracking

const (
	IndicatorsStatusBuy     IndicatorsStatus = 1
	IndicatorsStatusNeutral IndicatorsStatus = 0
	IndicatorsStatusSell    IndicatorsStatus = -1
)

const (
	BotStatusWaitPurchase   BotStatus = iota
	BotStatusActivePurchase BotStatus = iota
	BotStatusWaitSell       BotStatus = iota
	BotStatusActiveSell     BotStatus = iota
)
