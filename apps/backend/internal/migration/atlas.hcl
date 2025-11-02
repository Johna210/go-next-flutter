env "local" {
  url = "postgres://johna:postgres@localhost:5432/auth_db?sslmode=false"
  src = "./migrations"
  dev = "docker://postgres/18"
}
