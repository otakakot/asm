# asm - Agent Skills Manager

A CLI tool for managing agent skills published on GitHub repositories.

Skills are first downloaded from GitHub to a local cache (`~/.asm/skills/`), then installed or symlinked into your workspace (`.github/skills/`).

## Installation

```sh
go install github.com/otakakot/asm@latest
```

## Usage

```
asm download <github_repo_url_path>  Download a skill from GitHub to local cache
asm list                             List downloaded skills
asm install <skill_name>             Install a downloaded skill to workspace
asm link <skill_name>                Symlink a downloaded skill to workspace
asm workspace                        List skills installed in workspace
asm remove <skill_name>              Remove a skill from workspace
asm version                          Show version
```

### Download a skill

Specify a GitHub repository path to download a skill. You can include `tree/<branch>` to specify a branch.

```sh
asm download anthropics/skills/tree/main/skills/skill-creator
```

The skill files are saved to `~/.asm/skills/<skill_name>/` and the metadata is recorded in `~/.asm/asm.json`.

> The target path must contain a `SKILL.md` file, otherwise the download will fail.

### List downloaded skills

```sh
asm list
```

Lists skills stored in the local cache with their source repository.

### Install a skill to workspace

```sh
asm install skill-creator
```

Copies files from the local cache to `.github/skills/<skill_name>/`.

### Symlink a skill to workspace

```sh
asm link skill-creator
```

Creates a symbolic link from `.github/skills/<skill_name>` to the local cache.

### List workspace skills

```sh
asm workspace
```

Lists skills currently installed in `.github/skills/`.

### Remove a skill from workspace

```sh
asm remove skill-creator
```

### Show version

```sh
asm version
```

## References

- [Agent Skills (Claude)](https://platform.claude.com/docs/ja/agents-and-tools/agent-skills/overview)
- [Agent Skills (GitHub Copilot)](https://docs.github.com/ja/copilot/concepts/agents/about-agent-skills)
