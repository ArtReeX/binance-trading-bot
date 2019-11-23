# Бот для торговли на криптобирже Binance на основе нейронной сети

## Запуск

---

- Выполните:

```bash
git clone https://github.com/ArtReeX/binance-trading-bot.git
cd binance-trading-bot
go build
mkdir $HOME/.binance-trading-bot
cp config.json $HOME/.binance-trading-bot/config.json
open $HOME/.binance-trading-bot/config.json
```

- Настройте файл конфигурации учитывая следующие параметры.

  `key` - ваш KEY для биржи

  `secret` - ваш SECRET для биржи

  `fee` - размер комиссии за сделку (при наличии BNB комиссия составляет 0.075%)

  `pair` - валютная пара (пример: "BTCUSDT")

  `intervals` - временной промежуток для анализа валюты (пример: "15m")

  `priceForOneTransaction` - цена на одну сделку по выбранной паре (пример: 10)

- Запустите бота `binance-trading-bot`.
