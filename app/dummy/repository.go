package dummy

import (
	"context"
	"go-template/domain"
)

type Repository interface {
	CreateDummy(ctx context.Context, dummy *domain.Dummy) error
	GetDummy(ctx context.Context, id string) (*domain.Dummy, error)
}
