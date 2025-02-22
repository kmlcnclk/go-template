package product

import (
	"context"
	"go-template/domain"
)

type Repository interface {
	CreateProduct(ctx context.Context, product *domain.Product) error
	GetProduct(ctx context.Context, id string) (*domain.Product, error)
}
