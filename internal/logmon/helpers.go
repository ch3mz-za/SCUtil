package logmon

import (
	"strings"
	"time"
)

func roundTimeToSeconds(t string) string {
	parsedTime, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return "parse error"
	}
	return parsedTime.Round(time.Second).Format(time.RFC3339)
}

func trimAngleBrackets(tok string) string {
	if tok == "" {
		return ""
	}
	tok = strings.TrimPrefix(tok, "<")
	return strings.TrimSuffix(tok, ">")
}

func trimQuotes(tok string) string {
	if tok == "" {
		return ""
	}
	tok = strings.TrimPrefix(tok, "'")
	return strings.TrimSuffix(tok, "'")
}

func normalize(tok string) string {
	tok = trimQuotes(tok)
	if tok == "" {
		return ""
	}
	parts := strings.Split(tok, "_")
	if len(parts) > 1 && len(parts[len(parts)-1]) >= 12 {
		parts = parts[:len(parts)-1]
	}
	return strings.Join(parts, "_")
}
