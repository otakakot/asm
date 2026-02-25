# asm - Agent Skills Manager

A CLI tool for installing agent skills published on GitHub repositories into your local workspace.

## Installation

```sh
go install github.com/otakakot/asm@latest
```

## Usage

```
asm install <github_repo_url_path>   Install a skill from GitHub
asm list                             List installed skills
asm remove <skill_name>              Remove an installed skill
asm history                          Show global installation history
```

### Install a skill

Specify a GitHub repository path to install a skill. You can include `tree/<branch>` to specify a branch.

```sh
asm install anthropics/skills/tree/main/skills/skill-creator
```

Skills are installed into the `.github/skills/<skill_name>/` directory.

> The target path must contain a `SKILL.md` file, otherwise the install will fail.

### List installed skills

```sh
asm list
```

Lists directories under `.github/skills/`.

### Remove a skill

```sh
asm remove skill-creator
```

### View installation history

```sh
asm history
```

Displays all previously installed skills. History is stored globally at `~/.asm/history.json`.

## References

- [Agent Skills (Claude)](https://platform.claude.com/docs/ja/agents-and-tools/agent-skills/overview)
- [Agent Skills (GitHub Copilot)](https://docs.github.com/ja/copilot/concepts/agents/about-agent-skills)
