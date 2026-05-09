//go:generate go run install_tools.go

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ppowo/gwk/cmd"
	"github.com/ppowo/gwk/internal/config"
	"github.com/ppowo/gwk/internal/git"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "version":
		cmd.RunVersion(version, commit, date, os.Stdout)
	case "sync":
		if err := runSync(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "status":
		ok, err := runStatus()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if !ok {
			os.Exit(2)
		}
	case "example":
		if err := cmd.RunExample(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("gwk — Git Walk")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  gwk sync     Clone/update all configured repositories")
	fmt.Println("  gwk status   Check if local clones are up to date")
	fmt.Println("  gwk example  Create ~/.gwk.json with a sample entry")
	fmt.Println("  gwk version  Print version info")
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	return filepath.Join(home, ".gwk.json"), nil
}

func mirrorDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	return filepath.Join(home, "code-mirror"), nil
}

func runSync() error {
	cp, err := configPath()
	if err != nil {
		return err
	}
	repos, err := config.Load(cp)
	if err != nil {
		return err
	}

	md, err := mirrorDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(md, 0755); err != nil {
		return fmt.Errorf("creating code-mirror directory: %w", err)
	}

	g, err := git.NewExec(md)
	if err != nil {
		return fmt.Errorf("initializing git adapter: %w", err)
	}

	return cmd.RunSync(repos, g, md, os.Stdout)
}

func runStatus() (bool, error) {
	cp, err := configPath()
	if err != nil {
		return false, err
	}
	repos, err := config.Load(cp)
	if err != nil {
		return false, err
	}

	md, err := mirrorDir()
	if err != nil {
		return false, err
	}

	g, err := git.NewExec(md)
	if err != nil {
		return false, fmt.Errorf("initializing git adapter: %w", err)
	}

	return cmd.RunStatus(repos, g, md, os.Stdout)
}
