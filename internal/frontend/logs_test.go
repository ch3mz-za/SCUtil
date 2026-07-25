package frontend

import (
	"testing"

	"github.com/ch3mz-za/SCUtil/internal/logmon"
)

func TestStringToLogItem_ShortLine(t *testing.T) {
	if got := stringToLogItem("malformed"); got != nil {
		t.Fatalf("stringToLogItem returned %+v for malformed input", got)
	}
}

func TestFilterMatch_Search(t *testing.T) {
	item := logItemFixture()
	if !filterMatch(item, Search, "attacker") {
		t.Fatal("search should match attacker")
	}
	if filterMatch(item, Search, "missing") {
		t.Fatal("search should not match unrelated text")
	}
}

func logItemFixture() *logmon.LogItem {
	return &logmon.LogItem{Attacker: "Attacker", Victim: "Victim", Type: logmon.ActorDeath}
}
