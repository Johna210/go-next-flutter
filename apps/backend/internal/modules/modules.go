package modules

import (
	"go.uber.org/fx"

	"github.com/johna210/go-next-flutter/internal/modules/auth"
)

var Modules = fx.Options(
	auth.Module,
)
