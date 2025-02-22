package dummy

import (
	"context"
	"go-template/domain"

	"github.com/google/uuid"
)

type CreateDummyRequest struct {
	Name string `json:"name"`
}

type CreateDummyResponse struct {
	ID string `json:"id"`
}

type CreateDummyHandler struct {
	repository Repository
}

func NewCreateDummyHandler(repository Repository) *CreateDummyHandler {
	return &CreateDummyHandler{
		repository: repository,
	}
}

func (h *CreateDummyHandler) Handle(ctx context.Context, req *CreateDummyRequest) (*CreateDummyResponse, error) {
	dummyID := uuid.New().String()

	dummy := &domain.Dummy{
		ID:   dummyID,
		Name: req.Name,
	}

	err := h.repository.CreateDummy(ctx, dummy)
	if err != nil {
		return nil, err
	}

	return &CreateDummyResponse{ID: dummy.ID}, nil
}
