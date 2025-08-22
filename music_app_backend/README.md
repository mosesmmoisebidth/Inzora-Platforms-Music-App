# Music App Backend

A production-ready Go backend for a Flutter music application built with Gin framework, featuring robust authentication, music provider integration, and comprehensive API endpoints.

## ✨ Features

### 🔐 Authentication & Security
- **Email/Password Authentication** with Argon2id password hashing
- **JWT Access & Refresh Tokens** with automatic rotation
- **Google Sign-In Integration** with server-side ID token verification
- **Role-based Access Control** with middleware support
- **Rate Limiting** and **CORS** protection

### 🎵 Music Integration
- **Provider-agnostic Music API** with pluggable providers
- **iTunes Search API** integration (built-in)
- **Spotify Web API** support (configurable)
- **Search, Browse, and Discovery** features
- **Caching & Performance** optimization

### 📱 User Features
- **User Management** with profiles and preferences
- **Playlist Management** with full CRUD operations
- **Favorites & History** tracking
- **Download Management** (metadata only)
- **Social Features** with playlist sharing

### 🏗️ Architecture
- **Clean Architecture** with separated concerns
- **Dependency Injection** and interface-based design
- **Database Migrations** with GORM
- **Redis Caching** for performance
- **Structured Logging** with request tracing
- **Health Checks** and monitoring

## 🚀 Quick Start

### Prerequisites

- **Go 1.22+**
- **Docker & Docker Compose**
- **PostgreSQL 15+**
- **Redis 7+**

### Setup

1. **Clone and Navigate**
   ```bash
   cd music_app_backend
   ```

2. **Install Dependencies**
   ```bash
   make deps
   ```

3. **Environment Configuration**
   ```bash
   make setup
   # Edit .env file with your configuration
   ```

4. **Start Development Environment**
   ```bash
   make docker-up
   ```

5. **Access the API**
   - API: http://localhost:8080
   - Health: http://localhost:8080/healthz
   - Swagger: http://localhost:8080/docs

## 📖 API Documentation

### Authentication Endpoints

#### Register
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword",
  "display_name": "John Doe"
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword"
}
```

#### Google Sign-In
```http
POST /api/v1/auth/google
Content-Type: application/json

{
  "id_token": "google_id_token_here"
}
```

#### Refresh Token
```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "refresh_token_here"
}
```

### Music Discovery Endpoints

#### Search Tracks
```http
GET /api/v1/music/search?q=imagine%20dragons&page=1&size=20&provider=itunes
Authorization: Bearer your_access_token
```

#### Get Track Details
```http
GET /api/v1/music/tracks/123456?provider=itunes
Authorization: Bearer your_access_token
```

#### Get Top Charts
```http
GET /api/v1/music/top-charts?country=US&page=1&size=20&provider=itunes
Authorization: Bearer your_access_token
```

### Playlist Management

#### Create Playlist
```http
POST /api/v1/playlists
Authorization: Bearer your_access_token
Content-Type: application/json

{
  "title": "My Awesome Playlist",
  "description": "My favorite songs",
  "is_public": false
}
```

#### Add Track to Playlist
```http
POST /api/v1/playlists/playlist_id/tracks
Authorization: Bearer your_access_token
Content-Type: application/json

{
  "provider": "itunes",
  "provider_track_id": "123456",
  "title": "Song Title",
  "artist": "Artist Name",
  "album": "Album Name"
}
```

#### Get User Playlists
```http
GET /api/v1/playlists?mine=true
Authorization: Bearer your_access_token
```

### Library Management

#### Add to Favorites
```http
POST /api/v1/favorites
Authorization: Bearer your_access_token
Content-Type: application/json

{
  "provider": "itunes",
  "provider_track_id": "123456",
  "title": "Song Title",
  "artist": "Artist Name"
}
```

#### Get Listening History
```http
GET /api/v1/history?page=1&size=20
Authorization: Bearer your_access_token
```

## 🔧 Configuration

### Environment Variables

All configuration is done through environment variables with the `MUSIC_APP_` prefix:

#### Application Settings
```env
MUSIC_APP_APP_NAME=music-app-backend
MUSIC_APP_APP_ENVIRONMENT=development
MUSIC_APP_APP_LOG_LEVEL=info
```

#### Server Configuration
```env
MUSIC_APP_SERVER_PORT=8080
MUSIC_APP_SERVER_READ_TIMEOUT=30s
MUSIC_APP_SERVER_WRITE_TIMEOUT=30s
```

#### Database Settings
```env
MUSIC_APP_DATABASE_HOST=localhost
MUSIC_APP_DATABASE_PORT=5432
MUSIC_APP_DATABASE_USER=postgres
MUSIC_APP_DATABASE_PASSWORD=your_password
MUSIC_APP_DATABASE_NAME=music_app
```

#### Authentication Settings
```env
MUSIC_APP_AUTH_JWT_ACCESS_SECRET=your_secret_key
MUSIC_APP_AUTH_JWT_REFRESH_SECRET=your_refresh_key
MUSIC_APP_AUTH_ACCESS_TOKEN_TTL=15m
MUSIC_APP_AUTH_REFRESH_TOKEN_TTL=720h
```

#### Music Providers
```env
MUSIC_APP_PROVIDERS_ENABLED=itunes,spotify
MUSIC_APP_GOOGLE_CLIENT_ID=your_google_client_id
MUSIC_APP_SPOTIFY_CLIENT_ID=your_spotify_client_id
MUSIC_APP_SPOTIFY_CLIENT_SECRET=your_spotify_secret
```

## 🏃‍♂️ Development

### Available Commands

```bash
make help           # Show all available commands
make dev            # Run in development mode
make build          # Build the application
make test           # Run tests
make test-coverage  # Run tests with coverage
make lint           # Run linter
make docker-up      # Start with Docker
make docker-down    # Stop Docker containers
make migrate        # Run database migrations
make seed           # Seed database with sample data
```

### Project Structure

```
├── cmd/server/          # Application entry point
├── internal/
│   ├── auth/           # Authentication & JWT
│   ├── config/         # Configuration management
│   ├── library/        # Favorites, history, downloads
│   ├── middleware/     # HTTP middleware
│   ├── music/          # Music provider interfaces
│   ├── playlist/       # Playlist management
│   ├── server/         # HTTP server setup
│   ├── storage/        # Database & Redis
│   ├── transport/http/ # HTTP handlers & DTOs
│   └── user/           # User management
├── pkg/                # Shared packages
├── api/                # OpenAPI specifications
├── deploy/             # Docker & deployment files
├── scripts/            # Development scripts
└── migrations/         # Database migrations
```

### Code Quality

The project follows Go best practices:

- **Clean Architecture** with dependency injection
- **Interface-based design** for testability
- **Error handling** with structured errors
- **Logging** with structured fields and request tracing
- **Testing** with table-driven tests and mocks
- **Documentation** with OpenAPI/Swagger specs

## 🧪 Testing

### Running Tests
```bash
make test                # Run all tests
make test-coverage      # Run with coverage report
```

### Test Structure
- **Unit Tests**: Individual component testing
- **Integration Tests**: HTTP endpoint testing
- **Contract Tests**: OpenAPI specification compliance

## 📦 Deployment

### Docker Production Deployment

1. **Build Production Image**
   ```bash
   make docker-build
   ```

2. **Configure Environment**
   ```bash
   # Set production environment variables
   export MUSIC_APP_APP_ENVIRONMENT=production
   export MUSIC_APP_AUTH_JWT_ACCESS_SECRET=your_production_secret
   # ... other production configs
   ```

3. **Deploy with Docker Compose**
   ```bash
   docker-compose -f deploy/docker-compose.yml up -d
   ```

### Environment-Specific Configurations

#### Development
- Debug logging enabled
- Hot reload with `make dev`
- Test database seeding

#### Production
- JSON logging
- Secure headers enabled
- Performance optimizations
- Health checks configured

## 🔒 Security Considerations

### Authentication Security
- **Argon2id Password Hashing** with configurable parameters
- **JWT Token Rotation** with refresh token families
- **Google ID Token Verification** with proper audience validation
- **Rate Limiting** on authentication endpoints

### API Security
- **CORS** configuration for frontend origins
- **Security Headers** (CSP, HSTS, etc.)
- **Input Validation** with struct tags and middleware
- **Request/Response Logging** for audit trails

### Data Protection
- **Password fields** excluded from JSON serialization
- **Sensitive data** logged with appropriate levels
- **Database connections** with SSL in production
- **Environment variables** for secrets management

## 🎵 Music Provider Integration

### iTunes Search API
- **No API key required**
- **Search, lookup, and browse** functionality
- **High-quality artwork** URL optimization
- **Rate limiting** and error handling

### Spotify Web API (Optional)
- **OAuth client credentials** or **Authorization Code with PKCE**
- **Search, tracks, playlists** support
- **Caching** for performance
- **Configurable** via environment variables

### Adding New Providers
1. Implement the `MusicProvider` interface
2. Add provider initialization in `registry.go`
3. Configure provider settings in environment
4. Add provider to enabled list

## 📊 Monitoring & Health Checks

### Health Check Endpoint
```http
GET /healthz
```

Response includes:
- Database connectivity
- Redis availability
- Music provider status
- Overall system health

### Logging & Observability
- **Structured JSON logging** in production
- **Request ID tracing** across components
- **Performance metrics** and duration tracking
- **Error aggregation** with context

## 🤝 Contributing

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Commit changes**: `git commit -m 'Add amazing feature'`
4. **Push to branch**: `git push origin feature/amazing-feature`
5. **Open a Pull Request**

### Development Guidelines
- Follow Go conventions and best practices
- Add tests for new functionality
- Update documentation as needed
- Run linter and tests before submitting

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **Documentation**: Check the API documentation at `/docs`
- **Issues**: Report bugs and feature requests on GitHub
- **Health Check**: Monitor `/healthz` for system status

---

**Built with ❤️ using Go, Gin, PostgreSQL, Redis, and modern development practices.**
