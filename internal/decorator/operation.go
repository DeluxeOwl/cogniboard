package decorator

import (
	"context"
	"log/slog"
)

func ApplyOperationDecorators[H any, R any](
	handler OperationHandler[H, R],
	log *slog.Logger,
) OperationHandler[H, R] {
	return withOperationLogger[H, R]{
		base: handler,
		log:  log,
	}
}

type OperationHandler[Q any, R any] interface {
	Handle(ctx context.Context, q Q) (R, error)
}
