package server

import (
	"go-template/app/dummy"
	"go-template/app/healthcheck"
	"go-template/pkg/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitRouters(app *fiber.App, getDummyHandler *dummy.GetDummyHandler, createDummyHandler *dummy.CreateDummyHandler, healthcheckHandler *healthcheck.HealthCheckHandler) {

	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
	app.Get("/healthcheck", handler.Handle[healthcheck.HealthCheckRequest, healthcheck.HealthCheckResponse](healthcheckHandler))

	app.Get("/dummys/:id", handler.Handle[dummy.GetDummyRequest, dummy.GetDummyResponse](getDummyHandler))
	app.Post("/dummys", handler.Handle[dummy.CreateDummyRequest, dummy.CreateDummyResponse](createDummyHandler))

}
