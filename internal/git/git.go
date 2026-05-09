package git

// Git is the interface for git operations used by commands.
type Git interface {
	Clone(url, branch, dst string) error
	UpdateToLatest(repoPath, branch string) error
	RemoteHead(url, branch string) (string, error)
	LocalHead(repoPath string) (string, error)
}
