// Package aggregator группирует события по серверам и персонажам
// и считает статистику продаж и наказаний.
package aggregator

import (
	"sort"
	"strings"
	"time"

	"market_tg/lib/event"
)

// ItemStats — суммарная статистика по одному предмету.
type ItemStats struct {
	Count int
	Sum   float64
}

// SaleKey — составной ключ агрегации: тип продажи + название предмета.
type SaleKey struct {
	SaleType string
	Name     string
}

// Character — персонаж с накопленной статистикой продаж и наказаний.
type Character struct {
	ID          string
	Name        string
	LastSeen    time.Time
	Sales       map[SaleKey]*ItemStats
	PunTotal    int
	PunByType   map[string]int
	PunLastTime time.Time
}

// Server — игровой сервер со списком персонажей.
type Server struct {
	Name       string
	Characters map[string]*Character
}

// ItemStatsRow — строка в топ-листе предметов.
type ItemStatsRow struct {
	Label string
	Count int
	Sum   float64
}

// WindowByCode преобразует строковый код периода в читаемую метку и длительность.
// Нулевая длительность означает «весь период» (без ограничения).
func WindowByCode(code string) (label string, duration time.Duration) {
	switch code {
	case "day":
		return "день", 24 * time.Hour
	case "week":
		return "неделя", 7 * 24 * time.Hour
	case "month":
		return "месяц", 30 * 24 * time.Hour
	case "year":
		return "год", 365 * 24 * time.Hour
	default:
		return "весь период", 0
	}
}

// BuildAggregates группирует события по схеме сервер → персонаж.
// window=0 — учитывать все события без ограничения по времени.
func BuildAggregates(events []event.Event, now time.Time, window time.Duration) map[string]*Server {
	servers := make(map[string]*Server)

	ensureChar := func(serverName, characterFull string, t time.Time) *Character {
		namePart, idPart := SplitCharacter(characterFull)
		if idPart == "" {
			idPart = namePart
		}
		srv := servers[serverName]
		if srv == nil {
			srv = &Server{Name: serverName, Characters: make(map[string]*Character)}
			servers[serverName] = srv
		}
		ch := srv.Characters[idPart]
		if ch == nil {
			ch = &Character{
				ID:        idPart,
				Name:      namePart,
				LastSeen:  t,
				Sales:     make(map[SaleKey]*ItemStats),
				PunByType: make(map[string]int),
			}
			srv.Characters[idPart] = ch
		} else if t.After(ch.LastSeen) {
			ch.Name = namePart
			ch.LastSeen = t
		}
		return ch
	}

	for _, e := range events {
		if window > 0 && now.Sub(e.Time) > window {
			continue
		}
		ch := ensureChar(e.Server, e.Character, e.Time)
		switch e.Kind {
		case event.KindItemSale, event.KindPropertySale:
			key := SaleKey{SaleType: e.SaleType, Name: e.ItemName}
			st := ch.Sales[key]
			if st == nil {
				st = &ItemStats{}
				ch.Sales[key] = st
			}
			st.Count += e.Quantity
			st.Sum += e.Price
		case event.KindPunishment:
			ch.PunTotal++
			ch.PunByType[e.PunType]++
			if ch.PunLastTime.IsZero() || e.Time.After(ch.PunLastTime) {
				ch.PunLastTime = e.Time
			}
		}
	}

	return servers
}

// TopItems возвращает топ-5 предметов по сумме и по количеству за указанный период.
func TopItems(events []event.Event, now time.Time, window time.Duration) (bySum, byCount []ItemStatsRow) {
	m := map[string]*ItemStats{}
	for _, e := range events {
		if !e.IsSale() {
			continue
		}
		if window > 0 && now.Sub(e.Time) > window {
			continue
		}
		label := e.SaleType + ": " + e.ItemName
		st := m[label]
		if st == nil {
			st = &ItemStats{}
			m[label] = st
		}
		st.Count += e.Quantity
		st.Sum += e.Price
	}

	for label, st := range m {
		row := ItemStatsRow{Label: label, Count: st.Count, Sum: st.Sum}
		bySum = append(bySum, row)
		byCount = append(byCount, row)
	}

	sort.Slice(bySum, func(i, j int) bool {
		if bySum[i].Sum != bySum[j].Sum {
			return bySum[i].Sum > bySum[j].Sum
		}
		return bySum[i].Label < bySum[j].Label
	})
	sort.Slice(byCount, func(i, j int) bool {
		if byCount[i].Count != byCount[j].Count {
			return byCount[i].Count > byCount[j].Count
		}
		return byCount[i].Label < byCount[j].Label
	})

	if len(bySum) > 5 {
		bySum = bySum[:5]
	}
	if len(byCount) > 5 {
		byCount = byCount[:5]
	}
	return
}

// SortedServerKeys возвращает ключи карты серверов в алфавитном порядке.
func SortedServerKeys(m map[string]*Server) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// SortedCharIDs возвращает ID персонажей сервера, отсортированные по имени.
func SortedCharIDs(srv *Server) []string {
	keys := make([]string, 0, len(srv.Characters))
	for id := range srv.Characters {
		keys = append(keys, id)
	}
	sort.Slice(keys, func(i, j int) bool {
		return srv.Characters[keys[i]].Name < srv.Characters[keys[j]].Name
	})
	return keys
}

// SplitCharacter разбирает строку вида "Имя #ID" на имя и ID.
func SplitCharacter(full string) (name, id string) {
	if i := strings.LastIndex(full, "#"); i != -1 {
		return strings.TrimSpace(full[:i]), strings.TrimSpace(full[i+1:])
	}
	return strings.TrimSpace(full), ""
}
