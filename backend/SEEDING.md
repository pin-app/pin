# Development Data Seeding

This document explains the dummy data seeding system for development mode.

## Overview

When running the server in development mode (`DEV_MODE=true`), the system automatically seeds the database with realistic dummy data to facilitate frontend development and testing.

## What Gets Seeded

### Users (5 users)
- **Alice Johnson** (`alice@example.com`) - Coffee enthusiast and urban explorer
- **Bob Smith** (`bob@example.com`) - Foodie and travel blogger  
- **Charlie Brown** (`charlie@example.com`) - Photographer and nature lover
- **Diana Prince** (`diana@example.com`) - Fitness enthusiast and wellness advocate
- **Eve Wilson** (`eve@example.com`) - Art lover and museum enthusiast

Each user has:
- Unique profile information (bio, location, display name)
- Profile picture URLs from Unsplash
- Realistic creation timestamps

### Places (6 places)
- **Blue Bottle Coffee** (San Francisco) - Coffee shop with full address and hours
- **Central Park** (New York) - Public park with visitor information
- **Pike Place Market** (Seattle) - Historic market with contact details
- **Griffith Observatory** (Los Angeles) - Observatory with viewing hours
- **Art Institute of Chicago** (Chicago) - Museum with admission info
- **Golden Gate Bridge** (San Francisco) - Landmark with 24/7 access

Each place includes:
- PostGIS geometry coordinates
- Detailed properties (address, city, category, rating, phone, website, hours)
- Realistic metadata

### Posts (6 posts)
- One post per user, each associated with a different place
- Realistic descriptions with emojis
- 1-3 images per post with captions
- Images sourced from Unsplash with proper attribution

### Ratings & Comparisons
- Each user has rated every place (ratings 60-100)
- Place comparisons showing user preferences
- Realistic rating distribution

### Comments (12+ comments)
- Comments on posts from different users
- Hierarchical comment structure with replies
- Realistic conversation threads

## Usage

### Prerequisites
1. Docker and Docker Compose installed
2. Database URL configured (automatically set by test script)
3. Dev mode enabled

### Running with Seeding

#### Option 1: Using Docker (Recommended)
```bash
# Run the complete test setup with Docker
./test_seeding.sh

# In another terminal, test the API
./test_api.sh
```

#### Option 2: Manual Setup
```bash
# Start database with Docker Compose
cd /Users/race/Documents/pin
docker-compose up -d db

# Set up environment
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/pin?sslmode=disable"
export DEV_MODE=true

# Run server (seeding happens automatically)
cd backend
go run cmd/server/main.go
```

## API Testing

Once seeded, you can test the data via these endpoints:

```bash
# Get all users
curl http://localhost:8080/api/users

# Get all places  
curl http://localhost:8080/api/places

# Get all posts
curl http://localhost:8080/api/posts

# Get ratings for a specific place
curl http://localhost:8080/api/places/{place_id}/ratings

# Get comments for a specific post
curl http://localhost:8080/api/posts/{post_id}/comments
```

## Development Notes

### Seeding Logic
- Seeding only occurs in dev mode when `DEV_MODE=true`
- Requires a valid database connection
- Seeding happens once per server startup
- Data is not cleared between runs (additive)

### Data Relationships
- Users are created first
- Places are created second  
- Posts reference both users and places
- Ratings and comparisons link users to places
- Comments link users to posts with hierarchical structure

### Error Handling
- Seeding errors are logged but don't stop server startup
- Missing database connection causes server to exit in dev mode
- Individual seeding failures are logged with context

## Customization

To modify the seeded data, edit `/internal/seed/seed.go`:

- **Users**: Modify the `seedUsers()` function
- **Places**: Modify the `seedPlaces()` function  
- **Posts**: Modify the `seedPosts()` function
- **Ratings**: Modify the `seedRatings()` function
- **Comments**: Modify the `seedComments()` function

## Troubleshooting

### Common Issues

1. **"Docker is not running"**
   - Start Docker Desktop or Docker daemon
   - Ensure Docker is accessible from command line

2. **"dev mode requires DATABASE_URL"**
   - Use the test script: `./test_seeding.sh`
   - Or manually set `DATABASE_URL` environment variable

3. **"SSL is not enabled on the server"**
   - The test script automatically includes `?sslmode=disable`
   - For manual setup, add `?sslmode=disable` to your DATABASE_URL

4. **"failed to seed development data"**
   - Check that Docker container is running: `docker-compose ps`
   - Verify database connection: `docker-compose exec db pg_isready -U postgres`
   - Check server logs for specific errors

5. **No data appears in API responses**
   - Verify seeding completed successfully in logs
   - Check that you're using dev mode endpoints
   - Ensure database connection is working

### Logs to Check

Look for these log messages:
- `"seeding development data"` - Seeding started
- `"development data seeded successfully"` - Seeding completed
- `"failed to seed development data"` - Seeding failed
