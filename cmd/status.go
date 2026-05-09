package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ppowo/gwk/internal/config"
	"github.com/ppowo/gwk/internal/git"
)

func RunStatus(repos []config.Repo, g git.Git, mirrorDir string, out io.Writer) (bool, error) {
	allUpToDate := true
	for _, repo := range repos {
		dst := filepath.Join(mirrorDir, repo.DirectoryName)
		if _, err := os.Stat(dst); os.IsNotExist(err) {
			fmt.Fprintf(out, "%s: missing\n", repo.DirectoryName)
			allUpToDate = false
			continue
		} else if err != nil {
			return false, fmt.Errorf("stat %s: %w", repo.DirectoryName, err)
		}

		remoteHead, err := g.RemoteHead(repo.URL, repo.Branch)
		if err != nil {
			return false, fmt.Errorf("remote head %s: %w", repo.DirectoryName, err)
		}

		localHead, err := g.LocalHead(dst)
		if err != nil {
			return false, fmt.Errorf("local head %s: %w", repo.DirectoryName, err)
		}

		if remoteHead == localHead {
			fmt.Fprintf(out, "%s: up to date\n", repo.DirectoryName)
		} else {
			fmt.Fprintf(out, "%s: behind\n", repo.DirectoryName)
			allUpToDate = false
		}
	}
	return allUpToDate, nil
}
