package manifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Skill represents a downloaded skill entry.
type Skill struct {
	Name   string `json:"name"`
	Source string `json:"source"`
}

// Manifest holds the list of downloaded skills.
type Manifest struct {
	Skills []Skill `json:"skills"`
}

func filePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}

	return filepath.Join(home, ".asm", "asm.json"), nil
}

// Load reads the manifest from ~/.asm/asm.json.
// If the file does not exist, an empty Manifest is returned.
func Load() (*Manifest, error) {
	p, err := filePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return &Manifest{}, nil
		}

		return nil, fmt.Errorf("reading manifest: %w", err)
	}

	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parsing manifest: %w", err)
	}

	return &m, nil
}

// Save writes the manifest to ~/.asm/asm.json.
func (m *Manifest) Save() error {
	p, err := filePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling manifest: %w", err)
	}

	if err := os.WriteFile(p, data, 0o644); err != nil {
		return fmt.Errorf("writing manifest: %w", err)
	}

	return nil
}

// Add inserts or updates a skill entry by name.
func (m *Manifest) Add(name, source string) {
	for i, s := range m.Skills {
		if s.Name == name {
			m.Skills[i].Source = source
			return
		}
	}

	m.Skills = append(m.Skills, Skill{Name: name, Source: source})
}

// Remove deletes a skill entry by name.
func (m *Manifest) Remove(name string) {
	for i, s := range m.Skills {
		if s.Name == name {
			m.Skills = append(m.Skills[:i], m.Skills[i+1:]...)
			return
		}
	}
}

// Find returns a pointer to the skill entry with the given name, or nil.
func (m *Manifest) Find(name string) *Skill {
	for i, s := range m.Skills {
		if s.Name == name {
			return &m.Skills[i]
		}
	}

	return nil
}
