<p align="center">
  <img src="logo.png" alt="chatparser" width="160"/>
  <br/>
  <sub>made by godless</sub>
</p>

# chatparser

Разбирает HTML и JSON экспорты чата с [@MajesticRolePlayBot](https://t.me/MajesticRolePlayBot).  
Передаёшь пути к файлам — получаешь готовые Go-структуры в памяти.

Лицензия: [MIT](LICENSE)

---

## Содержание

- [Установка](#установка)
- [Как экспортировать чат](#как-экспортировать-чат)
- [Использование](#использование)
- [Result](#result)
- [Типы уведомлений](#типы-уведомлений)
- [English](#english)

---

## Установка

```bash
go get github.com/gxdlxss/majestic-bot-parser/parser
go get github.com/gxdlxss/majestic-bot-parser/notification
```

---

## Как экспортировать чат

1. Открой чат с **@MajesticRolePlayBot** в Telegram Desktop
2. Нажми на имя бота → **⋯** → **Экспортировать историю чата**
3. Убери галочку с **Photos**, выбери формат **HTML** или **JSON**
4. Нажми **Экспорт**

Если история большая — Telegram разобьёт её на несколько файлов  
(`messages.html`, `messages2.html`, …). Библиотека принимает все сразу.

---

## Использование

```go
import "github.com/gxdlxss/majestic-bot-parser/parser"
```

### HTML

```go
// один файл
result, err := parser.ParseHTMLFiles("messages.html")

// несколько частей одной переписки — объединяются в один результат
result, err := parser.ParseHTMLFiles("messages.html", "messages2.html", "messages3.html")
```

### JSON

```go
result, err := parser.ParseJSONFiles("result.json")

result, err := parser.ParseJSONFiles("result.json", "result2.json")
```

> HTML и JSON в одном вызове не смешиваются.

### Работа с результатом

```go
if err != nil {
    log.Fatal(err)
}

fmt.Println("Всего:", result.Total())

for _, s := range result.ItemSales {
    fmt.Printf("[продажа] %s — %s × %d за $%.2f\n",
        s.Character, s.ItemName, s.Quantity, s.SalePrice)
}

for _, p := range result.Punishments {
    fmt.Printf("[наказание] %s — %s (%.1f ч), причина: %s\n",
        p.Character, p.PunishmentType, p.DurationHours, p.Reason)
}

for _, r := range result.Rentals {
    fmt.Printf("[аренда] %s → %s, $%.2f, %.1f ч\n",
        r.VehicleName, r.Renter, r.Price, r.DurationHours)
}

for _, a := range result.OrgAttacks {
    fmt.Printf("[атака] %s напала на %s, квадрат %s\n",
        a.OrgName, a.EnemyName, a.SquareName)
}

for _, d := range result.OrgDefends {
    fmt.Printf("[защита] на %s напала %s\n",
        d.OrgName, d.AttackerName)
}

for _, w := range result.Warehouse {
    fmt.Printf("[склад] %s — %s × %d, забрать за %d ч\n",
        w.Character, w.ItemName, w.Quantity, w.PickupDeadlineHours)
}

for _, e := range result.SubExpiring {
    fmt.Printf("[подписка] %s истекает через %d дн\n", e.Login, e.DaysLeft)
}
```

---

## Result

```go
type Result struct {
    ItemSales     []notification.ItemSoldNotification
    Warehouse     []notification.ItemInWarehouseNotification
    Punishments   []notification.PunishmentNotification
    PropertySales []notification.PropertySoldNotification
    Rentals       []notification.VehicleRentedNotification
    OrgAttacks    []notification.OrgAttackNotification
    OrgDefends    []notification.OrgDefendNotification
    SubExpiring   []notification.SubscriptionExpiringNotification
    SubWarning    []notification.SubscriptionWarningNotification
    SubFrozen     []notification.SubscriptionFrozenNotification
}
```

| Поле            | Что в нём                                |
|-----------------|------------------------------------------|
| `ItemSales`     | Продажи предметов на маркете             |
| `Warehouse`     | Предметы, доставленные на склад          |
| `Punishments`   | Наказания персонажей                     |
| `PropertySales` | Продажи имущества                        |
| `Rentals`       | Сдача транспорта в аренду                |
| `OrgAttacks`    | Нападения вашей организации на чужую     |
| `OrgDefends`    | Нападения чужой организации на вашу      |
| `SubExpiring`   | Подписка заканчивается через N дней      |
| `SubWarning`    | Подписка скоро заканчивается             |
| `SubFrozen`     | Подписка заморожена                      |

```go
result.Total() // int — сумма длин всех срезов
```

---

## Типы уведомлений

Каждый тип встраивает `BaseNotification`:

```go
type BaseNotification struct {
    Type          NotificationType // "item_sold", "punishment", …
    Server        string           // название сервера
    CharID        string           // статический ID персонажа, напр. "12345"
    Timestamp     time.Time        // время события
    HTMLMessageID int              // ID сообщения в экспорте
}
```

---

### ItemSoldNotification

```go
Character string
ItemName  string
Quantity  int
SalePrice float64
Buyer     string
```

---

### ItemInWarehouseNotification

```go
Character           string
ItemName            string
Quantity            int
PickupDeadlineHours int
```

---

### PunishmentNotification

```go
Character      string
PunishmentType string
Admin          string
Reason         string
DurationHours  float64
```

---

### PropertySoldNotification

```go
Character string
Name      string
SalePrice float64
Buyer     string
```

---

### VehicleRentedNotification

```go
Character     string
VehicleName   string
VehicleNumber string
Price         float64
DurationHours float64
Renter        string
```

---

### OrgAttackNotification

```go
Character     string
OrgName       string    // ваша организация
EnemyName     string    // на кого напали
AttackStart   time.Time
SquareName    string
SquareNumber  string
AttackerCount int
WeaponCaliber string    // пусто = семейный кап
```

---

### OrgDefendNotification

```go
Character     string
OrgName       string    // ваша организация
AttackerName  string    // кто напал
AttackStart   time.Time
SquareName    string
SquareNumber  string
AttackerCount int
WeaponCaliber string    // пусто = семейный кап
```

---

### SubscriptionExpiringNotification

```go
Login    string
DaysLeft int
```

---

### SubscriptionWarningNotification / SubscriptionFrozenNotification

Только поля из `BaseNotification`.

---

---

## English

Parses Telegram chat exports from [@MajesticRolePlayBot](https://t.me/MajesticRolePlayBot).  
Pass file paths — get back plain Go structs. No BSON, no files written.

License: [MIT](LICENSE)

### Install

```bash
go get github.com/gxdlxss/majestic-bot-parser/parser
go get github.com/gxdlxss/majestic-bot-parser/notification
```

### Usage

```go
import "github.com/gxdlxss/majestic-bot-parser/parser"

// HTML (one or more files from the same chat)
result, err := parser.ParseHTMLFiles("messages.html", "messages2.html")

// JSON (one or more files from the same chat)
result, err := parser.ParseJSONFiles("result.json")

fmt.Println(result.Total())

for _, s := range result.ItemSales {
    fmt.Println(s.Character, s.ItemName, s.Quantity, s.SalePrice)
}
for _, p := range result.Punishments {
    fmt.Println(p.Character, p.PunishmentType, p.DurationHours)
}
for _, r := range result.Rentals {
    fmt.Println(r.VehicleName, r.Renter, r.Price)
}
```

> HTML and JSON cannot be mixed in a single call.

### Result fields

| Field | Type | Description |
|---|---|---|
| `ItemSales` | `[]ItemSoldNotification` | Items sold on the market |
| `Warehouse` | `[]ItemInWarehouseNotification` | Items at warehouse |
| `Punishments` | `[]PunishmentNotification` | Character punishments |
| `PropertySales` | `[]PropertySoldNotification` | Property sold |
| `Rentals` | `[]VehicleRentedNotification` | Vehicles rented out |
| `OrgAttacks` | `[]OrgAttackNotification` | Your org attacked another |
| `OrgDefends` | `[]OrgDefendNotification` | Your org was attacked |
| `SubExpiring` | `[]SubscriptionExpiringNotification` | Subscription expiring in N days |
| `SubWarning` | `[]SubscriptionWarningNotification` | Subscription expiring soon |
| `SubFrozen` | `[]SubscriptionFrozenNotification` | Subscription frozen |

All types embed `BaseNotification`: `Type`, `Server`, `CharID`, `Timestamp`, `HTMLMessageID`.
