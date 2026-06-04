// Package parser разбирает HTML-экспорт чата Telegram и возвращает события.
package parser

import (
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"market_tg/lib/event"
)

var (
	itemSaleRe = regexp.MustCompile(`(?s)Сервер:\s*(.+?)\s*Персонаж:\s*(.+?)\s*(?:Название|Предмет):\s*(.+?)\s*(?:Кол-во|Количество):\s*([0-9]+)\s*Цена продажи:\s*\$([0-9\s,]+)`)
	propertySaleRe = regexp.MustCompile(`(?s)Сервер:\s*(.+?)\s*Персонаж:\s*(.+?)\s*Название:\s*(.+?)\s*Цена продажи:\s*\$([0-9\s,]+)\s*(?:Покупатель:\s*(.+?))?\s*$`)
	punRe          = regexp.MustCompile(`(?s)Сервер:\s*(.+?)\s*Персонаж:\s*(.+?)\s*Тип:\s*(.+?)\s*Администратор:\s*(.+?)\s*Причина:\s*(.+?)\s*Длительность:\s*([0-9]+)\s*минут`)
)

// ParseHTML читает HTML-файл экспорта и возвращает все распознанные события
// в хронологическом порядке. Неизвестные типы сообщений молча пропускаются.
func ParseHTML(r io.Reader) ([]event.Event, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	var events []event.Event

	doc.Find("div.message").Each(func(_ int, msg *goquery.Selection) {
		text := msg.Find("div.text").Text()

		dateTitle, ok := msg.Find("div.pull_right.date.details").Attr("title")
		if !ok {
			return
		}
		ts := strings.Split(dateTitle, " UTC")[0]
		msgTime, err := time.ParseInLocation("02.01.2006 15:04:05", ts, time.Local)
		if err != nil {
			return
		}

		switch {
		case strings.Contains(text, "Вы успешно продали предмет"):
			if e, ok := parseItemSale(text, msgTime); ok {
				events = append(events, e)
			}
		case strings.Contains(text, "Вы успешно продали имущество"):
			if e, ok := parsePropertySale(text, msgTime); ok {
				events = append(events, e)
			}
		case strings.Contains(text, "Вы получили наказание"):
			if e, ok := parsePunishment(text, msgTime); ok {
				events = append(events, e)
			}
		}
	})

	return events, nil
}

func parseItemSale(text string, t time.Time) (event.Event, bool) {
	m := itemSaleRe.FindStringSubmatch(text)
	if len(m) != 6 {
		return event.Event{}, false
	}
	name := strings.TrimSpace(m[3])
	if name == "Улучшенный эпинефрин" {
		name = "Адреналин"
	}
	qty, _ := strconv.Atoi(m[4])
	price, _ := strconv.ParseFloat(normalisePrice(m[5]), 64)
	return event.Event{
		Kind:      event.KindItemSale,
		Time:      t,
		Server:    strings.TrimSpace(m[1]),
		Character: strings.TrimSpace(m[2]),
		SaleType:  "Предмет",
		ItemName:  name,
		Quantity:  qty,
		Price:     price,
	}, true
}

func parsePropertySale(text string, t time.Time) (event.Event, bool) {
	m := propertySaleRe.FindStringSubmatch(text)
	if len(m) < 5 {
		return event.Event{}, false
	}
	price, _ := strconv.ParseFloat(normalisePrice(m[4]), 64)
	buyer := ""
	if len(m) >= 6 {
		buyer = strings.TrimSpace(m[5])
	}
	return event.Event{
		Kind:      event.KindPropertySale,
		Time:      t,
		Server:    strings.TrimSpace(m[1]),
		Character: strings.TrimSpace(m[2]),
		SaleType:  "Имущество",
		ItemName:  strings.TrimSpace(m[3]),
		Quantity:  1,
		Price:     price,
		Buyer:     buyer,
	}, true
}

func parsePunishment(text string, t time.Time) (event.Event, bool) {
	m := punRe.FindStringSubmatch(text)
	if len(m) != 7 {
		return event.Event{}, false
	}
	dur, _ := strconv.Atoi(strings.TrimSpace(m[6]))
	return event.Event{
		Kind:      event.KindPunishment,
		Time:      t,
		Server:    strings.TrimSpace(m[1]),
		Character: strings.TrimSpace(m[2]),
		PunType:   strings.TrimSpace(m[3]),
		Admin:     strings.TrimSpace(m[4]),
		Reason:    strings.TrimSpace(m[5]),
		DurationM: dur,
	}, true
}

func normalisePrice(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), ",", ".")
}
