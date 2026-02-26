package skill_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otakakot/asm/internal/skill"
)

func TestRemove_Success(t *testing.T) {
	tmp := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmp)

	dir := filepath.Join(tmp, ".github", "skills", "my-skill")
	os.MkdirAll(dir, 0o755)
	os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# Skill"), 0o644)

	if err := skill.Remove("my-skill"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Error("skill directory should have been removed")
	}
}

func TestRemove_NotInstalled(t *testing.T) {
	tmp := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmp)

	err := skill.Remove("nonexistent")
	if err == nil {
		t.Fatal("expected error for non-existent skill, got nil")
	}
}

func TestRemove_FileNotDir(t *testing.T) {
	tmp := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmp)

	dir := filepath.Join(tmp, ".github", "skills")
	os.MkdirAll(dir, 0o755)
	os.WriteFile(filepath.Join(dir, "not-a-dir"), []byte("x"), 0o644)

	err := skill.Remove("not-a-dir")
	if err == nil {
		t.Fatal("expected error when target is a file, got nil")
	}
}
