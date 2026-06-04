package parser

import (
	"encoding/json"
	"io"
	"sort"
	"strings"
	"time"
)

type jsonExport struct {
	Messages []jsonMessage `json:"messages"`
}

type jsonMessage struct {
	ID   int             `json:"id"`
	Date string          `json:"date"`
	Text json.RawMessage `json:"text"`
}

// ParseJSON разбирает JSON-экспорт Telegram.
// Возвращает уведомления с ID > lastID, новый максимальный ID и ошибку.
func ParseJSON(r io.Reader, lastID int) ([]interface{}, int, error) {
	var export jsonExport
	if err := json.NewDecoder(r).Decode(&export); err != nil {
		return nil, lastID, err
	}

	sort.Slice(export.Messages, func(i, j int) bool {
		return export.Messages[i].ID < export.Messages[j].ID
	})

	maxID := lastID
	var notifications []interface{}

	for _, m := range export.Messages {
		if m.ID <= lastID {
			continue
		}
		text := extractJSONText(m.Text)
		if text == "" {
			if m.ID > maxID {
				maxID = m.ID
			}
			continue
		}
		ts, _ := time.ParseInLocation("2006-01-02T15:04:05", m.Date, time.Local)
		raw := rawMessage{HTMLID: m.ID, Timestamp: ts, Text: text}
		if n := buildNotification(raw); n != nil {
			notifications = append(notifications, n)
		}
		if m.ID > maxID {
			maxID = m.ID
		}
	}

	return notifications, maxID, nil
}

func extractJSONText(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return strings.TrimSpace(s)
	}
	var parts []json.RawMessage
	if json.Unmarshal(raw, &parts) != nil {
		return ""
	}
	var sb strings.Builder
	for _, p := range parts {
		var str string
		if json.Unmarshal(p, &str) == nil {
			sb.WriteString(str)
			continue
		}
		var obj struct {
			Text string `json:"text"`
		}
		if json.Unmarshal(p, &obj) == nil {
			sb.WriteString(obj.Text)
		}
	}
	return strings.TrimSpace(sb.String())
}
