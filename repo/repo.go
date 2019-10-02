package repo

// Repository is an abstraction over a (git) repository
type Repository interface {

	// GetCurrentBranch returns the name of the current branch
	GetCurrentBranch() (string, error)

	// CheckoutSpinOffBranch creates a new branch at current HEAD
	CheckoutSpinOffBranch(branch string) error

	// CommitAllChanged stages all changes in repo and commits them
	CommitAllChanged(msg string) error
}
