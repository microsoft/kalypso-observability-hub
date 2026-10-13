# Kalypso Observability Hub HTTP API

This HTTP API provides external access to deployment observability data stored in the Kalypso Observability Hub. It enables external systems like GitHub CD workflows to query deployment states and related information.

## Features

- **Deployment State Queries**: Query the deployment state for specific manifests endpoints and commit IDs
- **Environment Management**: Retrieve environment information
- **Deployment Information**: Access deployment details and status
- **Health Checks**: Monitor API server health

## API Endpoints

### Health Check
- `GET /health` - Returns API server health status

### Deployment State
- `GET /api/v1/deployment-state?manifests_endpoint={url}&commit_id={id}` - Get deployment state for a specific manifests endpoint and commit

### Environments
- `GET /api/v1/environments` - List all environments (placeholder)
- `GET /api/v1/environments/{name}` - Get specific environment by name

### Deployments
- `GET /api/v1/deployments` - List all deployments (placeholder)
- `GET /api/v1/deployments/{id}` - Get specific deployment by ID

## Configuration

The HTTP API server can be enabled in the main Kalypso Observability Hub application with the following flags:

```bash
./manager \
  --enable-http-api \
  --http-api-bind-address=:8082 \
  --postgres-host=localhost \
  --postgres-port=5432 \
  --postgres-user=postgres \
  --postgres-password=password \
  --postgres-dbname=kalypso \
  --postgres-sslmode=disable
```

### Configuration Options

- `--enable-http-api`: Enable the HTTP API server (default: false)
- `--http-api-bind-address`: Address for the HTTP API server to bind to (default: :8082)
- `--postgres-host`: PostgreSQL host for database connections (default: localhost)
- `--postgres-port`: PostgreSQL port (default: 5432)
- `--postgres-user`: PostgreSQL username (default: postgres)
- `--postgres-password`: PostgreSQL password (default: empty)
- `--postgres-dbname`: PostgreSQL database name (default: postgres)
- `--postgres-sslmode`: PostgreSQL SSL mode (default: disable)

## OpenAPI Specification

The API is documented using OpenAPI 3.0 specification. See [openapi.yaml](./openapi.yaml) for the complete specification.

## Example Usage

### Get Deployment State

```bash
curl "http://localhost:8082/api/v1/deployment-state?manifests_endpoint=https://github.com/example/gitops&commit_id=abc123"
```

Response:
```json
{
  "total_subscribers": 5,
  "total_succeeded_subscribers": 4,
  "total_failed_subscribers": 1,
  "total_in_progress_subscribers": 0,
  "succeeded_subscribers": [],
  "failed_subscribers": [
    {
      "name": "cluster-1",
      "status_message": "Deployment failed: timeout"
    }
  ],
  "in_progress_subscribers": []
}
```

### Get Environment

```bash
curl "http://localhost:8082/api/v1/environments/production"
```

Response:
```json
{
  "id": 1,
  "name": "production",
  "description": "Production environment for live workloads"
}
```

### Health Check

```bash
curl "http://localhost:8082/health"
```

Response:
```json
{
  "status": "healthy"
}
```

## Development

### Running Tests

```bash
go test ./api/http/
```

### Building

The HTTP API is built as part of the main application:

```bash
make build
```

## Integration

The HTTP API integrates with the existing Kalypso Observability Hub storage layer and uses the same PostgreSQL database as the gRPC storage API. It provides a RESTful interface for external systems to query deployment observability data without needing to use gRPC or direct database access.