package decorator

import (
	"context"
	"log/slog"
)

func ApplyQueryDecorators[H any, R any](
	handler QueryHandler[H, R],
	log *slog.Logger,
) QueryHandler[H, R] {
	return withQueryLogger[H, R]{
		base: handler,
		log:  log,
	}
}

type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, q Q) (R, error)
}
