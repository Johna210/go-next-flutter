package bootstrap

import (
	"go.uber.org/fx"

	"github.com/johna210/go-next-flutter/internal/core"
	"github.com/johna210/go-next-flutter/internal/modules"
)

type Application struct {
	*fx.App
}

func NewApp() *Application {
	app := fx.New(
		// Core module includes (database, logger, cache, migrations)
		core.Module,

		// Application modules (auth etc..)
		modules.Modules,
	)

	return &Application{app}
}
