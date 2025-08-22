# Music App Backend API Testing Script (PowerShell)
# This script demonstrates how to test the implemented APIs

$BaseUrl = "http://localhost:8085/api/v1"
$HealthUrl = "http://localhost:8085"

Write-Host "🎵 Music App Backend API Testing Script" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green

# Test health endpoint
Write-Host ""
Write-Host "1. Testing Health Check..." -ForegroundColor Yellow
try {
    $healthResponse = Invoke-RestMethod -Uri "$HealthUrl/healthz" -Method Get
    $healthResponse | ConvertTo-Json -Depth 10
} catch {
    Write-Host "❌ Health check failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Write-Host "2. Testing Version Info..." -ForegroundColor Yellow
try {
    $versionResponse = Invoke-RestMethod -Uri "$HealthUrl/version" -Method Get
    $versionResponse | ConvertTo-Json -Depth 10
} catch {
    Write-Host "❌ Version check failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test user registration
Write-Host ""
Write-Host "3. Testing User Registration..." -ForegroundColor Yellow

$registerBody = @{
    email = "test@example.com"
    password = "securepassword123"
    display_name = "Test User"
} | ConvertTo-Json

try {
    $registerResponse = Invoke-RestMethod -Uri "$BaseUrl/auth/register" -Method Post -Body $registerBody -ContentType "application/json"
    $registerResponse | ConvertTo-Json -Depth 10
    
    $accessToken = $registerResponse.data.access_token
    $refreshToken = $registerResponse.data.refresh_token
    
    if ($accessToken) {
        Write-Host ""
        Write-Host "✅ Registration successful! Access token obtained." -ForegroundColor Green
        
        # Test authenticated endpoint
        Write-Host ""
        Write-Host "4. Testing Get Current User (Protected Endpoint)..." -ForegroundColor Yellow
        try {
            $headers = @{ Authorization = "Bearer $accessToken" }
            $userResponse = Invoke-RestMethod -Uri "$BaseUrl/users/me" -Method Get -Headers $headers
            $userResponse | ConvertTo-Json -Depth 10
        } catch {
            Write-Host "❌ Get current user failed: $($_.Exception.Message)" -ForegroundColor Red
        }
        
        # Test music search
        Write-Host ""
        Write-Host "5. Testing Music Search..." -ForegroundColor Yellow
        try {
            $searchResponse = Invoke-RestMethod -Uri "$BaseUrl/music/search?q=imagine%20dragons" -Method Get -Headers $headers
            $searchResponse | ConvertTo-Json -Depth 10
        } catch {
            Write-Host "❌ Music search failed: $($_.Exception.Message)" -ForegroundColor Red
        }
        
        # Test playlist creation
        Write-Host ""
        Write-Host "6. Testing Playlist Creation..." -ForegroundColor Yellow
        $playlistBody = @{
            title = "My Test Playlist"
            description = "A playlist for testing"
            is_public = $false
        } | ConvertTo-Json
        
        try {
            $playlistResponse = Invoke-RestMethod -Uri "$BaseUrl/playlists" -Method Post -Body $playlistBody -ContentType "application/json" -Headers $headers
            $playlistResponse | ConvertTo-Json -Depth 10
        } catch {
            Write-Host "❌ Playlist creation failed: $($_.Exception.Message)" -ForegroundColor Red
        }
        
        # Test token refresh
        Write-Host ""
        Write-Host "7. Testing Token Refresh..." -ForegroundColor Yellow
        if ($refreshToken) {
            $refreshBody = @{ refresh_token = $refreshToken } | ConvertTo-Json
            try {
                $refreshResponse = Invoke-RestMethod -Uri "$BaseUrl/auth/refresh" -Method Post -Body $refreshBody -ContentType "application/json"
                $refreshResponse | ConvertTo-Json -Depth 10
            } catch {
                Write-Host "❌ Token refresh failed: $($_.Exception.Message)" -ForegroundColor Red
            }
        } else {
            Write-Host "❌ No refresh token available" -ForegroundColor Red
        }
        
    } else {
        Write-Host "❌ Registration failed or no access token received" -ForegroundColor Red
        
        # Try login instead
        Write-Host ""
        Write-Host "4. Testing User Login..." -ForegroundColor Yellow
        $loginBody = @{
            email = "test@example.com"
            password = "securepassword123"
        } | ConvertTo-Json
        
        try {
            $loginResponse = Invoke-RestMethod -Uri "$BaseUrl/auth/login" -Method Post -Body $loginBody -ContentType "application/json"
            $loginResponse | ConvertTo-Json -Depth 10
        } catch {
            Write-Host "❌ Login failed: $($_.Exception.Message)" -ForegroundColor Red
        }
    }
} catch {
    Write-Host "❌ Registration failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Write-Host "8. Testing Public Endpoints (No Authentication)..." -ForegroundColor Yellow

Write-Host ""
Write-Host "   8a. Music Categories..." -ForegroundColor Cyan
try {
    $categoriesResponse = Invoke-RestMethod -Uri "$BaseUrl/music/categories" -Method Get
    $categoriesResponse | ConvertTo-Json -Depth 10
} catch {
    Write-Host "❌ Categories failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Write-Host "   8b. Top Charts..." -ForegroundColor Cyan
try {
    $chartsResponse = Invoke-RestMethod -Uri "$BaseUrl/music/top-charts?country=US" -Method Get
    $chartsResponse | ConvertTo-Json -Depth 10
} catch {
    Write-Host "❌ Top charts failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Write-Host "✅ API Testing Complete!" -ForegroundColor Green
Write-Host ""
Write-Host "📝 Notes:" -ForegroundColor Blue
Write-Host "   - Authentication endpoints are fully functional" -ForegroundColor White
Write-Host "   - Music, playlist, and library endpoints return placeholder responses" -ForegroundColor White
Write-Host "   - Health check and version endpoints work" -ForegroundColor White
Write-Host "   - All endpoints follow the defined API structure" -ForegroundColor White
Write-Host ""
Write-Host "🔧 Next Steps:" -ForegroundColor Blue
Write-Host "   1. Implement remaining endpoint handlers" -ForegroundColor White
Write-Host "   2. Add Swagger documentation" -ForegroundColor White
Write-Host "   3. Set up proper database with Docker" -ForegroundColor White
Write-Host "   4. Add comprehensive testing" -ForegroundColor White
