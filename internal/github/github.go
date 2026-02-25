package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const apiBase = "https://api.github.com"

// RepoRef holds parsed GitHub repository reference information.
type RepoRef struct {
	Owner  string
	Repo   string
	Branch string
	Path   string
}

// ContentEntry represents a single item returned by the GitHub Contents API.
type ContentEntry struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	DownloadURL string `json:"download_url"`
}

// ParseRepoPath parses a GitHub URL path into a RepoRef.
//
// Accepted formats:
//
//	owner/repo
//	owner/repo/tree/branch/path/to/skill
func ParseRepoPath(raw string) (RepoRef, error) {
	raw = strings.TrimPrefix(raw, "https://github.com/")
	raw = strings.TrimPrefix(raw, "http://github.com/")
	raw = strings.TrimPrefix(raw, "github.com/")
	raw = strings.Trim(raw, "/")

	parts := strings.Split(raw, "/")
	if len(parts) < 2 {
		return RepoRef{}, fmt.Errorf("invalid repo path %q: need at least owner/repo", raw)
	}

	ref := RepoRef{
		Owner:  parts[0],
		Repo:   parts[1],
		Branch: "main",
	}

	rest := parts[2:]
	if len(rest) >= 2 && rest[0] == "tree" {
		ref.Branch = rest[1]
		rest = rest[2:]
	}

	if len(rest) > 0 {
		ref.Path = strings.Join(rest, "/")
	}

	return ref, nil
}

// SkillName returns the directory name to use for the installed skill.
func (r RepoRef) SkillName() string {
	if r.Path != "" {
		parts := strings.Split(strings.TrimRight(r.Path, "/"), "/")
		return parts[len(parts)-1]
	}
	return r.Repo
}

// Client is a thin wrapper around the GitHub REST API.
type Client struct {
	httpClient *http.Client
}

// NewClient returns a Client that uses the given http.Client.
func NewClient(hc *http.Client) *Client {
	if hc == nil {
		hc = http.DefaultClient
	}
	return &Client{httpClient: hc}
}

// ListContents returns the directory entries at the given path.
func (c *Client) ListContents(ctx context.Context, ref RepoRef) ([]ContentEntry, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s?ref=%s",
		apiBase, ref.Owner, ref.Repo, ref.Path, ref.Branch)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error %d: %s", resp.StatusCode, body)
	}

	var entries []ContentEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return entries, nil
}

// downloadFile downloads the raw content of a file from its download URL.
func (c *Client) downloadFile(ctx context.Context, downloadURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download error %d for %s", resp.StatusCode, downloadURL)
	}

	return io.ReadAll(resp.Body)
}

// FetchAllFiles recursively fetches all files under the given RepoRef.
func (c *Client) FetchAllFiles(ctx context.Context, ref RepoRef) (map[string][]byte, error) {
	files := make(map[string][]byte)
	if err := c.fetchRecursive(ctx, ref, ref.Path, files); err != nil {
		return nil, err
	}
	return files, nil
}

func (c *Client) fetchRecursive(ctx context.Context, ref RepoRef, currentPath string, files map[string][]byte) error {
	subRef := RepoRef{
		Owner:  ref.Owner,
		Repo:   ref.Repo,
		Branch: ref.Branch,
		Path:   currentPath,
	}

	entries, err := c.ListContents(ctx, subRef)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		switch entry.Type {
		case "file":
			data, err := c.downloadFile(ctx, entry.DownloadURL)
			if err != nil {
				return fmt.Errorf("downloading %s: %w", entry.Path, err)
			}
			relPath := strings.TrimPrefix(entry.Path, ref.Path+"/")
			if relPath == entry.Path {
				relPath = entry.Name
			}
			files[relPath] = data

		case "dir":
			if err := c.fetchRecursive(ctx, ref, entry.Path, files); err != nil {
				return err
			}
		}
	}

	return nil
}
