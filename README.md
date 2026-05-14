# te_demo — TraderEvolution CLI

Command-line interface для TradeRevolution REST API.

## Установка

### Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/YOUR_ORG/te_demo/main/install.ps1 | iex
```

### Linux / macOS
```sh
curl -fsSL https://raw.githubusercontent.com/YOUR_ORG/te_demo/main/install.sh | sh
```

## Настройка

Получите токен и сохраните его:
```powershell
$r = Invoke-RestMethod -Method POST -Uri "https://webhooks-clientapi.traderevolution.com/traderevolution/v1/authorize?login=YOUR_LOGIN&password=YOUR_PASSWORD"
te_demo config --token $r.d.access_token
```

## Использование

```bash
te_demo accounts                          # список аккаунтов
te_demo accounts state <id>              # баланс аккаунта
te_demo accounts positions <id>          # открытые позиции
te_demo accounts orders <id>             # активные ордера
te_demo accounts orders-history <id>     # история ордеров
te_demo accounts executions <id>         # исполнения
te_demo accounts instruments <id>        # инструменты

te_demo quotes --tradableInstrumentId 100 --accountId 12345
te_demo history --tradableInstrumentId 100 --accountId 12345 --resolution 1H --from 1700000000000 --to 1710000000000

te_demo order place 12345 --tradableInstrumentId 100 --side buy --type market --qty 1 --validity DAY
te_demo order cancel 67890
te_demo order cancel-all 12345

te_demo position close 111
te_demo position close-all 12345
te_demo position modify 111 --stopLoss 1.2000 --takeProfit 1.3000
```

## Конфигурация

```bash
te_demo config                    # показать настройки
te_demo config --token <token>    # установить токен
te_demo config --url <url>        # сменить базовый URL
```

Конфиг хранится в `~/.te_demo/config.json`.
