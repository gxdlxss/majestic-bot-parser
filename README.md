<p align="center">
  <img src="logo.png" alt="chatparser" width="160"/>
  <br/>
  <sub>made by godless</sub>
</p>

# chatparser

Библиотека для разбора HTML-экспорта чата с ботом [@MajesticRolePlayBot](https://t.me/MajesticRolePlayBot).  
Извлекает сырые события — продажи предметов, продажи имущества, наказания — и при необходимости агрегирует их в статистику.

Лицензия: [MIT](../LICENSE)

---

## Структура пакетов

```
lib/
├── event/       — тип Event, константы Kind, Filter, ItemNames
├── parser/      — ParseHTML: разбирает HTML-файл, возвращает []Event
└── aggregator/  — BuildAggregates, TopItems, WindowByCode и сопутствующие типы
```

---

## Установка

```bash
go get market_tg/lib/event
go get market_tg/lib/parser
go get market_tg/lib/aggregator
```

---

## Быстрый старт

```go
import (
    "market_tg/lib/event"
    "market_tg/lib/parser"
)

f, _ := os.Open("messages.html")
defer f.Close()

events, err := parser.ParseHTML(f)
if err != nil {
    log.Fatal(err)
}

for _, e := range events {
    switch e.Kind {
    case event.KindItemSale:
        fmt.Printf("[продажа] %s продал %s × %d за $%.2f\n",
            e.Character, e.ItemName, e.Quantity, e.Price)
    case event.KindPropertySale:
        fmt.Printf("[имущество] %s продал %s за $%.2f (покупатель: %s)\n",
            e.Character, e.ItemName, e.Price, e.Buyer)
    case event.KindPunishment:
        fmt.Printf("[наказание] %s — %s (%d мин) от %s\n",
            e.Character, e.PunType, e.DurationM, e.Admin)
    }
}
```

---

## Пакет `event`

### Типы событий

```
Kind                  Описание
────────────────────  ──────────────────────────
KindItemSale          продажа предмета на маркете
KindPropertySale      продажа имущества
KindPunishment        получение наказания
```

### Структура `Event`

| Поле        | Тип         | Заполняется для                     |
|-------------|-------------|-------------------------------------|
| `Kind`      | `Kind`      | всегда                              |
| `Time`      | `time.Time` | всегда                              |
| `Server`    | `string`    | всегда                              |
| `Character` | `string`    | всегда                              |
| `SaleType`  | `string`    | `KindItemSale`, `KindPropertySale`  |
| `ItemName`  | `string`    | `KindItemSale`, `KindPropertySale`  |
| `Quantity`  | `int`       | `KindItemSale`, `KindPropertySale`  |
| `Price`     | `float64`   | `KindItemSale`, `KindPropertySale`  |
| `Buyer`     | `string`    | `KindPropertySale`                  |
| `PunType`   | `string`    | `KindPunishment`                    |
| `Admin`     | `string`    | `KindPunishment`                    |
| `Reason`    | `string`    | `KindPunishment`                    |
| `DurationM` | `int`       | `KindPunishment` (минуты)           |

```go
e.IsSale()       // true для KindItemSale и KindPropertySale
e.IsPunishment() // true для KindPunishment
```

### Функции

#### `Filter(events []Event, kinds ...Kind) []Event`

Оставляет только события указанных типов.

```go
sales := event.Filter(events, event.KindItemSale, event.KindPropertySale)
puns  := event.Filter(events, event.KindPunishment)
```

#### `ItemNames(events []Event) []string`

Отсортированные уникальные названия предметов из `KindItemSale`.

```go
items := event.ItemNames(events)
// ["Адреналин", "Бронежилет", "Граната", ...]
```

---

## Пакет `parser`

#### `ParseHTML(r io.Reader) ([]Event, error)`

Читает HTML-файл экспорта и возвращает все распознанные события.  
Неизвестные типы сообщений молча пропускаются.

---

## Пакет `aggregator`

### Агрегация по серверам

#### `BuildAggregates(events []Event, now time.Time, window time.Duration) map[string]*Server`

Группирует события по схеме **сервер → персонаж**.  
`window = 0` — учитывать весь период.

```go
import "market_tg/lib/aggregator"

agg := aggregator.BuildAggregates(events, time.Now(), 7*24*time.Hour)

for _, srvName := range aggregator.SortedServerKeys(agg) {
    srv := agg[srvName]
    for _, charID := range aggregator.SortedCharIDs(srv) {
        ch := srv.Characters[charID]
        fmt.Printf("%s / %s: наказаний %d\n", srvName, ch.Name, ch.PunTotal)
        for key, st := range ch.Sales {
            fmt.Printf("  %s %s: %d шт, $%.2f\n", key.SaleType, key.Name, st.Count, st.Sum)
        }
    }
}
```

### Топ предметов

#### `TopItems(events []Event, now time.Time, window time.Duration) (bySum, byCount []ItemStatsRow)`

Топ-5 предметов по сумме и по количеству.

```go
bySum, byCount := aggregator.TopItems(events, time.Now(), 0)
for i, r := range bySum {
    fmt.Printf("%d. %s — $%.2f (%d шт)\n", i+1, r.Label, r.Sum, r.Count)
}
```

### Период

#### `WindowByCode(code string) (label string, duration time.Duration)`

| Код      | Метка       | Длительность      |
|----------|-------------|-------------------|
| `day`    | день        | 24h               |
| `week`   | неделя      | 7 × 24h           |
| `month`  | месяц       | 30 × 24h          |
| `year`   | год         | 365 × 24h         |
| любой    | весь период | 0 (без фильтра)   |

### Вспомогательные функции

```go
aggregator.SortedServerKeys(agg)     // серверы в алфавитном порядке
aggregator.SortedCharIDs(srv)        // персонажи отсортированы по имени
aggregator.SplitCharacter("Имя #ID") // → name="Имя", id="ID"
```

---

---

# chatparser — English

A library for parsing HTML chat exports from [@MajesticRolePlayBot](https://t.me/MajesticRolePlayBot).  
Extracts raw events — item sales, property sales, punishments — and optionally aggregates them into statistics.

License: [MIT](../LICENSE)

## Package layout

```
lib/
├── event/       — Event type, Kind constants, Filter, ItemNames
├── parser/      — ParseHTML: reads the HTML file, returns []Event
└── aggregator/  — BuildAggregates, TopItems, WindowByCode and related types
```

## Install

```bash
go get market_tg/lib/event
go get market_tg/lib/parser
go get market_tg/lib/aggregator
```

## Quick start

```go
import (
    "market_tg/lib/event"
    "market_tg/lib/parser"
)

events, err := parser.ParseHTML(f)

for _, e := range events {
    switch e.Kind {
    case event.KindItemSale:
        fmt.Printf("[sale] %s sold %s × %d for $%.2f\n",
            e.Character, e.ItemName, e.Quantity, e.Price)
    case event.KindPropertySale:
        fmt.Printf("[property] %s sold %s for $%.2f (buyer: %s)\n",
            e.Character, e.ItemName, e.Price, e.Buyer)
    case event.KindPunishment:
        fmt.Printf("[punishment] %s — %s (%d min) by %s\n",
            e.Character, e.PunType, e.DurationM, e.Admin)
    }
}
```

## API summary

| Package      | Function / Type           | Description                                      |
|--------------|---------------------------|--------------------------------------------------|
| `event`      | `Event`, `Kind`           | Core types                                       |
| `event`      | `Filter(events, kinds…)`  | Keep only events of given kinds                  |
| `event`      | `ItemNames(events)`       | Sorted unique item names from `KindItemSale`     |
| `parser`     | `ParseHTML(r)`            | Parse HTML export → `[]Event`                    |
| `aggregator` | `BuildAggregates(…)`      | Group events by server → character               |
| `aggregator` | `TopItems(…)`             | Top-5 items by revenue and quantity              |
| `aggregator` | `WindowByCode(code)`      | Map `"day/week/month/year/all"` → duration       |
| `aggregator` | `SortedServerKeys(agg)`   | Server names sorted alphabetically               |
| `aggregator` | `SortedCharIDs(srv)`      | Character IDs sorted by name                     |
| `aggregator` | `SplitCharacter("N #ID")` | Split into name and ID                           |

Field reference: see the Russian section above — field names are identical.
