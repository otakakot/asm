package skill_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/otakakot/asm/internal/skill"
)

func TestList_NoDir(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := skill.List()

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	if got := buf.String(); got != "No skills downloaded.\n" {
		t.Errorf("output = %q, want %q", got, "No skills downloaded.\n")
	}
}

func TestList_WithSkills(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	dir := filepath.Join(tmp, ".asm", "skills")
	os.MkdirAll(filepath.Join(dir, "alpha"), 0o755)
	os.MkdirAll(filepath.Join(dir, "beta"), 0o755)
	os.WriteFile(filepath.Join(dir, "alpha", "SKILL.md"), []byte("# Alpha"), 0o644)
	os.WriteFile(filepath.Join(dir, "beta", "SKILL.md"), []byte("# Beta"), 0o644)
	// no SKILL.md -> should be ignored
	os.MkdirAll(filepath.Join(dir, "no-skill-md"), 0o755)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := skill.List()

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	got := buf.String()

	if want := "alpha\nbeta\n"; got != want {
		t.Errorf("output = %q, want %q", got, want)
	}
}

func TestList_EmptyDir(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	os.MkdirAll(filepath.Join(tmp, ".asm", "skills"), 0o755)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := skill.List()

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	if got := buf.String(); got != "No skills downloaded.\n" {
		t.Errorf("output = %q, want %q", got, "No skills downloaded.\n")
	}
}
