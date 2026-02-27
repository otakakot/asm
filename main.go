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
  asm download <github_repo_url_path>  Download a skill from GitHub to local cache
  asm list                             List downloaded skills
  asm install <skill_name>             Install a downloaded skill to workspace
  asm link <skill_name>                Symlink a downloaded skill to workspace
  asm workspace                        List skills installed in workspace
  asm remove <skill_name>              Remove a skill from workspace
  asm version                          Show version

Examples:
  asm download anthropics/skills/tree/main/skills/skill-creator
  asm list
  asm install skill-creator
  asm link skill-creator
  asm workspace
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
	case "download":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: asm download <github_repo_url_path>")
			os.Exit(1)
		}
		err = skill.Download(ctx, os.Args[2])

	case "list":
		err = skill.List()

	case "install":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: asm install <skill_name>")
			os.Exit(1)
		}
		err = skill.Install(os.Args[2])

	case "link":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: asm link <skill_name>")
			os.Exit(1)
		}
		err = skill.Link(os.Args[2])

	case "workspace":
		err = skill.Workspace()

	case "remove":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: asm remove <skill_name>")
			os.Exit(1)
		}
		err = skill.Remove(os.Args[2])

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
