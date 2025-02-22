package product

import (
	"context"
	"go-template/domain"

	"github.com/google/uuid"
)

type CreateProductRequest struct {
	Name string `json:"name"`
}

type CreateProductResponse struct {
	ID string `json:"id"`
}

type CreateProductHandler struct {
	repository Repository
}

func NewCreateProductHandler(repository Repository) *CreateProductHandler {
	return &CreateProductHandler{
		repository: repository,
	}
}

func (h *CreateProductHandler) Handle(ctx context.Context, req *CreateProductRequest) (*CreateProductResponse, error) {
	productID := uuid.New().String()

	product := &domain.Product{
		ID:   productID,
		Name: req.Name,
	}

	err := h.repository.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	return &CreateProductResponse{ID: product.ID}, nil
}
