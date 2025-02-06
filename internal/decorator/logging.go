package decorator

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/sanity-io/litter"
)

type withCmdLogger[C any] struct {
	base CommandHandler[C]
	log  *slog.Logger
}

func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}

func (d withCmdLogger[C]) Handle(ctx context.Context, cmd C) (err error) {
	handlerType := generateActionName(cmd)

	logger := d.log.With(
		"command", handlerType,
		"command_body", litter.Sdump(cmd),
	)

	logger.Debug("executing command")
	defer func() {
		if err == nil {
			logger.Info("command executed successfully")
		} else {
			logger.Error("execute command", "error", err)
		}
	}()

	return d.base.Handle(ctx, cmd)
}

type withQueryLogger[C any, R any] struct {
	base QueryHandler[C, R]
	log  *slog.Logger
}

func (d withQueryLogger[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	logger := d.log.With(
		"query", generateActionName(cmd),
		"query_body", fmt.Sprintf("%#v", cmd),
	)

	logger.Debug("executing query")
	defer func() {
		if err == nil {
			logger.Info("query executed successfully")
		} else {
			logger.Error("execute query", "error", err)
		}
	}()

	return d.base.Handle(ctx, cmd)
}
