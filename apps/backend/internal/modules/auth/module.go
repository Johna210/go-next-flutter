package auth

import (
	"go.uber.org/fx"

	"github.com/johna210/go-next-flutter/internal/core"
	"github.com/johna210/go-next-flutter/internal/modules/auth/domain/entity"
)

var Module = fx.Module("auth",
	// Provide entity provider
	fx.Provide(
		fx.Annotate(
			entity.NewEntityProvider,
			fx.As(new(core.EntityProvider)),
		),
	),

	// Auto-register with schema manager
	fx.Invoke(func(sm *core.SchemaManager, provider core.EntityProvider) {
		if err := sm.RegisterProvider(provider); err != nil {
			panic(err)
		}
	}),
)
