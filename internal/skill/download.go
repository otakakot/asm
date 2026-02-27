package skill

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/otakakot/asm/internal/github"
	"github.com/otakakot/asm/internal/manifest"
)

// localSkillsDir returns the path to the global skills cache directory.
func localSkillsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	return filepath.Join(home, ".asm", "skills"), nil
}

// Download fetches a skill from GitHub and saves it to the local cache (~/.asm/skills/).
func Download(ctx context.Context, rawPath string) error {
	ref, err := github.ParseRepoPath(rawPath)
	if err != nil {
		return err
	}

	name := ref.SkillName()
	fmt.Printf("Downloading skill %q from %s/%s (branch: %s, path: %s)...\n",
		name, ref.Owner, ref.Repo, ref.Branch, ref.Path)

	client := github.NewClient(nil)

	// Check that the target directory contains a SKILL.md file.
	entries, err := client.ListContents(ctx, ref)
	if err != nil {
		return fmt.Errorf("listing contents: %w", err)
	}

	hasSkillMD := false
	for _, e := range entries {
		if e.Type == "file" && e.Name == "SKILL.md" {
			hasSkillMD = true
			break
		}
	}

	if !hasSkillMD {
		return fmt.Errorf("SKILL.md not found at %s: not a valid skill", rawPath)
	}

	files, err := client.FetchAllFiles(ctx, ref)
	if err != nil {
		return fmt.Errorf("fetching skill files: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no files found at %s", rawPath)
	}

	base, err := localSkillsDir()
	if err != nil {
		return err
	}

	destDir := filepath.Join(base, name)

	for relPath, data := range files {
		dest := filepath.Join(destDir, relPath)

		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return fmt.Errorf("creating directory for %s: %w", dest, err)
		}

		if err := os.WriteFile(dest, data, 0o644); err != nil {
			return fmt.Errorf("writing %s: %w", dest, err)
		}

		fmt.Printf("  %s\n", dest)
	}

	fmt.Printf("Skill %q downloaded successfully.\n", name)

	man, err := manifest.Load()
	if err != nil {
		return fmt.Errorf("loading manifest: %w", err)
	}

	man.Add(name, rawPath)

	if err := man.Save(); err != nil {
		return fmt.Errorf("saving manifest: %w", err)
	}

	return nil
}
