package repository

import (
	"log/slog"
	"os"
	"slices"
)

type Repository struct {
	contentPath string // replace with https://github.com/fsnotify/fsnotify ?
	logger      *slog.Logger
}

func NewRepository(contentPath string, logger *slog.Logger) *Repository {
	return &Repository{
		contentPath: contentPath,
		logger:      logger,
	}
}

func (r *Repository) GetStreamNames() ([]string, error) {
	files, err := os.ReadDir(r.contentPath)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(files))

	for _, file := range files {
		if file.IsDir() {
			names = append(names, file.Name())
		}
	}

	slices.Sort(names)

	return names, nil
}
