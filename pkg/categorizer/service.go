package categorizer

import (
	"context"
)

type Categorizer interface {
	CategorizeOne(ctx context.Context, input string) (*TransactionRecord, error)
	CategorizeMany(ctx context.Context, inputs []string) ([]TransactionRecord, error)
}

type Service struct {
	provider Categorizer
}

func NewService(provider Categorizer) *Service {
	return &Service{provider}
}

func (s *Service) CategorizeOne(ctx context.Context, input string) (*TransactionRecord, error) {
	return s.provider.CategorizeOne(ctx, input)
}

func (s *Service) CategorizeMany(ctx context.Context, inputs []string) ([]TransactionRecord, error) {
	return s.provider.CategorizeMany(ctx, inputs)
}
