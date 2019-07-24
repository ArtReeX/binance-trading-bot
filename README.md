# binance-trading-bot

Бот для торговли на криптобирже Binance по индикаторам.

Перед запуском создайте файл "config.json" со следующим содержанием:
{
"API": {
"Binance": {
"key": "YOUR_KEY_BINANCE_API",
"secret": "YOU_SECRET_BINANCE_API"
}
},
"directions": [
{
"base": "BTC",
"quote": "USDT",
"intervals": ["1m", "5m", "15m"],
"accuracyQuantity": 6,
"accuracyPrice": 2
}
]

}

base - базовая валюта (пример: "BTC")
quote - валюта котировки (пример: "USDT")
accuracyQuantity - количество знаков после запятой допустимых биржой для основной валюты
accuracyPrice - количество знаков после запятой допустимых биржой для котировочной валюты
intervals - задаёт временной промежуток для анализа валюты (пример: "15m")
