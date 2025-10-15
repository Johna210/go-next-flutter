#!/bin/bash

# Export OpenAPI specification
echo "📝 Exporting OpenAPI specification..."

# Create api-spec directory if it doesn't exist
mkdir -p ../api-spec

# Run the export command
go run cmd/export-spec/main.go

echo "✅ OpenAPI spec exported successfully!"
echo "📄 File location: ../api-spec/openapi.json"
