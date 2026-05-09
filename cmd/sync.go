package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ppowo/gwk/internal/config"
	"github.com/ppowo/gwk/internal/git"
)

func RunSync(repos []config.Repo, g git.Git, mirrorDir string, out io.Writer) error {
	for _, repo := range repos {
		dst := filepath.Join(mirrorDir, repo.DirectoryName)
		if _, err := os.Stat(dst); os.IsNotExist(err) {
			fmt.Fprintf(out, "+ %s (cloning)\n", repo.DirectoryName)
			if err := g.Clone(repo.URL, repo.Branch, dst); err != nil {
				return fmt.Errorf("clone %s: %w", repo.DirectoryName, err)
			}
		} else if err != nil {
			return fmt.Errorf("stat %s: %w", repo.DirectoryName, err)
		} else {
			fmt.Fprintf(out, "  %s (updating)\n", repo.DirectoryName)
			if err := g.UpdateToLatest(dst, repo.Branch); err != nil {
				return fmt.Errorf("update %s: %w", repo.DirectoryName, err)
			}
		}
	}
	return nil
}
