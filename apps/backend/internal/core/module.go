package core

import (
	"context"

	"go.uber.org/fx"
)

var Module = fx.Module("core",
	fx.Provide(
		NewConfig,
		NewLogger,
	),
	fx.Invoke(registerLifecycleHooks),
)

func registerLifecycleHooks(
	lc fx.Lifecycle,
	cfg *Config,
	log Logger,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Core Module starting",
				String("cfg", cfg.App.Name),
				String("env", cfg.App.Environment),
			)
			// Run Health checks
			//
			log.Info("Core Module started successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Shutting down core module")

			// Close Connections
			// Sync logger last
			return nil
		},
	})
}
