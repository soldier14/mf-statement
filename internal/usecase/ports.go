package usecase

import (
	"context"
	"io"
	"mf-statement/internal/domain"
)

type Source interface {
	Open(ctx context.Context, uri string) (io.ReadCloser, error)
}

type Parser interface {
	Parse(ctx context.Context, reader io.Reader) ([]domain.Transaction, error)
}
