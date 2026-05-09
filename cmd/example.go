package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

func RunExample() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting home directory: %w", err)
	}
	path := filepath.Join(home, ".gwk.json")

	if _, err := os.Stat(path); err == nil {
		fmt.Println("~/.gwk.json already exists, not overwriting")
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("checking existing config: %w", err)
	}

	content := `[
    {
        "url": "github.com/example-org/example-repo.git",
        "branch": "main",
        "tags": ["sample"]
    }
]
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing example config: %w", err)
	}
	fmt.Printf("Created %s with a sample repository entry\n", path)
	return nil
}
