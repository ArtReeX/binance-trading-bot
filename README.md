# Бот для торговли на криптобирже Binance на основе нейронной сети #

### Запуск: ###
***
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

  `base` - базовая валюта (пример: "BTC")

  `quote` - валюта котировки (пример: "USDT")

  `intervals` - временной промежуток для анализа валюты (пример: "15m")

  `priceForOneTransaction` - цена на одну сделку по выбранной паре (пример: 10)

- Запустите бота `binance-trading-bot`.