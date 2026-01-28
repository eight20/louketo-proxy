# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Louketo Proxy is an OpenID Connect (OIDC) authentication proxy written in Go. It protects backend services by validating JWT tokens via cookies or bearer tokens, integrating with Keycloak or other OIDC providers.

**Note:** This project reached End of Life in November 2020.

## Build and Development Commands

```bash
# Build the binary (outputs to bin/louketo-proxy)
make build

# Build static binary (for containers)
make static

# Run all tests
make test

# Run only unit tests (verbose)
go test -v

# Run a single test
go test -v -run TestFunctionName

# Run tests with coverage
make cover

# Generate coverage report (HTML)
make coverage

# Format code before committing
make format

# Lint the code
make lint

# Run golangci-lint
make verify

# Run benchmarks
make bench

# Check spelling in Go and Markdown files
make spelling

# Build container image (supports podman or docker)
make docker
```

## Architecture

This is a single-package (`main`) Go application with all source files at the repository root.

### Core Components

- **main.go / cli.go**: Entry point and CLI setup using `urfave/cli`. Configuration is parsed via reflection from struct tags in `doc.go`.

- **doc.go**: Defines the `Config` struct with all configuration options, `Resource` for protected URL patterns, and `userContext` for token claims. Configuration can come from YAML/JSON files, CLI flags, or environment variables (prefix `PROXY_`).

- **server.go**: Contains `oauthProxy` struct - the main service object. Initializes either reverse proxy or forwarding proxy mode. Uses `go-chi/chi` for routing and `goproxy` for the underlying proxy.

- **oauth.go**: OpenID Connect operations - token retrieval, refresh, and validation using `coreos/go-oidc`.

- **handlers.go**: HTTP handlers for OAuth endpoints (`/oauth/authorize`, `/oauth/callback`, `/oauth/login`, `/oauth/logout`, `/oauth/token`).

- **middleware.go**: Authentication and authorization middleware chain. Key middlewares:
  - `authenticationMiddleware`: Validates JWT tokens
  - `admissionMiddleware`: Checks roles/groups against protected resources
  - `identityHeadersMiddleware`: Injects user claims into upstream headers

- **forwarding.go**: Forwarding proxy mode - signs outbound requests with OAuth tokens.

- **stores.go / store_boltdb.go / store_redis.go**: Token storage backends for refresh tokens (BoltDB, Redis, or encrypted cookies).

- **session.go / cookies.go**: Session management and cookie handling for access/refresh tokens.

### Proxy Modes

1. **Reverse Proxy** (default): Sits in front of upstream services, validates tokens, enforces authorization rules
2. **Forwarding Proxy**: Signs outbound requests with OAuth tokens for clients

### Key OAuth Endpoints

All under configurable `--oauth-uri` prefix (default `/oauth`):
- `/authorize` - Initiates OAuth flow
- `/callback` - OAuth callback handler
- `/login` - Resource owner password grant
- `/logout` - Logout and token revocation
- `/token` - Returns current access token
- `/health` - Health check
- `/metrics` - Prometheus metrics

## Configuration

Configuration supports YAML files (`--config`), CLI flags, and environment variables. The `Config` struct in `doc.go` defines all options via struct tags (`yaml`, `usage`, `env`).

## Testing

Tests use `stretchr/testify` for assertions. Most tests are in `*_test.go` files corresponding to their source. `server_test.go` contains integration-style tests that spin up the full proxy.
