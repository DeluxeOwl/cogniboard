package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/DeluxeOwl/cogniboard/internal/postgres"
	"github.com/DeluxeOwl/cogniboard/internal/project/adapters"
	"github.com/DeluxeOwl/cogniboard/internal/project/app"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/golang-cz/devslog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sanity-io/litter"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

const (
	EnvDev     = "dev"
	APIVersion = "v1/api"
)

type Options struct {
	PostgresDSN string `help:"The postgres connection string"`
	Host        string `help:"The host:port to listen on" default:"127.0.0.1:8888"`
	Env         string `help:"The environment to run in" default:"dev"`
}

func main() {
	litter.Config.Compact = true
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		e := setupEcho()
		api := setupAPI(e, options.Host)
		logger := setupLogger(options.Env)
		app := setupApplication(options.PostgresDSN, logger)
		setupHTTPHandlers(api, app)

		hooks.OnStart(func() {
			logger.Info("Server started", "host", options.Host)
			saveOpenAPISpec(api)
			http.ListenAndServe(options.Host, e)
		})
	})

	cli.Run()
}

func setupEcho() *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"}, //TODO: hardcoded
	}))
	v1 := e.Group("/" + APIVersion)

	// Setup documentation endpoint
	v1.GET("/docs", func(c echo.Context) error {
		return c.HTML(http.StatusOK, fmt.Sprintf(`<!doctype html>
<html>
  <head>
    <title>API Reference</title>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script
      id="api-reference"
      data-url="/%s/openapi.json"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`, APIVersion))
	})

	return e
}

func setupAPI(e *echo.Echo, host string) huma.API {
	config := huma.DefaultConfig("CogniBoard", "0.0.1")
	config.Servers = []*huma.Server{
		{URL: fmt.Sprintf("http://%s/%s", host, APIVersion)},
	}
	config.DocsPath = ""

	return humaecho.NewWithGroup(e, e.Group("/"+APIVersion), config)
}

func setupLogger(env string) *slog.Logger {
	var handler slog.Handler
	if env == EnvDev {
		handler = devslog.NewHandler(os.Stdout, nil)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, nil)
	}
	return slog.New(handler)
}

func setupApplication(dsn string, logger *slog.Logger) *app.Application {
	db, err := postgres.NewPostgresWithMigrate(dsn)
	if err != nil {
		panic(err)
	}

	repo, err := adapters.NewPostgresTaskRepository(db)
	if err != nil {
		panic(err)
	}

	app, err := app.New(repo, logger)
	if err != nil {
		panic(err)
	}

	return app
}

func setupHTTPHandlers(api huma.API, app *app.Application) {
	projectHTTP := adapters.NewHuma(api, app)
	projectHTTP.Register()
}

func saveOpenAPISpec(api huma.API) {
	openapiSpec, err := api.OpenAPI().YAML()
	if err != nil {
		panic(err)
	}
	filePath := filepath.Join("openapi3.yaml")
	err = os.WriteFile(filePath, openapiSpec, 0644)
	if err != nil {
		panic(err)
	}
}
