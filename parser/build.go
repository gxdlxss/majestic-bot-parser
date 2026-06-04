package parser

import (
	"strconv"
	"strings"
	"time"

	"github.com/gxdlxss/majestic-bot-parser/notification"
)

type rawMessage struct {
	HTMLID    int
	Timestamp time.Time
	Text      string
}

func buildNotification(m rawMessage) interface{} {
	lines := strings.Split(m.Text, "\n")
	var parts []string
	for _, l := range lines {
		if l = strings.TrimSpace(l); l != "" {
			parts = append(parts, l)
		}
	}
	if len(parts) == 0 {
		return nil
	}

	title := parts[0]
	fields := parseFields(parts[1:])

	base := notification.BaseNotification{
		Timestamp:     m.Timestamp,
		Server:        fields["Сервер"],
		HTMLMessageID: m.HTMLID,
	}

	switch {
	case strings.HasPrefix(title, "Вы успешно продали имущество"):
		base.Type = notification.TypePropertySold
		charStr := fields["Персонаж"]
		base.CharID = extractCharID(charStr)
		return &notification.PropertySoldNotification{
			BaseNotification: base,
			Character:        extractCharName(charStr),
			Name:             fields["Название"],
			SalePrice:        parsePrice(fields["Цена продажи"]),
			Buyer:            fields["Покупатель"],
		}

	case strings.HasPrefix(title, "Вы успешно продали"):
		base.Type = notification.TypeItemSold
		charStr := fields["Персонаж"]
		base.CharID = extractCharID(charStr)
		qty, _ := strconv.Atoi(fields["Кол-во"])
		return &notification.ItemSoldNotification{
			BaseNotification: base,
			Character:        extractCharName(charStr),
			ItemName:         fields["Название"],
			Quantity:         qty,
			SalePrice:        parsePrice(fields["Цена продажи"]),
			Buyer:            fields["Покупатель"],
		}

	case strings.HasPrefix(title, "Предмет на складе"):
		base.Type = notification.TypeItemInWarehouse
		charStr := fields["Персонаж"]
		base.CharID = extractCharID(charStr)
		qty, _ := strconv.Atoi(fields["Количество"])
		return &notification.ItemInWarehouseNotification{
			BaseNotification:    base,
			Character:           extractCharName(charStr),
			ItemName:            fields["Предмет"],
			Quantity:            qty,
			PickupDeadlineHours: extractDeadlineHours(m.Text),
		}

	case strings.HasPrefix(title, "Вы получили наказание"):
		base.Type = notification.TypePunishment
		charStr := fields["Персонаж"]
		base.CharID = extractCharID(charStr)
		return &notification.PunishmentNotification{
			BaseNotification: base,
			Character:        extractCharName(charStr),
			PunishmentType:   fields["Тип"],
			Admin:            fields["Администратор"],
			Reason:           fields["Причина"],
			DurationHours:    parseDurationHours(fields["Длительность"]),
		}

	case strings.Contains(title, "осталось") && strings.Contains(title, "дн"):
		base.Type = notification.TypeSubscriptionExpiring
		if s := serverFromSubscriptionTitle(title); s != "" {
			base.Server = s
		}
		return &notification.SubscriptionExpiringNotification{
			BaseNotification: base,
			Login:            fields["Логин"],
			DaysLeft:         extractDaysLeft(title),
		}

	case strings.Contains(title, "заканчивается через"):
		base.Type = notification.TypeSubscriptionWarning
		if s := serverFromSubscriptionTitle(title); s != "" {
			base.Server = s
		}
		return &notification.SubscriptionWarningNotification{BaseNotification: base}

	case strings.Contains(title, "заморожена"):
		base.Type = notification.TypeSubscriptionFrozen
		if s := serverFromSubscriptionTitle(title); s != "" {
			base.Server = s
		}
		return &notification.SubscriptionFrozenNotification{BaseNotification: base}

	case strings.HasPrefix(title, "Транспорт сдан в аренду"):
		base.Type = notification.TypeVehicleRented
		charStr := fields["Персонаж"]
		base.CharID = extractCharID(charStr)
		return &notification.VehicleRentedNotification{
			BaseNotification: base,
			Character:        extractCharName(charStr),
			VehicleName:      fields["Транспорт"],
			VehicleNumber:    fields["Номер транспорта"],
			Price:            parsePrice(fields["Цена"]),
			DurationHours:    parseDurationHours(fields["Длительность"]),
			Renter:           fields["Арендатор"],
		}

	case strings.HasPrefix(title, "Ваша организация") && strings.Contains(title, "напала"):
		base.Type = notification.TypeOrgAttack
		charStr := fields["Персонаж"]
		base.CharID = extractCharID(charStr)
		orgName, enemyName := extractOrgAndEnemy(title, "напала на")
		count, _ := strconv.Atoi(fields["Количество нападающих"])
		return &notification.OrgAttackNotification{
			BaseNotification: base,
			Character:        extractCharName(charStr),
			OrgName:          orgName,
			EnemyName:        enemyName,
			AttackStart:      parseAttackTime(fields["Начало"]),
			SquareName:       fields["Название квадрата"],
			SquareNumber:     fields["Номер квадрата"],
			AttackerCount:    count,
			WeaponCaliber:    fields["Калибр оружия"],
		}

	case strings.HasPrefix(title, "На вашу организацию") && strings.Contains(title, "напали"):
		base.Type = notification.TypeOrgDefend
		charStr := fields["Персонаж"]
		base.CharID = extractCharID(charStr)
		orgName, attackerName := extractOrgAndEnemy(title, "напали")
		count, _ := strconv.Atoi(fields["Количество нападающих"])
		return &notification.OrgDefendNotification{
			BaseNotification: base,
			Character:        extractCharName(charStr),
			OrgName:          orgName,
			AttackerName:     attackerName,
			AttackStart:      parseAttackTime(fields["Начало"]),
			SquareName:       fields["Название квадрата"],
			SquareNumber:     fields["Номер квадрата"],
			AttackerCount:    count,
			WeaponCaliber:    fields["Калибр оружия"],
		}
	}

	return nil
}

func parseFields(lines []string) map[string]string {
	fields := make(map[string]string, len(lines))
	for _, line := range lines {
		if idx := strings.Index(line, ": "); idx > 0 {
			fields[strings.TrimSpace(line[:idx])] = strings.TrimSpace(line[idx+2:])
		}
	}
	return fields
}

func extractCharID(character string) string {
	idx := strings.LastIndex(character, "#")
	if idx < 0 || idx+1 >= len(character) {
		return ""
	}
	return strings.TrimSpace(character[idx+1:])
}

func extractCharName(character string) string {
	if idx := strings.LastIndex(character, " #"); idx > 0 {
		return strings.TrimSpace(character[:idx])
	}
	return character
}

func parsePrice(s string) float64 {
	s = strings.TrimPrefix(s, "$")
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, ",", ".")
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func parseDurationHours(s string) float64 {
	parts := strings.Fields(s)
	if len(parts) < 2 {
		n, _ := strconv.ParseFloat(s, 64)
		return n / 60.0
	}
	n, _ := strconv.ParseFloat(parts[0], 64)
	unit := strings.ToLower(parts[1])
	switch {
	case strings.HasPrefix(unit, "мин"):
		return n / 60.0
	case strings.HasPrefix(unit, "час"):
		return n
	case strings.HasPrefix(unit, "ден"), strings.HasPrefix(unit, "дн"), strings.HasPrefix(unit, "сут"):
		return n * 24.0
	}
	return n / 60.0
}

func extractDeadlineHours(text string) int {
	idx := strings.Index(text, "в течение ")
	if idx < 0 {
		return 48
	}
	parts := strings.Fields(text[idx+len("в течение "):])
	if len(parts) >= 1 {
		if n, err := strconv.Atoi(parts[0]); err == nil {
			return n
		}
	}
	return 48
}

func extractDaysLeft(title string) int {
	idx := strings.LastIndex(title, "осталось ")
	if idx < 0 {
		return 0
	}
	parts := strings.Fields(title[idx+len("осталось "):])
	if len(parts) >= 1 {
		n, _ := strconv.Atoi(parts[0])
		return n
	}
	return 0
}

func serverFromSubscriptionTitle(title string) string {
	idx := strings.Index(title, "на сервере ")
	if idx < 0 {
		return ""
	}
	s := title[idx+len("на сервере "):]
	if end := strings.IndexAny(s, " !.,"); end > 0 {
		return s[:end]
	}
	return s
}

func extractOrgAndEnemy(title, verb string) (orgName, otherName string) {
	var orgStart int
	switch {
	case strings.HasPrefix(title, "Ваша организация "):
		orgStart = len("Ваша организация ")
	case strings.HasPrefix(title, "На вашу организацию "):
		orgStart = len("На вашу организацию ")
	default:
		return "", ""
	}
	rest := title[orgStart:]
	sep := " " + verb + " "
	sepIdx := strings.Index(rest, sep)
	if sepIdx < 0 {
		sep = " " + verb
		sepIdx = strings.LastIndex(rest, sep)
	}
	if sepIdx < 0 {
		return strings.TrimSuffix(rest, "!"), ""
	}
	orgName = rest[:sepIdx]
	otherName = strings.TrimSuffix(strings.TrimSpace(rest[sepIdx+len(sep):]), "!")
	return
}

func parseAttackTime(s string) time.Time {
	t, _ := time.ParseInLocation("02.01.2006 15:04:05", s, time.Local)
	return t
}
