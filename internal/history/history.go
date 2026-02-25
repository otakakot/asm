package history

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Entry represents a previously installed skill.
type Entry struct {
	Name   string `json:"name"`
	Source string `json:"source"`
}

// History holds all history entries.
type History struct {
	Entries []Entry `json:"entries"`
}

func dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	return filepath.Join(home, ".asm"), nil
}

func filePath() (string, error) {
	d, err := dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, "history.json"), nil
}

// Load reads the global history file.
func Load() (*History, error) {
	path, err := filePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return &History{}, nil
		}
		return nil, fmt.Errorf("reading history: %w", err)
	}

	var h History
	if err := json.Unmarshal(data, &h); err != nil {
		return nil, fmt.Errorf("parsing history: %w", err)
	}

	return &h, nil
}

// save writes the history to the global history file.
func (h *History) save() error {
	path, err := filePath()
	if err != nil {
		return err
	}

	d := filepath.Dir(path)
	if err := os.MkdirAll(d, 0o755); err != nil {
		return fmt.Errorf("creating directory %s: %w", d, err)
	}

	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling history: %w", err)
	}
	data = append(data, '\n')

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing history: %w", err)
	}

	return nil
}

// Record adds a skill to the history if it is not already present.
func Record(name, source string) error {
	h, err := Load()
	if err != nil {
		return err
	}

	for _, e := range h.Entries {
		if e.Name == name {
			return nil // already recorded
		}
	}

	h.Entries = append(h.Entries, Entry{
		Name:   name,
		Source: source,
	})

	return h.save()
}
