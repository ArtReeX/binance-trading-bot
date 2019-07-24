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
"intervals": ["1m", "5m", "15m"]
}
]
}

base - базовая валюта (пример: "BTC")
quote - валюта котировки (пример: "USDT")
intervals - задаёт временной промежуток для анализа валюты (пример: "15m")
