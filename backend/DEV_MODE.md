# Development Mode

This document explains how to use the development mode for the Pin API backend.

## Enabling Dev Mode

Set the `DEV_MODE` environment variable to `true`:

```bash
export DEV_MODE=true
go run cmd/server/main.go
```

Or run directly:

```bash
DEV_MODE=true go run cmd/server/main.go
```

## How Dev Mode Works

When dev mode is enabled:

1. **Authentication Bypass**: Most endpoints that normally require authentication will work without a valid session token
2. **Dev User**: A default dev user (`dev@localhost`) is automatically created
3. **Dev User ID Header**: You can specify a specific user ID using the `X-Dev-User-ID` header
4. **Query Parameter**: Alternatively, use the `dev_user_id` query parameter

## Usage Examples

### Using Default Dev User

```bash
# This will use the default dev user
curl http://localhost:8080/api/users

# This will also use the default dev user
curl http://localhost:8080/api/places
```

### Using Custom Dev User ID

```bash
# Using header
curl -H "X-Dev-User-ID: 123e4567-e89b-12d3-a456-426614174000" \
     http://localhost:8080/api/users

# Using query parameter
curl "http://localhost:8080/api/users?dev_user_id=123e4567-e89b-12d3-a456-426614174000"
```

### Testing Protected Endpoints

```bash
# Create a place (normally requires auth)
curl -X POST http://localhost:8080/api/places \
     -H "Content-Type: application/json" \
     -H "X-Dev-User-ID: 123e4567-e89b-12d3-a456-426614174000" \
     -d '{
       "name": "Test Place",
       "geometry": "POINT(-122.4194 37.7749)",
       "properties": {
         "address": "123 Test St",
         "city": "Test City"
       }
     }'
```

## OAuth in Dev Mode

OAuth endpoints still work in dev mode, but you can also test without them:

```bash
# OAuth still works
curl http://localhost:8080/api/auth/google

# But you can also test with dev mode
curl -H "X-Dev-User-ID: 123e4567-e89b-12d3-a456-426614174000" \
     http://localhost:8080/api/users/me
```

## Environment Variables for Dev Mode

```bash
# Required
export DATABASE_URL="postgres://user:password@localhost/pin_db?sslmode=disable"
export DEV_MODE=true

# Optional
export PORT=8080
```

**Note**: In dev mode, a database connection is required for seeding dummy data. Make sure your PostgreSQL database is running and accessible.

## Disabling Dev Mode

To run in production mode, simply don't set `DEV_MODE` or set it to `false`:

```bash
export DEV_MODE=false
go run cmd/server/main.go
```

In production mode, all protected endpoints will require valid OAuth authentication.
