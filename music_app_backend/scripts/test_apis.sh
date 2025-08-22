#!/bin/bash

# Music App Backend API Testing Script
# This script demonstrates how to test the implemented APIs

BASE_URL="http://localhost:8085/api/v1"
HEALTH_URL="http://localhost:8085"

echo "🎵 Music App Backend API Testing Script"
echo "========================================"

# Test health endpoint
echo ""
echo "1. Testing Health Check..."
curl -s "$HEALTH_URL/healthz" | jq '.'

echo ""
echo "2. Testing Version Info..."
curl -s "$HEALTH_URL/version" | jq '.'

# Test user registration
echo ""
echo "3. Testing User Registration..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "securepassword123",
    "display_name": "Test User"
  }')

echo "$REGISTER_RESPONSE" | jq '.'

# Extract access token for further tests
ACCESS_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.data.access_token // empty')
REFRESH_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.data.refresh_token // empty')

if [ -n "$ACCESS_TOKEN" ] && [ "$ACCESS_TOKEN" != "null" ]; then
  echo ""
  echo "✅ Registration successful! Access token obtained."
  
  # Test authenticated endpoint
  echo ""
  echo "4. Testing Get Current User (Protected Endpoint)..."
  curl -s -H "Authorization: Bearer $ACCESS_TOKEN" "$BASE_URL/users/me" | jq '.'
  
  # Test music search (currently returns placeholder)
  echo ""
  echo "5. Testing Music Search..."
  curl -s -H "Authorization: Bearer $ACCESS_TOKEN" "$BASE_URL/music/search?q=imagine%20dragons" | jq '.'
  
  # Test playlist creation (currently returns placeholder)
  echo ""
  echo "6. Testing Playlist Creation..."
  curl -s -X POST "$BASE_URL/playlists" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
      "title": "My Test Playlist",
      "description": "A playlist for testing",
      "is_public": false
    }' | jq '.'
    
  # Test token refresh
  echo ""
  echo "7. Testing Token Refresh..."
  if [ -n "$REFRESH_TOKEN" ] && [ "$REFRESH_TOKEN" != "null" ]; then
    curl -s -X POST "$BASE_URL/auth/refresh" \
      -H "Content-Type: application/json" \
      -d "{
        \"refresh_token\": \"$REFRESH_TOKEN\"
      }" | jq '.'
  else
    echo "❌ No refresh token available"
  fi
  
else
  echo "❌ Registration failed or no access token received"
  
  # Try login instead
  echo ""
  echo "4. Testing User Login..."
  LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
      "email": "test@example.com",
      "password": "securepassword123"
    }')
  
  echo "$LOGIN_RESPONSE" | jq '.'
fi

echo ""
echo "8. Testing Public Endpoints (No Authentication)..."
echo ""
echo "   8a. Music Categories..."
curl -s "$BASE_URL/music/categories" | jq '.'

echo ""
echo "   8b. Top Charts..."
curl -s "$BASE_URL/music/top-charts?country=US" | jq '.'

echo ""
echo "✅ API Testing Complete!"
echo ""
echo "📝 Notes:"
echo "   - Authentication endpoints are fully functional"
echo "   - Music, playlist, and library endpoints return placeholder responses"
echo "   - Health check and version endpoints work"
echo "   - All endpoints follow the defined API structure"
echo ""
echo "🔧 Next Steps:"
echo "   1. Implement remaining endpoint handlers"
echo "   2. Add Swagger documentation"
echo "   3. Set up proper database with Docker"
echo "   4. Add comprehensive testing"
