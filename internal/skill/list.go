package skill

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

// List prints all installed skills by scanning the .github/skills directory.
func List() error {
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			fmt.Println("No skills installed.")
			return nil
		}
		return fmt.Errorf("reading skills directory: %w", err)
	}

	var skills []string
	for _, e := range entries {
		if e.IsDir() {
			skills = append(skills, e.Name())
		}
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
