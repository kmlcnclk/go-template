package main

import (
	httpclient "go-template/pkg/http_client"
	"go-template/pkg/tracer"
	"time"

	"github.com/gofiber/fiber/v2"

	"go.uber.org/zap"

	"go-template/app/dummy"
	"go-template/app/healthcheck"
	"go-template/infra/couchbase"
	"go-template/infra/server"
	"go-template/pkg/config"
	_ "go-template/pkg/log"
)

func main() {
	appConfig := config.Read()
	defer zap.L().Sync()

	zap.L().Info("app starting...")

	tp := tracer.InitTracer()
	httpClient := httpclient.HttpClient()

	couchbaseRepository := couchbase.NewCouchbaseRepository(tp)

	// Dependency Injection
	getDummyHandler := dummy.NewGetDummyHandler(couchbaseRepository, httpClient)
	createDummyHandler := dummy.NewCreateDummyHandler(couchbaseRepository)
	healthcheckHandler := healthcheck.NewHealthCheckHandler()

	// Init Fiber app
	app := fiber.New(fiber.Config{
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		Concurrency:  256 * 1024,
	})

	server.InitMiddlewares(app)

	server.InitRouters(app, getDummyHandler, createDummyHandler, healthcheckHandler)

	server.Start(app, appConfig)

	server.GracefulShutdown(app)
}
