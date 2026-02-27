package skill

import (
	"fmt"
	"os"
	"path/filepath"
)

// Remove deletes an installed skill from the workspace.
func Remove(name string) error {
	dir := filepath.Join(workspaceSkillsDir, name)

	info, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("skill %q is not installed", name)
	}

	if !info.IsDir() {
		return fmt.Errorf("skill %q is not installed", name)
	}

	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("removing %s: %w", dir, err)
	}

	fmt.Printf("Skill %q removed.\n", name)

	return nil
}
