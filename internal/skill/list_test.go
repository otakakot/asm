package skill_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/otakakot/asm/internal/skill"
)

func TestList_NoDir(t *testing.T) {
	tmp := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmp)

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
	if got := buf.String(); got != "No skills installed.\n" {
		t.Errorf("output = %q, want %q", got, "No skills installed.\n")
	}
}

func TestList_WithSkills(t *testing.T) {
	tmp := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmp)

	dir := filepath.Join(tmp, ".github", "skills")
	os.MkdirAll(filepath.Join(dir, "alpha"), 0o755)
	os.MkdirAll(filepath.Join(dir, "beta"), 0o755)
	// SKILL.md を持つディレクトリのみスキルとみなされる
	os.WriteFile(filepath.Join(dir, "alpha", "SKILL.md"), []byte("# Alpha"), 0o644)
	os.WriteFile(filepath.Join(dir, "beta", "SKILL.md"), []byte("# Beta"), 0o644)
	// SKILL.md がないディレクトリとファイルは無視される
	os.MkdirAll(filepath.Join(dir, "no-skill-md"), 0o755)
	os.WriteFile(filepath.Join(dir, "not-a-skill.txt"), []byte("x"), 0o644)

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
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmp)

	os.MkdirAll(filepath.Join(tmp, ".github", "skills"), 0o755)

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
	if got := buf.String(); got != "No skills installed.\n" {
		t.Errorf("output = %q, want %q", got, "No skills installed.\n")
	}
}
