package repositories

type Repository interface {
	Fetch() []string
}
