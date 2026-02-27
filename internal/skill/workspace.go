package skill

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const workspaceSkillsDir = ".github/skills"

// Workspace prints all skills installed in the current workspace.
func Workspace() error {
	entries, err := os.ReadDir(workspaceSkillsDir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			fmt.Println("No skills installed.")
			return nil
		}
		return fmt.Errorf("reading skills directory: %w", err)
	}

	var skills []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if _, err := os.Stat(filepath.Join(workspaceSkillsDir, e.Name(), "SKILL.md")); err != nil {
			continue
		}
		skills = append(skills, e.Name())
	}

	if len(skills) == 0 {
		fmt.Println("No skills installed.")
		return nil
	}

	for _, name := range skills {
		fmt.Println(name)
	}

	return nil
}
