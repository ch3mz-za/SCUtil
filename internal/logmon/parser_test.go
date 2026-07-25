package logmon

import "testing"

func TestGameParser_SkipsBlankAndShortLines(t *testing.T) {
	parse := GameParser(DefaultAD(), DefaultVD())

	for _, line := range []string{"", "   ", "Actor Death"} {
		if item, ok := parse(line); ok || item != nil {
			t.Fatalf("parse(%q) = (%v, %v), want no item", line, item, ok)
		}
	}
}

func TestParseMode(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  Mode
	}{
		{"once", ModeOnce},
		{"TAIL", ModeTail},
		{"Both", ModeBoth},
	} {
		got, err := ParseMode(tc.input)
		if err != nil || got != tc.want {
			t.Fatalf("ParseMode(%q) = %v, %v; want %v", tc.input, got, err, tc.want)
		}
	}
	if _, err := ParseMode("invalid"); err == nil {
		t.Fatal("ParseMode should reject invalid modes")
	}
}
