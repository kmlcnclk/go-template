package product

import (
	"context"
	"go-template/domain"
	"io"
	"net/http"
)

type GetProductRequest struct {
	ID string `json:"id" param:"id"`
}

type GetProductResponse struct {
	Product *domain.Product `json:"product"`
}

type GetProductHandler struct {
	repository Repository
	httpClient *http.Client
}

func NewGetProductHandler(repository Repository, httpClient *http.Client) *GetProductHandler {
	return &GetProductHandler{
		repository: repository,
		httpClient: httpClient,
	}
}

func (h *GetProductHandler) Handle(ctx context.Context, req *GetProductRequest) (*GetProductResponse, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.google.com", nil)
	if err != nil {
		return nil, err
	}

	resp, err := h.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if _, err = io.ReadAll(resp.Body); err != nil {
		return nil, err
	}

	product, err := h.repository.GetProduct(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &GetProductResponse{Product: product}, nil
}
