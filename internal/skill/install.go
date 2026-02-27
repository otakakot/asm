package skill

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Install copies a downloaded skill from the local cache to the workspace.
func Install(name string) error {
	base, err := localSkillsDir()
	if err != nil {
		return err
	}

	srcDir := filepath.Join(base, name)

	if _, err := os.Stat(filepath.Join(srcDir, "SKILL.md")); err != nil {
		return fmt.Errorf("skill %q is not downloaded", name)
	}

	destDir := filepath.Join(workspaceSkillsDir, name)

	if err := os.MkdirAll(filepath.Dir(destDir), 0o755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	if err := copyDir(srcDir, destDir); err != nil {
		return fmt.Errorf("copying skill: %w", err)
	}

	fmt.Printf("Skill %q installed to workspace.\n", name)

	return nil
}

// copyDir recursively copies src directory to dst.
func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		target := filepath.Join(dst, rel)

		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}

		return os.WriteFile(target, data, 0o644)
	})
}
