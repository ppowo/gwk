package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
)

// Repo represents a configured repository with its generated directory name.
type Repo struct {
	URL           string   `json:"url"`
	Branch        string   `json:"branch"`
	Tags          []string `json:"tags"`
	DirectoryName string   `json:"-"`
}

type jsonRepo struct {
	URL    string   `json:"url"`
	Branch string   `json:"branch"`
	Tags   []string `json:"tags"`
}

// Load reads and validates the config at the given path.
// Returns a slice of repos with sanitized tags and pre-computed directory names.
func Load(configPath string) ([]Repo, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var rawRepos []jsonRepo
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&rawRepos); err != nil {
		return nil, fmt.Errorf("invalid JSON in config: %w", err)
	}

	seenDirs := make(map[string]int)
	repos := make([]Repo, 0, len(rawRepos))

	for i, raw := range rawRepos {
		repo, err := validateAndConvert(i, raw)
		if err != nil {
			return nil, err
		}
		if firstIdx, ok := seenDirs[repo.DirectoryName]; ok {
			return nil, fmt.Errorf("config entry %d produces duplicate directory name (same as entry %d): %s", i+1, firstIdx+1, repo.DirectoryName)
		}
		seenDirs[repo.DirectoryName] = i
		repos = append(repos, repo)
	}

	return repos, nil
}

func validateAndConvert(idx int, raw jsonRepo) (Repo, error) {
	if raw.URL == "" {
		return Repo{}, fmt.Errorf("config entry %d: url is required", idx+1)
	}
	if raw.Branch == "" {
		return Repo{}, fmt.Errorf("config entry %d: branch is required", idx+1)
	}
	if len(raw.Tags) == 0 {
		return Repo{}, fmt.Errorf("config entry %d: tags must contain at least one tag", idx+1)
	}

	sanitized := make([]string, 0, len(raw.Tags))
	for j, tag := range raw.Tags {
		s := sanitizeTag(tag)
		if s == "" {
			return Repo{}, fmt.Errorf("config entry %d, tag %d: tag is empty after sanitization", idx+1, j+1)
		}
		sanitized = append(sanitized, s)
	}

	repoName := extractRepoName(raw.URL)
	sanitizedBranch := sanitizeBranch(raw.Branch)
	dirName := fmt.Sprintf("%s-%s-%s", repoName, sanitizedBranch, strings.Join(sanitized, "-"))

	return Repo{
		URL:           raw.URL,
		Branch:        raw.Branch,
		Tags:          sanitized,
		DirectoryName: dirName,
	}, nil
}

func sanitizeTag(s string) string {
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "/", "")
	return s
}

func sanitizeBranch(s string) string {
	return strings.ReplaceAll(s, "/", "-")
}

func extractRepoName(url string) string {
	// Remove trailing .git if present
	url = strings.TrimSuffix(url, ".git")
	// Take the last path segment
	return path.Base(url)
}
