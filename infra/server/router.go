package server

import (
	"go-template/app/dummy"
	"go-template/app/healthcheck"
	"go-template/pkg/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitRouters(app *fiber.App, getDummyHandler *dummy.GetDummyHandler, createDummyHandler *dummy.CreateDummyHandler, healthcheckHandler *healthcheck.HealthCheckHandler, sendRequestToRabbitMQ *dummy.SendRequestToRabbitMQHandler) {

	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
	app.Get("/healthcheck", handler.Handle[healthcheck.HealthCheckRequest, healthcheck.HealthCheckResponse](healthcheckHandler))

	app.Get("/dummy/:id", handler.Handle[dummy.GetDummyRequest, dummy.GetDummyResponse](getDummyHandler))
	app.Post("/dummy", handler.Handle[dummy.CreateDummyRequest, dummy.CreateDummyResponse](createDummyHandler))
	app.Get("/dummy/send-request-to-rabbit-mq", handler.Handle[dummy.SendRequestToRabbitMQRequest, dummy.SendRequestToRabbitMQResponse](sendRequestToRabbitMQ))

}
