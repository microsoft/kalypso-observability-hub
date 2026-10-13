#!/bin/bash

# Demo script to show HTTP API functionality
# This script demonstrates how to test the HTTP API endpoints manually

set -e

echo "=== Kalypso Observability Hub HTTP API Demo ==="
echo

# Check if the manager binary exists
if [ ! -f "bin/manager" ]; then
    echo "Building the manager binary..."
    make build
fi

echo "The HTTP API can be started with the main manager binary using these flags:"
echo
echo "./bin/manager \\"
echo "  --enable-http-api \\"
echo "  --http-api-bind-address=:8082 \\"
echo "  --postgres-host=localhost \\"
echo "  --postgres-port=5432 \\"
echo "  --postgres-user=postgres \\"
echo "  --postgres-password=password \\"
echo "  --postgres-dbname=kalypso \\"
echo "  --postgres-sslmode=disable"
echo

echo "Available API endpoints:"
echo "- GET /health                              - Health check"
echo "- GET /api/v1/deployment-state             - Deployment state (requires manifests_endpoint and commit_id params)"
echo "- GET /api/v1/environments                 - List environments (placeholder)"
echo "- GET /api/v1/environments/{name}          - Get environment by name" 
echo "- GET /api/v1/deployments                  - List deployments (placeholder)"
echo "- GET /api/v1/deployments/{id}             - Get deployment by ID"
echo

echo "Example curl commands (when server is running):"
echo
echo "# Health check"
echo "curl http://localhost:8082/health"
echo
echo "# Get deployment state"
echo "curl 'http://localhost:8082/api/v1/deployment-state?manifests_endpoint=https://github.com/example/gitops&commit_id=abc123'"
echo
echo "# Get environment (will fail if not in database)"
echo "curl http://localhost:8082/api/v1/environments/production"
echo
echo "# Get deployment (will fail if not in database)"
echo "curl http://localhost:8082/api/v1/deployments/1"
echo

echo "Note: The HTTP API requires a PostgreSQL database connection to function properly."
echo "The deployment-state endpoint will work with the existing storage layer."
echo "Other endpoints require data to be present in the database."
echo

echo "To view the complete API specification, see:"
echo "- api/http/openapi.yaml - OpenAPI 3.0 specification"
echo "- api/http/README.md - Detailed documentation"