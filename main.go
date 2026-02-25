package main

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/otakakot/asm/internal/skill"
)

const usage = `asm - Agent Skills Manager

Usage:
  asm install <github_repo_url_path>   Install a skill from GitHub
  asm list                             List installed skills
  asm remove <skill_name>              Remove an installed skill
  asm history                          Show global installation history
  asm version                          Show version

Examples:
  asm install anthropics/skills/tree/main/skills/skill-creator
  asm list
  asm remove skill-creator
`

func main() {
	if len(os.Args) < 2 {
		fmt.Print(usage)
		os.Exit(1)
	}

	ctx := context.Background()

	var err error

	switch os.Args[1] {
	case "install":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: asm install <github_repo_url_path>")
			os.Exit(1)
		}
		err = skill.Install(ctx, os.Args[2])

	case "list":
		err = skill.List()

	case "remove":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: asm remove <skill_name>")
			os.Exit(1)
		}
		err = skill.Remove(os.Args[2])

	case "history":
		err = skill.History()

	case "version":
		version := "dev"
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
			version = info.Main.Version
		}
		fmt.Println(version)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", os.Args[1])
		fmt.Print(usage)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
