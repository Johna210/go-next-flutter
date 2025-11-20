package main

import (
	"context"
	"flag"

	"go.uber.org/fx"

	"github.com/johna210/go-next-flutter/internal/core"
	"github.com/johna210/go-next-flutter/internal/modules"
)

func main() {
	var migrator *core.Migrator
	var logger core.Logger

	app := fx.New(
		core.Module,
		modules.Modules,
		fx.Populate(&migrator, &logger),
		fx.NopLogger,
	)

	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		logger.Fatal(err.Error())
	}
	defer func() {
		if err := app.Stop(ctx); err != nil {
			logger.Error(err.Error())
		}
	}()

	action := flag.String("action", "", "Action: list, generate, apply, status")
	name := flag.String("name", "", "Migration name (for generate)")
	modules := flag.String("modules", "", "Comma-separated module names (empty = all)")
	flag.Parse()

	switch *action {
	case "list":
		migrator.ListModules()
	case "generate":
		if *name == "" {
			logger.Fatal("Migration name is required for generate action ")
		}
		err := migrator.GenerateMigration(*name, *modules)
		if err != nil {
			logger.Fatal(err.Error())
			panic(err)
		}
	case "apply":
		err := migrator.ApplyMigrations()
		if err != nil {
			logger.Fatal(err.Error())
			panic(err)
		}
	case "status":
		migrator.CheckStatus()
	default:
		logger.Fatal("Invalid action. Use: list, generate, apply, or status")
	}
}

// func (m *Migrator) listModules() {
// 	m.Logger.Info("Listing registered modules")

// 	fmt.Println("\n📦 Registered Modules:")
// 	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

// 	totalEntities := 0
// 	for module, count := range m.SchemaManager.GetModuleInfo() {
// 		fmt.Printf("  %-15s %d entities\n", module, count)
// 		m.Logger.Debug("Module registered",
// 			core.String("module", module),
// 			core.Int("entities", count))
// 		totalEntities += count
// 	}

// 	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
// 	fmt.Printf("  Total:          %d entities\n\n", totalEntities)
// }

// func (m *Migrator) generateMigration(name, envOverride, moduleFilter string) {
// 	env := envOverride
// 	if env == "" {
// 		env = m.Config.App.Environment
// 		m.Logger.Info("Using environment from config", core.String("env", env))
// 	}

// 	var entities []interface{}
// 	var targetModules string

// 	if moduleFilter == "" {
// 		entities = m.SchemaManager.GetAllEntities()
// 		targetModules = "ALL"
// 		m.Logger.Info("Generating migration for all module",
// 			core.Int("total_entities", len(entities)))
// 	} else {
// 		mods := strings.Split(moduleFilter, ",")
// 		for i, mod := range mods {
// 			mods[i] = strings.TrimSpace(mod)
// 		}
// 		entities = m.SchemaManager.GetEntitiesByModules(mods...)
// 		targetModules = strings.Join(mods, ",")
// 		m.Logger.Info("Generating migration for specific modules",
// 			core.Int("total_entities", len(entities)))
// 	}

// 	fmt.Printf("\n🔨 Generating Migration\n")
// 	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
// 	fmt.Printf("  Name:     %s\n", name)
// 	fmt.Printf("  Modules:  %s\n", targetModules)
// 	fmt.Printf("  Entities: %d\n", len(entities))
// 	fmt.Printf("  Env:      %s\n", env)
// 	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

// 	dsn := m.Config.GetDSN()

// 	// nolint:gosec // G204: Arguments are derived from validated application configuration, not untrusted user input.
// 	cmd := exec.Command("atlas", "migrate", "diff", name,
// 		"--dir", "file://migrations",
// 		"--dev-url", dsn,
// 		"--env", env,
// 	)

// 	m.Logger.Info("Running Atlas migration generation",
// 		core.String("migration_name", name),
// 		core.String("environment", env))

// 	if err := cmd.Run(); err != nil {
// 		m.Logger.Error("Migration generation failed", core.Error(err))
// 	}

// 	m.Logger.Info("Migration generated successfully", core.String("name", name))
// 	fmt.Println("\n✅ Migration generated successfully!")
// }

// func (m *Migrator) applyMigrations(envOverride string) {
// 	env := envOverride
// 	if env == "" {
// 		env = m.Config.App.Environment
// 		m.Logger.Info("Using environment from config", core.String("env", env))
// 	}

// 	fmt.Printf("\n🔄 Applying Migrations\n")
// 	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
// 	fmt.Printf("  Environment: %s\n", env)
// 	fmt.Printf("  Database: %s", m.Config.GetDSN())
// 	fmt.Printf("  Modules:     %d registered\n", len(m.SchemaManager.ListModules()))
// 	fmt.Printf("  Entities:    %d total\n", len(m.SchemaManager.GetAllEntities()))
// 	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

// 	m.Logger.Info("Applying database migrations",
// 		core.String("environment", env),
// 		core.String("database", m.Config.Database.DBName))

// 	dsn := m.Config.GetDSN()

// 	// nolint:gosec // G204: Arguments are derived from validated application configuration, not untrusted user input.
// 	cmd := exec.Command("atlas", "migrate", "apply",
// 		"--dir", "file://migrations",
// 		"--dev-url", dsn,
// 		"--env", env,
// 	)

// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr

// 	if err := cmd.Run(); err != nil {
// 		m.Logger.Error("Migration application failed", core.Error(err))
// 	}

// 	m.Logger.Info("Migrations applied successfully")
// }

// func (m *Migrator) checkStatus(envOverride string) {
// 	env := envOverride
// 	if env == "" {
// 		env = m.Config.App.Environment
// 	}

// 	fmt.Printf("\n📊 Migration Status\n")
// 	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
// 	fmt.Printf("  Environment: %s\n", env)
// 	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

// 	dsn := m.Config.GetDSN()

// 	// nolint:gosec // G204: Arguments are derived from validated application configuration, not untrusted user input.
// 	cmd := exec.Command("atlas", "migrate", "status",
// 		"--dir", "file://migrations",
// 		"--env", env,
// 		"--url", dsn,
// 	)

// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr

// 	if err := cmd.Run(); err != nil {
// 		m.Logger.Error("Failed to check migration status", core.Error(err))
// 		m.Logger.Fatal("❌ Status check failed")
// 	}
// }
