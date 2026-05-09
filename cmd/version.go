package cmd

import (
	"fmt"
	"io"
)

func RunVersion(version, commit, date string, out io.Writer) {
	fmt.Fprintf(out, "gwk %s (commit: %s, built: %s)\n", version, commit, date)
}
