package manifest_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otakakot/asm/internal/manifest"
)

func TestLoad_NoFile(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	man, err := manifest.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(man.Skills) != 0 {
		t.Fatalf("expected empty skills, got %d", len(man.Skills))
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	p := filepath.Join(home, ".asm", "asm.json")
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(p, []byte("{invalid"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := manifest.Load()
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestSave_RoundTrip(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	man := &manifest.Manifest{
		Skills: []manifest.Skill{
			{Name: "foo", Source: "owner/repo"},
		},
	}

	if err := man.Save(); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded, err := manifest.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if len(loaded.Skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(loaded.Skills))
	}

	if loaded.Skills[0].Name != "foo" || loaded.Skills[0].Source != "owner/repo" {
		t.Fatalf("unexpected skill: %+v", loaded.Skills[0])
	}
}

func TestAdd_NewAndUpdate(t *testing.T) {
	man := &manifest.Manifest{}

	man.Add("a", "source-a")
	if len(man.Skills) != 1 {
		t.Fatalf("expected 1, got %d", len(man.Skills))
	}

	man.Add("a", "source-a-updated")
	if len(man.Skills) != 1 {
		t.Fatalf("expected 1 after update, got %d", len(man.Skills))
	}

	if man.Skills[0].Source != "source-a-updated" {
		t.Fatalf("expected updated source, got %s", man.Skills[0].Source)
	}

	man.Add("b", "source-b")
	if len(man.Skills) != 2 {
		t.Fatalf("expected 2, got %d", len(man.Skills))
	}
}

func TestRemove(t *testing.T) {
	man := &manifest.Manifest{
		Skills: []manifest.Skill{
			{Name: "a", Source: "s-a"},
			{Name: "b", Source: "s-b"},
		},
	}

	man.Remove("a")
	if len(man.Skills) != 1 || man.Skills[0].Name != "b" {
		t.Fatalf("unexpected skills after remove: %+v", man.Skills)
	}

	man.Remove("nonexistent")
	if len(man.Skills) != 1 {
		t.Fatalf("remove nonexistent changed count: %d", len(man.Skills))
	}
}

func TestFind(t *testing.T) {
	man := &manifest.Manifest{
		Skills: []manifest.Skill{
			{Name: "x", Source: "s-x"},
		},
	}

	if s := man.Find("x"); s == nil || s.Source != "s-x" {
		t.Fatalf("expected to find x, got %v", s)
	}

	if s := man.Find("y"); s != nil {
		t.Fatalf("expected nil, got %+v", s)
	}
}
