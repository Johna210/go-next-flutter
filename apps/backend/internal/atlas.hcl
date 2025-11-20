# Variable for environment selection
variable "current_env" {
  type = string
  default = "local"
}

# Export GORM schema from your Go code
data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./cmd/schema",
    "--action=export-schema",
    "--env=${var.current_env}"
  ]
}

# Local environment configuration
env "local" {
  # Source: Your GORM models (desired state)
  src = data.external_schema.gorm.url
  # Target: Your actual database (read from env or hardcode)
  url = getenv("DATABASE_URL")
  # Dev: Same as target for local development
  dev = getenv("DATABASE_URL")
  # Migration settings
  migration {
    dir = "file://migrations"
  }
}

# Development environment configuration
env "development" {
  src = data.external_schema.gorm.url
  url = getenv("DATABASE_URL")
  dev = getenv("DATABASE_URL")
  migration {
    dir = "file://migrations"
  }
}

# Production environment configuration
env "production" {
  src = data.external_schema.gorm.url
  url = getenv("DATABASE_URL")
  dev = getenv("DATABASE_URL")
  migration {
    dir = "file://migrations"
    revisions_schema = "atlas_schema_revisions"
  }
}
