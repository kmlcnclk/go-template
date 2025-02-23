package dummy

import (
	"context"
	"go-template/domain"
	"go-template/infra/rabbitmq"
	"net/http"
)

type SendRequestToRabbitMQRequest struct {
	ID string `json:"id" param:"id"`
}

type SendRequestToRabbitMQResponse struct {
	Dummy *domain.Dummy `json:"dummy"`
}

type SendRequestToRabbitMQHandler struct {
	repository Repository
	httpClient *http.Client
	rmq        *rabbitmq.RabbitMQ
}

func NewSendRequestToRabbitMQHandler(repository Repository, httpClient *http.Client, rmq *rabbitmq.RabbitMQ) *SendRequestToRabbitMQHandler {
	return &SendRequestToRabbitMQHandler{
		repository: repository,
		httpClient: httpClient,
		rmq:        rmq,
	}
}

func (h *SendRequestToRabbitMQHandler) Handle(ctx context.Context, req *SendRequestToRabbitMQRequest) (*SendRequestToRabbitMQResponse, error) {
	err := h.rmq.Publish("my_exchange", "", []byte("test dummy"))

	if err != nil {
		return nil, err
	}
	return &SendRequestToRabbitMQResponse{}, nil
}
