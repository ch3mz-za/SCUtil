package logmon

import (
	"strings"
	"time"

	"github.com/ch3mz-za/SCUtil/internal/util"
)

func RoundTimeToSeconds(t string) *time.Time {
	parsedTime, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return nil
	}
	return util.Ptr(parsedTime.Round(time.Second))
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
