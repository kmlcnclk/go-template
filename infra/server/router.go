package server

import (
	"go-template/app/healthcheck"
	"go-template/app/product"
	"go-template/pkg/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitRouters(app *fiber.App, getProductHandler *product.GetProductHandler, createProductHandler *product.CreateProductHandler, healthcheckHandler *healthcheck.HealthCheckHandler) {

	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
	app.Get("/healthcheck", handler.Handle[healthcheck.HealthCheckRequest, healthcheck.HealthCheckResponse](healthcheckHandler))

	app.Get("/products/:id", handler.Handle[product.GetProductRequest, product.GetProductResponse](getProductHandler))
	app.Post("/products", handler.Handle[product.CreateProductRequest, product.CreateProductResponse](createProductHandler))

}
