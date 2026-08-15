# cbr-app

A Go REST API service that serves currency exchange rates from the Central Bank of Russia. The service fetches the official XML feed (`http://www.cbr.ru/scripts/XML_daily.asp`), parses it, and returns the data as JSON.

## Features

- Get exchange rates for all currencies on a given date (or today's rate if no date is given)
- Get the rate for a specific currency by its character code (e.g. `USD`, `EUR`)
- Info endpoint returning the service version and author
- Containerized via Docker / Docker Compose with resource limits and a non-root user
- CI pipeline (Jenkins) with build, tests, SBOM generation, and vulnerability scanning (Trivy)

## Endpoints

### `GET /info`

Returns service information.

```json
{
  "version": "1.0.0",
  "service": "currency",
  "author": "a.dashchinsky"
}
```

### `GET /info/currency`

Returns CBR currency exchange rates.

**Query parameters:**

| Parameter  | Required | Description                                                                |
|------------|----------|------------------------------------------------------------------------------|
| `currency` | no       | Currency character code (`USD`, `EUR`, ...). If omitted, all are returned.  |
| `date`     | no       | Date in `YYYY-MM-DD` format. If omitted, the current CBR rate is used.      |

Examples:

```bash
# Rates for all currencies
curl "http://localhost:8000/info/currency"

# Rate for a specific currency
curl "http://localhost:8000/info/currency?currency=USD"

# Rate for a currency on a specific date
curl "http://localhost:8000/info/currency?currency=EUR&date=2024-01-15"
```

Response:

```json
{
  "service": "currency",
  "data": {
    "USD": 91.23,
    "EUR": 98.76
  }
}
```

If the requested currency code isn't found, the service returns `404 Not Found`.

## Stack

- Go 1.25
- `golang.org/x/net/html/charset` — for correctly decoding the CBR XML response's charset
- standard library `net/http` — no third-party router/framework

## Project structure

```
main.go              entry point, HTTP handler registration
api/api.go           HTTP handlers and request/response models
cbr/cbr.go            client for the CBR XML feed and rate parsing
tests/                shell tests for the endpoints
ci/                   Jenkins pipelines
Dockerfile
docker-compose.yml
```

## Running locally

Requires Go 1.25+.

```bash
git clone https://github.com/SANEKDASH/cbr-app.git
cd cbr-app
go mod download
go run main.go
```

By default the service listens on port `:8000`.

## Configuration (environment variables)

| Variable  | Default           | Description                          |
|-----------|--------------------|--------------------------------------|
| `PORT`    | `8000`            | Port the service listens on          |
| `AUTHOR`  | `a.dashchinsky`   | Author, returned from `/info`        |
| `VERSION` | `1.0.0`           | Version, returned from `/info`       |

See `.env.example` for a sample env file.

## Running with Docker

```bash
docker compose up --build
```

The container runs as a non-root user, with a read-only filesystem and no extra Linux capabilities. Default limits: 6 MB memory, 0.2 CPU (see `docker-compose.yml`).

Build the image manually:

```bash
docker build -t cbr-app .
docker run -p 8000:8000 -e PORT=8000 cbr-app
```

## Tests

Endpoint tests live in `tests/currency` and `tests/info` and run as shell scripts against a running instance of the service.

## CI/CD

The pipeline is set up with Jenkins (see `ci/` — build and parameterized deploy pipelines). The repository also keeps `sbom.json` (build SBOM) and `trivy-report.json` (image vulnerability scan) as build artifacts.
