package frontend

import (
	"fmt"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ch3mz-za/SCUtil/internal/logmon"
)

// Helper functions for creating rich text segments
func newTextSegment(text string) *widget.TextSegment {
	return &widget.TextSegment{
		Text:  text,
		Style: widget.RichTextStyle{Inline: true},
	}
}

func newBoldSegment(text string) *widget.TextSegment {
	return &widget.TextSegment{
		Text: text,
		Style: widget.RichTextStyle{
			TextStyle: fyne.TextStyle{Bold: true},
			Inline:    true,
		},
	}
}

func newColoredBoldSegment(text string, color fyne.ThemeColorName) *widget.TextSegment {
	return &widget.TextSegment{
		Text: text,
		Style: widget.RichTextStyle{
			ColorName: color,
			TextStyle: fyne.TextStyle{Bold: true},
			Inline:    true,
		},
	}
}

func formatTimestamp(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Round(time.Second).Format(time.RFC3339)
}

func parseLogResult(i *logmon.LogItem) fyne.CanvasObject {
	r := widget.NewRichText()
	r.Wrapping = fyne.TextWrapWord

	switch i.Type {
	case logmon.ActorDeath:
		r.Segments = []widget.RichTextSegment{
			newTextSegment(fmt.Sprintf("%s [AD] ", formatTimestamp(i.Time))),
			newColoredBoldSegment(i.Attacker, theme.ColorNameError),
			newTextSegment(" killed "),
			newColoredBoldSegment(i.Victim, theme.ColorNamePrimary),
			newTextSegment(" using "),
			newBoldSegment(i.Weapon),
		}

	case logmon.VehicleDestruction:
		r.Segments = []widget.RichTextSegment{
			newTextSegment(fmt.Sprintf("%s [VD] ", formatTimestamp(i.Time))),
			newColoredBoldSegment(i.Attacker, theme.ColorNameError),
			newTextSegment(" destroyed "),
			newColoredBoldSegment(i.Vehicle, theme.ColorNamePrimary),
			newTextSegment(" using "),
			newBoldSegment(i.Weapon),
			newTextSegment(" at "),
			newBoldSegment(i.Location),
		}

	default:
		r.Segments = []widget.RichTextSegment{
			&widget.TextSegment{Text: fmt.Sprintf("Unknown Event: %s", i.Type)},
		}
	}

	return r
}

func generateLogName() string {
	return fmt.Sprintf("%s %s", prefixAggregateLog, time.Now().Format("2006-01-02T15.04.05.log"))
}

func getGameLogPath(selectionGameVersion *widget.Select) string {

	return filepath.Join(getGameDir(nil), selectionGameVersion.Selected, "Game.log")
}
