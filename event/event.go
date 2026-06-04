// Package event определяет типы событий из HTML-экспорта чата
// и вспомогательные функции для работы с ними.
package event

import (
	"sort"
	"time"
)

// Kind — тип события.
type Kind string

const (
	KindItemSale     Kind = "item_sale"     // продажа предмета
	KindPropertySale Kind = "property_sale" // продажа имущества
	KindPunishment   Kind = "punishment"    // наказание
)

// Event — одно событие из экспорта чата.
// Поля заполняются в зависимости от Kind:
//   - KindItemSale, KindPropertySale → SaleType, ItemName, Quantity, Price, Buyer
//   - KindPunishment                 → PunType, Admin, Reason, DurationM
type Event struct {
	Kind      Kind
	Time      time.Time
	Server    string
	Character string

	// поля продажи
	SaleType string  // "Предмет" | "Имущество"
	ItemName string
	Quantity int
	Price    float64
	Buyer    string // только для KindPropertySale

	// поля наказания
	PunType   string
	Admin     string
	Reason    string
	DurationM int
}

func (e Event) IsSale() bool       { return e.Kind == KindItemSale || e.Kind == KindPropertySale }
func (e Event) IsPunishment() bool { return e.Kind == KindPunishment }

// Filter возвращает только события указанных типов.
func Filter(events []Event, kinds ...Kind) []Event {
	set := make(map[Kind]struct{}, len(kinds))
	for _, k := range kinds {
		set[k] = struct{}{}
	}
	out := make([]Event, 0, len(events))
	for _, e := range events {
		if _, ok := set[e.Kind]; ok {
			out = append(out, e)
		}
	}
	return out
}

// ItemNames возвращает отсортированный список уникальных названий предметов
// из событий KindItemSale.
func ItemNames(events []Event) []string {
	seen := make(map[string]struct{})
	for _, e := range events {
		if e.Kind == KindItemSale {
			seen[e.ItemName] = struct{}{}
		}
	}
	names := make([]string, 0, len(seen))
	for n := range seen {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
