package parser

import (
	"fmt"
	"os"

	"github.com/gxdlxss/majestic-bot-parser/notification"
)

// Result содержит все уведомления из файла(ов), разбитые по Go-типам.
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

// Total возвращает суммарное число уведомлений всех типов.
func (r *Result) Total() int {
	return len(r.ItemSales) + len(r.Warehouse) + len(r.Punishments) +
		len(r.PropertySales) + len(r.Rentals) + len(r.OrgAttacks) +
		len(r.OrgDefends) + len(r.SubExpiring) + len(r.SubWarning) + len(r.SubFrozen)
}

// ParseHTMLFiles разбирает один или несколько HTML-файлов экспорта чата
// и возвращает объединённый результат. Все файлы должны быть частями одной переписки.
func ParseHTMLFiles(paths ...string) (*Result, error) {
	result := &Result{}
	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("open %q: %w", path, err)
		}
		notifs, _, err := ParseHTML(f, 0)
		f.Close()
		if err != nil {
			return nil, fmt.Errorf("parse %q: %w", path, err)
		}
		distribute(result, notifs)
	}
	return result, nil
}

// ParseJSONFiles разбирает один или несколько JSON-файлов экспорта чата
// и возвращает объединённый результат. Все файлы должны быть частями одной переписки.
func ParseJSONFiles(paths ...string) (*Result, error) {
	result := &Result{}
	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("open %q: %w", path, err)
		}
		notifs, _, err := ParseJSON(f, 0)
		f.Close()
		if err != nil {
			return nil, fmt.Errorf("parse %q: %w", path, err)
		}
		distribute(result, notifs)
	}
	return result, nil
}

func distribute(r *Result, notifs []interface{}) {
	for _, n := range notifs {
		switch v := n.(type) {
		case *notification.ItemSoldNotification:
			r.ItemSales = append(r.ItemSales, *v)
		case *notification.ItemInWarehouseNotification:
			r.Warehouse = append(r.Warehouse, *v)
		case *notification.PunishmentNotification:
			r.Punishments = append(r.Punishments, *v)
		case *notification.PropertySoldNotification:
			r.PropertySales = append(r.PropertySales, *v)
		case *notification.VehicleRentedNotification:
			r.Rentals = append(r.Rentals, *v)
		case *notification.OrgAttackNotification:
			r.OrgAttacks = append(r.OrgAttacks, *v)
		case *notification.OrgDefendNotification:
			r.OrgDefends = append(r.OrgDefends, *v)
		case *notification.SubscriptionExpiringNotification:
			r.SubExpiring = append(r.SubExpiring, *v)
		case *notification.SubscriptionWarningNotification:
			r.SubWarning = append(r.SubWarning, *v)
		case *notification.SubscriptionFrozenNotification:
			r.SubFrozen = append(r.SubFrozen, *v)
		}
	}
}
