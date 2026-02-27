package skill_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otakakot/asm/internal/skill"
)

func TestInstall_Success(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	workspace := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(workspace)

	srcDir := filepath.Join(tmp, ".asm", "skills", "my-skill")
	os.MkdirAll(srcDir, 0o755)
	os.WriteFile(filepath.Join(srcDir, "SKILL.md"), []byte("# Skill"), 0o644)
	os.WriteFile(filepath.Join(srcDir, "prompt.md"), []byte("prompt"), 0o644)

	if err := skill.Install("my-skill"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	destSkill := filepath.Join(workspace, ".github", "skills", "my-skill", "SKILL.md")
	if _, err := os.Stat(destSkill); err != nil {
		t.Errorf("SKILL.md not found in workspace: %v", err)
	}

	destPrompt := filepath.Join(workspace, ".github", "skills", "my-skill", "prompt.md")
	if _, err := os.Stat(destPrompt); err != nil {
		t.Errorf("prompt.md not found in workspace: %v", err)
	}
}

func TestInstall_NotDownloaded(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	tmp := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmp)

	err := skill.Install("nonexistent")
	if err == nil {
		t.Fatal("expected error for non-downloaded skill, got nil")
	}
}
