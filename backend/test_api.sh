#!/bin/bash

# Test script to verify seeded data via API endpoints

BASE_URL="http://localhost:8080"

echo "Testing seeded data via API endpoints..."
echo ""

# Test users endpoint
echo "1. Testing users endpoint:"
curl -s "$BASE_URL/api/users" | jq '.[0:3] | .[] | {id, email, username, display_name, bio}' 2>/dev/null || curl -s "$BASE_URL/api/users"
echo ""
echo ""

# Test places endpoint
echo "2. Testing places endpoint:"
curl -s "$BASE_URL/api/places" | jq '.[0:3] | .[] | {id, name, properties: {address, city, category}}' 2>/dev/null || curl -s "$BASE_URL/api/places"
echo ""
echo ""

# Test posts endpoint
echo "3. Testing posts endpoint:"
curl -s "$BASE_URL/api/posts" | jq '.[0:3] | .[] | {id, description, user_id, place_id}' 2>/dev/null || curl -s "$BASE_URL/api/posts"
echo ""
echo ""

# Test ratings for a specific place
echo "4. Testing place ratings (using first place ID):"
PLACE_ID=$(curl -s "$BASE_URL/api/places" | jq -r '.[0].id' 2>/dev/null)
if [ "$PLACE_ID" != "null" ] && [ "$PLACE_ID" != "" ]; then
    echo "Testing ratings for place: $PLACE_ID"
    curl -s "$BASE_URL/api/places/$PLACE_ID/ratings" | jq '.[0:3] | .[] | {user_id, rating}' 2>/dev/null || curl -s "$BASE_URL/api/places/$PLACE_ID/ratings"
else
    echo "No places found to test ratings"
fi
echo ""
echo ""

# Test comments for a specific post
echo "5. Testing post comments (using first post ID):"
POST_ID=$(curl -s "$BASE_URL/api/posts" | jq -r '.[0].id' 2>/dev/null)
if [ "$POST_ID" != "null" ] && [ "$POST_ID" != "" ]; then
    echo "Testing comments for post: $POST_ID"
    curl -s "$BASE_URL/api/posts/$POST_ID/comments" | jq '.[0:3] | .[] | {id, content, user_id}' 2>/dev/null || curl -s "$BASE_URL/api/posts/$POST_ID/comments"
else
    echo "No posts found to test comments"
fi
echo ""
echo ""

echo "API testing complete!"
