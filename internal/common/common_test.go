package common

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindDir_FindsTargetAndReportsMissingRoot(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "Library", "StarCitizen", "LIVE")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}

	got, err := FindDir(root, filepath.Join("StarCitizen", "LIVE"))
	if err != nil {
		t.Fatalf("FindDir returned error: %v", err)
	}
	if got != target {
		t.Fatalf("FindDir = %q, want %q", got, target)
	}

	if _, err := FindDir(filepath.Join(root, "missing"), filepath.Join("StarCitizen", "LIVE")); err == nil {
		t.Fatal("FindDir should report a missing root")
	}
}

func TestUserHomeDir_IsAvailable(t *testing.T) {
	if got := UserHomeDir(); got == "" {
		t.Fatal("UserHomeDir returned an empty path")
	}
}
