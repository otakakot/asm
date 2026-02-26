package github_test

import (
	"testing"

	"github.com/otakakot/asm/internal/github"
)

func TestParseRepoPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    github.RepoRef
		wantErr bool
	}{
		{
			name:  "owner/repo only",
			input: "owner/repo",
			want:  github.RepoRef{Owner: "owner", Repo: "repo", Branch: "main"},
		},
		{
			name:  "with tree/branch/path",
			input: "anthropics/skills/tree/main/skills/skill-creator",
			want:  github.RepoRef{Owner: "anthropics", Repo: "skills", Branch: "main", Path: "skills/skill-creator"},
		},
		{
			name:  "with non-main branch",
			input: "owner/repo/tree/develop/some/path",
			want:  github.RepoRef{Owner: "owner", Repo: "repo", Branch: "develop", Path: "some/path"},
		},
		{
			name:  "https URL",
			input: "https://github.com/owner/repo/tree/main/dir",
			want:  github.RepoRef{Owner: "owner", Repo: "repo", Branch: "main", Path: "dir"},
		},
		{
			name:  "http URL",
			input: "http://github.com/owner/repo",
			want:  github.RepoRef{Owner: "owner", Repo: "repo", Branch: "main"},
		},
		{
			name:  "github.com prefix without scheme",
			input: "github.com/owner/repo",
			want:  github.RepoRef{Owner: "owner", Repo: "repo", Branch: "main"},
		},
		{
			name:  "trailing slash",
			input: "owner/repo/tree/main/path/",
			want:  github.RepoRef{Owner: "owner", Repo: "repo", Branch: "main", Path: "path"},
		},
		{
			name:    "too short",
			input:   "owner",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := github.ParseRepoPath(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestRepoRef_SkillName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ref  github.RepoRef
		want string
	}{
		{
			name: "path with multiple segments",
			ref:  github.RepoRef{Repo: "skills", Path: "skills/skill-creator"},
			want: "skill-creator",
		},
		{
			name: "single segment path",
			ref:  github.RepoRef{Repo: "skills", Path: "myskill"},
			want: "myskill",
		},
		{
			name: "no path falls back to repo",
			ref:  github.RepoRef{Repo: "my-skill"},
			want: "my-skill",
		},
		{
			name: "trailing slash in path",
			ref:  github.RepoRef{Repo: "repo", Path: "a/b/c/"},
			want: "c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.ref.SkillName()
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
