//go:generate go run install_tools.go

package main

import (
	"fmt"
	"os"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("gwk %s (commit: %s, built: %s)\n", version, commit, date)
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "inspect" {
		fmt.Println("gwk — walking repos...")
		fmt.Println("(coming soon)")
		return
	}

	fmt.Println("gwk — Git Walk")
	fmt.Println("Walk multiple git repos and inspect their latest commits.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  gwk inspect    Walk all configured repos")
	fmt.Println("  gwk version    Print version info")
}
