package skill

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/otakakot/asm/internal/manifest"
)

// List prints all skills downloaded to the local cache (~/.asm/skills/).
func List() error {
	base, err := localSkillsDir()
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(base)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			fmt.Println("No skills downloaded.")
			return nil
		}
		return fmt.Errorf("reading skills directory: %w", err)
	}

	man, _ := manifest.Load()

	var found bool
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if _, err := os.Stat(filepath.Join(base, e.Name(), "SKILL.md")); err != nil {
			continue
		}
		found = true
		if s := man.Find(e.Name()); s != nil {
			fmt.Printf("%s\t%s\n", e.Name(), s.Source)
		} else {
			fmt.Println(e.Name())
		}
	}

	if !found {
		fmt.Println("No skills downloaded.")
	}

	return nil
}
