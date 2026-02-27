package skill

import (
	"fmt"
	"os"
	"path/filepath"
)

// Link creates a symbolic link from the local skill cache to the workspace.
func Link(name string) error {
	base, err := localSkillsDir()
	if err != nil {
		return err
	}

	srcDir := filepath.Join(base, name)

	if _, err := os.Stat(filepath.Join(srcDir, "SKILL.md")); err != nil {
		return fmt.Errorf("skill %q is not downloaded", name)
	}

	if err := os.MkdirAll(workspaceSkillsDir, 0o755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	linkPath := filepath.Join(workspaceSkillsDir, name)

	if _, err := os.Lstat(linkPath); err == nil {
		return fmt.Errorf("skill %q already exists in workspace", name)
	}

	if err := os.Symlink(srcDir, linkPath); err != nil {
		return fmt.Errorf("creating symlink: %w", err)
	}

	fmt.Printf("Skill %q linked to workspace.\n", name)

	return nil
}
