package workerfx

import (
	"context"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/worker"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"worker",
		fx.Provide(worker.NewRegistry),
		fx.Invoke(func(lc fx.Lifecycle, registry worker.Registry) {
			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					stopCtx := ctx
					stopCancel := func() {}
					if ctx.Err() != nil {
						stopCtx, stopCancel = context.WithTimeout(context.Background(), 10*time.Second)
					}
					defer stopCancel()
					return registry.Stop(stopCtx)
				},
			})
		}),
	)
}
