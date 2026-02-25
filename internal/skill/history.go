package skill

import (
	"fmt"

	"github.com/otakakot/asm/internal/history"
)

// History prints the global installation/removal history.
func History() error {
	h, err := history.Load()
	if err != nil {
		return err
	}

	if len(h.Entries) == 0 {
		fmt.Println("No history.")
		return nil
	}

	fmt.Printf("%-20s %s\n", "NAME", "SOURCE")
	for _, e := range h.Entries {
		fmt.Printf("%-20s %s\n", e.Name, e.Source)
	}

	return nil
}
