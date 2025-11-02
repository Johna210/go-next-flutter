package core

import (
	"fmt"
	"os/exec"
)

func ApplyMigrations(cfg *Config, sm *SchemaManager) error {
	// nolint:gosec // G204: Arguments are derived from validated application configuration, not untrusted user input.
	cmd := exec.Command("atlas", "migrate", "apply",
		"--dir", "file://migrations",
		"--env", cfg.App.Environment,
		"--url", cfg.GetDSN(),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("atlas failed: %w, output: %s", err, string(output))
	}
	return nil
}
