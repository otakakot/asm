package skill

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/otakakot/asm/internal/github"
	"github.com/otakakot/asm/internal/history"
)

const skillsDir = ".github/skills"

// Install fetches a skill from GitHub and installs it into the workspace.
func Install(ctx context.Context, rawPath string) error {
	ref, err := github.ParseRepoPath(rawPath)
	if err != nil {
		return err
	}

	name := ref.SkillName()
	fmt.Printf("Installing skill %q from %s/%s (branch: %s, path: %s)...\n",
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

	destDir := filepath.Join(skillsDir, name)

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

	if err := history.Record(name, rawPath); err != nil {
		return fmt.Errorf("recording history: %w", err)
	}

	fmt.Printf("Skill %q installed successfully.\n", name)

	return nil
}
