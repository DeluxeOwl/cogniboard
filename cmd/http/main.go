package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/DeluxeOwl/cogniboard/internal/postgres"
	"github.com/DeluxeOwl/cogniboard/internal/project/adapters"
	"github.com/DeluxeOwl/cogniboard/internal/project/app"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/golang-cz/devslog"
	"github.com/labstack/echo/v4"
	"github.com/sanity-io/litter"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

const EnvDev = "dev"

type Options struct {
	PostgresDSN string `help:"The postgres connection string"`
	Host        string `help:"The host:port to listen on" default:"127.0.0.1:8888"`
	Env         string `help:"The environment to run in" default:"dev"`
}

func main() {
	litter.Config.Compact = true
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		e := echo.New()

		v1apiGroup := "v1/api"
		v1 := e.Group("/" + v1apiGroup)

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
</html>`, v1apiGroup))
		})

		config := huma.DefaultConfig("CogniBoard", "0.0.1")
		config.Servers = []*huma.Server{
			{URL: fmt.Sprintf("http://%s/%s", options.Host, v1apiGroup)},
		}
		config.DocsPath = ""

		dsn := options.PostgresDSN
		db, err := postgres.NewPostgresWithMigrate(dsn)
		if err != nil {
			panic(err)
		}

		repo, err := adapters.NewPostgresTaskRepository(db)
		if err != nil {
			panic(err)
		}

		var handler slog.Handler
		if options.Env == EnvDev {
			handler = devslog.NewHandler(os.Stdout, nil)
		} else {
			handler = slog.NewJSONHandler(os.Stdout, nil)
		}

		logger := slog.New(handler)
		app, err := app.New(repo, logger)
		if err != nil {
			panic(err)
		}

		api := humaecho.NewWithGroup(e, v1, config)
		projectHTTP := adapters.NewHuma(api, app)

		projectHTTP.Register()

		hooks.OnStart(func() {
			logger.Info("Server started", "host", options.Host)
			http.ListenAndServe(options.Host, e)
		})
	})

	cli.Run()
}
