package dummy

import (
	"context"
	"go-template/domain"
	"io"
	"net/http"
)

type GetDummyRequest struct {
	ID string `json:"id" param:"id"`
}

type GetDummyResponse struct {
	Dummy *domain.Dummy `json:"dummy"`
}

type GetDummyHandler struct {
	repository Repository
	httpClient *http.Client
}

func NewGetDummyHandler(repository Repository, httpClient *http.Client) *GetDummyHandler {
	return &GetDummyHandler{
		repository: repository,
		httpClient: httpClient,
	}
}

func (h *GetDummyHandler) Handle(ctx context.Context, req *GetDummyRequest) (*GetDummyResponse, error) {
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

	dummy, err := h.repository.GetDummy(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &GetDummyResponse{Dummy: dummy}, nil
}
