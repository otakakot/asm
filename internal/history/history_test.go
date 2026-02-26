package history_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otakakot/asm/internal/history"
)

func TestLoad_NoFile(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	h, err := history.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(h.Entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(h.Entries))
	}
}

func TestRecord_And_Load(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	if err := history.Record("skill-a", "owner/repo"); err != nil {
		t.Fatalf("Record: %v", err)
	}

	h, err := history.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(h.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(h.Entries))
	}
	if h.Entries[0].Name != "skill-a" {
		t.Errorf("name = %q, want %q", h.Entries[0].Name, "skill-a")
	}
	if h.Entries[0].Source != "owner/repo" {
		t.Errorf("source = %q, want %q", h.Entries[0].Source, "owner/repo")
	}
}

func TestRecord_Deduplication(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	if err := history.Record("skill-a", "owner/repo"); err != nil {
		t.Fatalf("first Record: %v", err)
	}
	if err := history.Record("skill-a", "owner/repo"); err != nil {
		t.Fatalf("second Record: %v", err)
	}

	h, err := history.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(h.Entries) != 1 {
		t.Errorf("expected 1 entry after dedup, got %d", len(h.Entries))
	}
}

func TestRecord_MultipleSkills(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	if err := history.Record("skill-a", "owner/repo-a"); err != nil {
		t.Fatalf("Record skill-a: %v", err)
	}
	if err := history.Record("skill-b", "owner/repo-b"); err != nil {
		t.Fatalf("Record skill-b: %v", err)
	}

	h, err := history.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(h.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(h.Entries))
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	dir := filepath.Join(tmp, ".asm")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "history.json"), []byte("{invalid"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := history.Load()
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
