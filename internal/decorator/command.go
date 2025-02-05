package decorator

import (
	"context"
	"log/slog"
)

func ApplyCommandDecorators[H any](handler CommandHandler[H], log *slog.Logger) CommandHandler[H] {
	return withCmdLogger[H]{
		base: handler,
		log:  log,
	}
}

type CommandHandler[C any] interface {
	Handle(ctx context.Context, cmd C) error
}
