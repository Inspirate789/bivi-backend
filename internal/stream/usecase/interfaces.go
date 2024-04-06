package usecase

type Repository interface {
	GetStreamNames() ([]string, error)
}

type StreamNameEncoder interface {
	EncodeToString(src []byte) string
}
