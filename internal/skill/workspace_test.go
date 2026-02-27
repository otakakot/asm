package skill_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/otakakot/asm/internal/skill"
)

func TestWorkspace_NoDir(t *testing.T) {
	tmp := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmp)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := skill.Workspace()

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

func TestWorkspace_WithSkills(t *testing.T) {
	tmp := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmp)

	dir := filepath.Join(tmp, ".github", "skills")
	os.MkdirAll(filepath.Join(dir, "alpha"), 0o755)
	os.MkdirAll(filepath.Join(dir, "beta"), 0o755)
	os.WriteFile(filepath.Join(dir, "alpha", "SKILL.md"), []byte("# Alpha"), 0o644)
	os.WriteFile(filepath.Join(dir, "beta", "SKILL.md"), []byte("# Beta"), 0o644)
	os.MkdirAll(filepath.Join(dir, "no-skill-md"), 0o755)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := skill.Workspace()

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	if want := "alpha\nbeta\n"; buf.String() != want {
		t.Errorf("output = %q, want %q", buf.String(), want)
	}
}
