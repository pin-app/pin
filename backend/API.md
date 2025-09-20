# API docs

## Base URL
```
http://localhost:8080/api
```

## Authentication

The API uses OAuth 2.0 for authentication with Google and Apple providers. Most endpoints require authentication via session tokens.

### Authentication Methods

1. **OAuth 2.0** - Use Google or Apple OAuth for user authentication
2. **Session Tokens** - Use Bearer tokens for authenticated requests
3. **Dev Mode** - In development, authentication can be bypassed using dev headers

### OAuth Flow

1. Redirect user to `/api/auth/google` or `/api/auth/apple`
2. User completes OAuth flow with provider
3. Provider redirects to callback endpoint
4. API returns session token and user information
5. Use session token in `Authorization: Bearer <token>` header for subsequent requests

### Dev Mode

Set `DEV_MODE=true` environment variable to enable development mode:
- Authentication is bypassed for requests with `X-Dev-User-ID` header
- A default dev user is created automatically
- Useful for development and testing

## API Endpoints

### Authentication

#### Google OAuth
```http
GET /api/auth/google?redirect_url=https://yourapp.com/callback
```
Redirects user to Google OAuth consent screen.

#### Google OAuth Callback
```http
GET /api/auth/google/callback?code=...&state=...
```
Handles Google OAuth callback and returns session token.

#### Apple OAuth
```http
GET /api/auth/apple?redirect_url=https://yourapp.com/callback
```
Redirects user to Apple OAuth consent screen.

#### Apple OAuth Callback
```http
GET /api/auth/apple/callback?code=...&state=...
```
Handles Apple OAuth callback and returns session token.

#### Logout
```http
POST /api/auth/logout
Authorization: Bearer <session_token>
```

Response:
```json
{
  "message": "Logged out successfully"
}
```

### Users

#### Create User
```http
POST /api/users
Content-Type: application/json

{
  "email": "user@example.com",
  "username": "username",
  "bio": "User bio",
  "location": "City, Country",
  "display_name": "Display Name",
  "pfp_url": "https://example.com/avatar.jpg"
}
```
*No authentication required*

#### Get User
```http
GET /api/users/{id}
```
*Optional authentication - returns additional data if authenticated*

#### Update User
```http
PUT /api/users/{id}
Authorization: Bearer <session_token>
Content-Type: application/json

{
  "username": "newusername",
  "bio": "Updated bio",
  "location": "New City, Country",
  "display_name": "New Display Name",
  "pfp_url": "https://example.com/new-avatar.jpg"
}
```
*Requires authentication - user can only update their own profile*

#### Delete User
```http
DELETE /api/users/{id}
Authorization: Bearer <session_token>
```
*Requires authentication - user can only delete their own profile*

#### List Users
```http
GET /api/users?limit=20&offset=0
```
*Optional authentication*

#### Search Users
```http
GET /api/users/search?q=searchterm&limit=20&offset=0
```
*Optional authentication*

### Places

#### Create Place
```http
POST /api/places
Authorization: Bearer <session_token>
Content-Type: application/json

{
  "name": "Place Name",
  "geometry": "POINT(-122.4194 37.7749)",
  "properties": {
    "address": "123 Main St",
    "city": "San Francisco",
    "state": "CA"
  }
}
```
*Requires authentication*

#### Get Place
```http
GET /api/places/{id}
```
*Optional authentication*

#### Update Place
```http
PUT /api/places/{id}
Authorization: Bearer <session_token>
Content-Type: application/json

{
  "name": "Updated Place Name",
  "geometry": "POINT(-122.4194 37.7749)",
  "properties": {
    "address": "456 New St",
    "city": "San Francisco",
    "state": "CA"
  }
}
```
*Requires authentication*

#### Delete Place
```http
DELETE /api/places/{id}
Authorization: Bearer <session_token>
```
*Requires authentication*

#### List Places
```http
GET /api/places?limit=20&offset=0
```
*Optional authentication*

#### Search Places
```http
GET /api/places/search?q=searchterm&limit=20&offset=0
```
*Optional authentication*

#### Search Nearby Places
```http
GET /api/places/nearby?lat=37.7749&lng=-122.4194&radius=10&limit=20
```
*Optional authentication*

### Posts

#### Create Post
```http
POST /api/posts
Content-Type: application/json

{
  "place_id": "uuid",
  "description": "Post description",
  "images": [
    "https://example.com/image1.jpg",
    "https://example.com/image2.jpg"
  ]
}
```

#### Get Post
```http
GET /api/posts/{id}
```

#### Update Post
```http
PUT /api/posts/{id}
Content-Type: application/json

{
  "description": "Updated post description"
}
```

#### Delete Post
```http
DELETE /api/posts/{id}
```

#### List Posts (Feed)
```http
GET /api/posts?limit=20&offset=0
```

#### List Posts by User
```http
GET /api/users/{id}/posts?limit=20&offset=0
```

#### List Posts by Place
```http
GET /api/places/{id}/posts?limit=20&offset=0
```

### Comments

#### Create Comment
```http
POST /api/comments
Content-Type: application/json

{
  "post_id": "uuid",
  "parent_id": "uuid", // Optional, for replies
  "content": "Comment content"
}
```

#### Get Comment
```http
GET /api/comments/{id}
```

#### Update Comment
```http
PUT /api/comments/{id}
Content-Type: application/json

{
  "content": "Updated comment content"
}
```

#### Delete Comment
```http
DELETE /api/comments/{id}
```

#### List Comments by Post
```http
GET /api/posts/{id}/comments?limit=20&offset=0
```

#### List Comments by User
```http
GET /api/users/{id}/comments?limit=20&offset=0
```

#### Get Comment Replies
```http
GET /api/comments/{id}/replies?limit=20&offset=0
```

### Ratings

#### Create/Update Place Rating
```http
POST /api/places/{id}/ratings
Content-Type: application/json

{
  "rating": 85
}
```

#### Get My Rating for Place
```http
GET /api/places/{id}/ratings/me
```

#### Update Place Rating
```http
PUT /api/places/{id}/ratings
Content-Type: application/json

{
  "rating": 90
}
```

#### Delete Place Rating
```http
DELETE /api/places/{id}/ratings
```

#### List Ratings by Place
```http
GET /api/places/{id}/ratings?limit=20&offset=0
```

#### Get Average Rating for Place
```http
GET /api/places/{id}/ratings/average
```

#### Create Place Comparison
```http
POST /api/places/compare
Content-Type: application/json

{
  "better_place_id": "uuid",
  "worse_place_id": "uuid"
}
```

#### List Comparisons by User
```http
GET /api/users/{id}/comparisons?limit=20&offset=0
```

## Response Format

All successful responses return JSON with the following structure:

### Success Response
```json
{
  "id": "uuid",
  "field1": "value1",
  "field2": "value2",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

### Error Response
```json
{
  "error": "Error message"
}
```

### List Response
```json
{
  "items": [...],
  "limit": 20,
  "offset": 0,
  "count": 15
}
```

## Status Codes

- `200 OK` - Request successful
- `201 Created` - Resource created successfully
- `204 No Content` - Resource deleted successfully
- `400 Bad Request` - Invalid request data
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## Database Schema

The API uses PostgreSQL with the following main tables:

- `users` - User profiles and authentication
- `oauth_accounts` - OAuth provider accounts
- `sessions` - User sessions
- `places` - Geographic locations with PostGIS support
- `place_relations` - Relationships between places
- `posts` - User posts about places
- `post_images` - Images associated with posts
- `comments` - Hierarchical comments on posts
- `place_ratings` - User ratings for places (0-100)
- `place_comparisons` - Relative comparisons between places

## Environment Variables

Required:
- `DATABASE_URL` - PostgreSQL connection string

Optional:
- `PORT` - Server port (default: 8080)
- `DEV_MODE` - Enable development mode (default: false)

OAuth Configuration:
- `GOOGLE_CLIENT_ID` - Google OAuth client ID
- `GOOGLE_CLIENT_SECRET` - Google OAuth client secret
- `GOOGLE_REDIRECT_URL` - Google OAuth redirect URL
- `APPLE_CLIENT_ID` - Apple OAuth client ID
- `APPLE_CLIENT_SECRET` - Apple OAuth client secret
- `APPLE_REDIRECT_URL` - Apple OAuth redirect URL

## Running the Server

1. Set up PostgreSQL database
2. Set `DATABASE_URL` environment variable
3. Configure OAuth providers (optional for dev mode)
4. Run migrations automatically on startup
5. Start server: `go run cmd/server/main.go`

The server will start on port 8080 by default (configurable via `PORT` environment variable).

### Development Mode

To run in development mode with authentication bypass:
```bash
DEV_MODE=true go run cmd/server/main.go
```

In dev mode, you can use the `X-Dev-User-ID` header to specify a user ID for requests.