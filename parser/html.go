// Package parser разбирает HTML и JSON экспорты чата Telegram и возвращает
// типизированные уведомления из пакета notification.
package parser

import (
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// ParseHTML разбирает HTML-экспорт Telegram.
// Возвращает уведомления с HTMLMessageID > lastID, новый максимальный ID и ошибку.
func ParseHTML(r io.Reader, lastID int) ([]interface{}, int, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, lastID, err
	}

	raw := collectMessages(doc)
	sort.Slice(raw, func(i, j int) bool { return raw[i].HTMLID < raw[j].HTMLID })

	maxID := lastID
	var notifications []interface{}

	for _, m := range raw {
		if m.HTMLID <= lastID {
			continue
		}
		if n := buildNotification(m); n != nil {
			notifications = append(notifications, n)
		}
		if m.HTMLID > maxID {
			maxID = m.HTMLID
		}
	}

	return notifications, maxID, nil
}

func collectMessages(doc *html.Node) []rawMessage {
	var messages []rawMessage
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if isMessageDiv(n) {
			if msg, ok := parseMessageNode(n); ok {
				messages = append(messages, msg)
			}
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	return messages
}

func isMessageDiv(n *html.Node) bool {
	if n.Type != html.ElementNode || n.Data != "div" {
		return false
	}
	cls := getAttr(n, "class")
	id := getAttr(n, "id")
	return (cls == "message default clearfix" || cls == "message default clearfix joined") &&
		strings.HasPrefix(id, "message")
}

func parseMessageNode(n *html.Node) (rawMessage, bool) {
	htmlID, err := strconv.Atoi(strings.TrimPrefix(getAttr(n, "id"), "message"))
	if err != nil {
		return rawMessage{}, false
	}

	var ts time.Time
	var text string

	var find func(*html.Node)
	find = func(child *html.Node) {
		if child.Type == html.ElementNode && child.Data == "div" {
			switch getAttr(child, "class") {
			case "pull_right date details":
				ts, _ = time.Parse("02.01.2006 15:04:05 UTC-07:00", getAttr(child, "title"))
			case "text":
				text = extractText(child)
				return
			}
		}
		for c := child.FirstChild; c != nil; c = c.NextSibling {
			find(c)
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		find(c)
	}

	if text == "" {
		return rawMessage{}, false
	}
	return rawMessage{HTMLID: htmlID, Timestamp: ts, Text: text}, true
}

func extractText(n *html.Node) string {
	var sb strings.Builder
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		switch {
		case node.Type == html.TextNode:
			sb.WriteString(node.Data)
		case node.Type == html.ElementNode && node.Data == "br":
			sb.WriteString("\n")
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return strings.TrimSpace(sb.String())
}

func getAttr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}
