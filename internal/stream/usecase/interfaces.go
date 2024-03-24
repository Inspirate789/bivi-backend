package usecase

type Repository interface {
	GetStreamNames() ([]string, error)
}
