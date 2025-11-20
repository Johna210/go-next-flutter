package core

import (
	"context"

	"go.uber.org/fx"
)

var Module = fx.Module("core",
	fx.Provide(
		NewConfig,
		NewLogger,
		NewDatabase,
		NewSchemaManager,
		NewMigrator,
	),
	fx.Invoke(registerLifecycleHooks),
)

func registerLifecycleHooks(
	lc fx.Lifecycle,
	cfg *Config,
	log Logger,
	db *Database,
	sm *SchemaManager,
	m *Migrator,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Core Module starting",
				String("cfg", cfg.App.Name),
				String("env", cfg.App.Environment),
			)
			// Run Health checks
			if err := db.Health(ctx); err != nil {
				return err
			}

			// RunMigrations - applies all migrations on startup
			// log.Info("Applying database migrations...")
			// log.Info("Managing",
			// 	Int("entities", len(sm.GetAllEntities())),
			// 	Int("modules", len(sm.ListModules())),
			// )

			// if err := m.applyMigrations(); err != nil {
			// 	return fmt.Errorf("migration failed: %w", err)
			// }
			// log.Info("Migrations applied successfully")

			// log.Info("Core Module started successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Shutting down core module")

			// Close Connections
			if err := db.Close(); err != nil {
				log.Error("Failed to close database", Error(err))
			}

			// Sync logger last
			return nil
		},
	})
}
