package repositories

type Repository interface {
	Fetch() []string
	DeleteBranch(branch string, promptForDeletion bool)
}
