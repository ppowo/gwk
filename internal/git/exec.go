package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Exec is the real git adapter that shells out to the git binary.
type Exec struct {
	MirrorDir string // e.g. ~/code-mirror (resolved to absolute)
}

func NewExec(mirrorDir string) (*Exec, error) {
	abs, err := filepath.Abs(mirrorDir)
	if err != nil {
		return nil, fmt.Errorf("resolving mirror directory: %w", err)
	}
	return &Exec{MirrorDir: abs}, nil
}

func (e *Exec) Clone(url, branch, dst string) error {
	if err := e.validateInMirror(dst); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("creating parent directory for clone: %w", err)
	}
	cmd := exec.Command("git", "clone", "--branch", branch, url, dst)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %w (stderr: %s)", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

func (e *Exec) UpdateToLatest(repoPath, branch string) error {
	if err := e.validateInMirror(repoPath); err != nil {
		return err
	}

	commands := [][]string{
		{"git", "-C", repoPath, "fetch", "origin"},
		{"git", "-C", repoPath, "checkout", branch},
		{"git", "-C", repoPath, "reset", "--hard", "origin/" + branch},
	}

	for _, args := range commands {
		cmd := exec.Command(args[0], args[1:]...)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("git %s failed: %w (stderr: %s)", strings.Join(args[1:], " "), err, strings.TrimSpace(stderr.String()))
		}
	}
	return nil
}

func (e *Exec) RemoteHead(url, branch string) (string, error) {
	cmd := exec.Command("git", "ls-remote", url, "refs/heads/"+branch)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git ls-remote failed: %w (stderr: %s)", err, strings.TrimSpace(stderr.String()))
	}

	output := strings.TrimSpace(stdout.String())
	if output == "" {
		return "", fmt.Errorf("remote branch %q not found at %s", branch, url)
	}
	// Format: "<sha>\trefs/heads/<branch>"
	parts := strings.SplitN(output, "\t", 2)
	if len(parts) < 1 {
		return "", fmt.Errorf("unexpected git ls-remote output: %q", output)
	}
	return parts[0], nil
}

func (e *Exec) LocalHead(repoPath string) (string, error) {
	if err := e.validateInMirror(repoPath); err != nil {
		return "", err
	}
	cmd := exec.Command("git", "-C", repoPath, "rev-parse", "HEAD")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git rev-parse failed: %w (stderr: %s)", err, strings.TrimSpace(stderr.String()))
	}
	return strings.TrimSpace(stdout.String()), nil
}

func (e *Exec) validateInMirror(target string) error {
	abs, err := filepath.Abs(target)
	if err != nil {
		return fmt.Errorf("resolving path: %w", err)
	}
	// Ensure the target is a child of the mirror directory, not merely a string prefix
	// like /home/me/code-mirror-evil.
	rel, err := filepath.Rel(e.MirrorDir, abs)
	if err != nil {
		return fmt.Errorf("checking mirror-relative path: %w", err)
	}
	if rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("path %q is not inside mirror directory %q", abs, e.MirrorDir)
	}
	return nil
}
