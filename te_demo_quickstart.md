# te_demo CLI — Quick Start

## 1. Установка

Откройте **PowerShell** и выполните одну команду:

```powershell
irm https://raw.githubusercontent.com/dmikhailenko-hash/te_demo/main/install.ps1 | iex
```

Закройте PowerShell и откройте **новый** — PATH обновится автоматически.

---

## 2. Получение токена

```powershell
$r = Invoke-RestMethod -Method POST -Uri "https://webhooks-clientapi.traderevolution.com/traderevolution/v1/authorize?login=YOUR_LOGIN&password=YOUR_PASSWORD"
te_demo config --token $r.d.access_token
```

Замените `YOUR_LOGIN` и `YOUR_PASSWORD` на ваши учётные данные.  
Токен сохраняется локально в `~/.te_demo/config.json`.

---

## 3. Первые запросы

```powershell
# Список аккаунтов
te_demo accounts

# Баланс аккаунта
te_demo accounts state 12345

# Открытые позиции
te_demo accounts positions 12345

# Активные ордера
te_demo accounts orders 12345

# Доступные инструменты
te_demo accounts instruments 12345
```

---

## 4. Торговля

```powershell
# Выставить ордер
te_demo order place 12345 --tradableInstrumentId 100 --side buy --type market --qty 1 --validity DAY

# Отменить ордер
te_demo order cancel 67890

# Закрыть позицию
te_demo position close 111
```

---

## 5. Рыночные данные

```powershell
# Котировки
te_demo quotes --tradableInstrumentId 100 --accountId 12345

# История баров (1H за период)
te_demo history --tradableInstrumentId 100 --accountId 12345 --resolution 1H --from 1700000000000 --to 1710000000000

# Стакан
te_demo depth --tradableInstrumentId 100 --accountId 12345
```

---

## Ссылки

- **Репозиторий:** https://github.com/dmikhailenko-hash/te_demo
- **Releases:** https://github.com/dmikhailenko-hash/te_demo/releases
- **API документация:** https://webhooks-clientapi.traderevolution.com/traderevolution/v1/swagger-ui/index.html

---

## Смена токена

Токен действует ~1 час. Для обновления:

```powershell
$r = Invoke-RestMethod -Method POST -Uri "https://webhooks-clientapi.traderevolution.com/traderevolution/v1/authorize?login=YOUR_LOGIN&password=YOUR_PASSWORD"
te_demo config --token $r.d.access_token
```
