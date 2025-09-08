# Pin Backend

A minimalist Go backend service for the Pin mobile app.

## Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── handlers/
│   │   ├── health.go        # Health check handler
│   │   └── routes.go        # Route registration
│   └── server/
│       └── server.go        # Server configuration and middleware
├── go.mod                   # Go module definition
├── Makefile                 # Build and run commands
└── README.md               # This file
```

## Quick Start

### Prerequisites

- Go 1.21 or later
- curl (for testing)

### Running the Server

1. **Using the run script:**
   ```bash
   ./scripts/run-server.sh
   ```

2. **Using Make:**
   ```bash
   cd backend
   make run
   ```

3. **Direct Go command:**
   ```bash
   cd backend
   go run cmd/server/main.go
   ```

The server will start on `http://localhost:8080` by default.

### Testing the API

1. **Health check:**
   ```bash
   ./scripts/health-check.sh
   ```

2. **Full API test:**
   ```bash
   ./scripts/test-api.sh
   ```

3. **Manual testing:**
   ```bash
   curl http://localhost:8080/health
   ```

## API Endpoints

- `GET /health` - Health check endpoint
- `GET /` - Redirects to health endpoint

## Environment Variables

- `PORT` - Server port (default: 8080)

## Development

### Building

```bash
make build
```

### Running Tests

```bash
make test
```

### Cleaning

```bash
make clean
```

## Next Steps

This is a minimal setup. Future enhancements will include:
- Database integration
- Authentication
- Business logic endpoints
- Configuration management
- Logging improvements
