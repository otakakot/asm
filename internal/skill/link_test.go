package skill_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otakakot/asm/internal/skill"
)

func TestLink_Success(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	workspace := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(workspace)

	srcDir := filepath.Join(tmp, ".asm", "skills", "my-skill")
	os.MkdirAll(srcDir, 0o755)
	os.WriteFile(filepath.Join(srcDir, "SKILL.md"), []byte("# Skill"), 0o644)

	if err := skill.Link("my-skill"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	linkPath := filepath.Join(workspace, ".github", "skills", "my-skill")
	info, err := os.Lstat(linkPath)
	if err != nil {
		t.Fatalf("link not found: %v", err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Error("expected symlink")
	}
}

func TestLink_NotDownloaded(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	tmp := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmp)

	err := skill.Link("nonexistent")
	if err == nil {
		t.Fatal("expected error for non-downloaded skill, got nil")
	}
}

func TestLink_AlreadyExists(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	workspace := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(workspace)

	srcDir := filepath.Join(tmp, ".asm", "skills", "my-skill")
	os.MkdirAll(srcDir, 0o755)
	os.WriteFile(filepath.Join(srcDir, "SKILL.md"), []byte("# Skill"), 0o644)

	dest := filepath.Join(workspace, ".github", "skills", "my-skill")
	os.MkdirAll(dest, 0o755)

	err := skill.Link("my-skill")
	if err == nil {
		t.Fatal("expected error when skill already exists, got nil")
	}
}
