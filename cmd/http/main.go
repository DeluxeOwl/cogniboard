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
	"github.com/labstack/echo/v4"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

type Options struct {
	PostgresDSN string `help:"The postgres connection string"`
	Host        string `help:"The host:port to listen on" default:"127.0.0.1:8888"`
}

func main() {

	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		e := echo.New()
		v1 := e.Group("/v1/api")

		config := huma.DefaultConfig("CogniBoard", "0.0.1")
		config.Servers = []*huma.Server{
			{URL: fmt.Sprintf("http://%s/v1/api", options.Host)},
		}

		dsn := options.PostgresDSN
		db, err := postgres.NewPostgresWithMigrate(dsn)
		if err != nil {
			panic(err)
		}

		repo, err := adapters.NewPostgresTaskRepository(db)
		if err != nil {
			panic(err)
		}

		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
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
