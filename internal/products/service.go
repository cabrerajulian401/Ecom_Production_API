package products

import (
	"context"

	repo "github.com/cabrerajulian401/ecom/internal/adapters/postgresql/sqlc"
)

type Service interface {
	// method signatures

	ListProducts(ctx context.Context) ([]repo.Product, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}

}

func (s *svc) ListProducts(ctx context.Context) ([]repo.Product, error) {
	return s.repo.ListProducts(ctx)

}
