# Go Auth API with Organization Switching

A modular authentication API built with Go, Gin, Uber Fx, GORM, MySQL, and Redis. Supports JWT authentication, organization switching, and rate limiting.

## Features
- User login, refresh, and logout
- Organization switching per user
- JWT-based authentication
- Redis-backed session and rate limiting
- Modular architecture (Gin + Uber Fx)
- MySQL database with GORM
- Swagger/OpenAPI documentation

## Seeded Users
- sohel@tenbyte.com / 123456
- jane@openresty.com / 123456
- riad@openresty.com / 123456
- tanvir@tenbyte.com / 123456

## Getting Started

### 1. Clone and Configure
```bash
git clone <your-repo-url>
cd go_auth
cp .env.example .env # or create your own .env
```

### 2. Run with Docker
```bash
docker compose up --build
```

- The API will be available at `http://localhost:8080`
- MySQL and Redis will be started as services

### 3. Run Locally (without Docker)
Make sure MySQL and Redis are running and your `.env` is configured.
```bash
go run cmd/server/main.go
```

## API Endpoints

| Method | Path              | Description                |
|--------|-------------------|----------------------------|
| POST   | /login            | User login                 |
| POST   | /refresh          | Refresh JWT tokens         |
| POST   | /logout           | Logout user                |
| GET    | /me               | Get current user info      |
| POST   | /orgs/switch      | Switch current organization|
| GET    | /health           | Health check               |

## Swagger/OpenAPI Docs
- After running the app, visit: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
- Docs are generated using [swaggo/swag](https://github.com/swaggo/swag)

## Project Structure
```
cmd/server/main.go         # Application entrypoint
internal/app/              # Server and lifecycle
internal/handlers/         # HTTP handlers
internal/service/          # Business logic
internal/repository/       # Data access
internal/models/           # GORM models
internal/middlewares/      # Gin middlewares
internal/utils/            # Utilities (JWT, etc)
pkg/database/              # MySQL connection
pkg/redis/                 # Redis connection
scripts/                   # Seed scripts
```

## Environment Variables
See `.env.example` for all configuration options.

## Database
- MySQL is used for persistent storage
- Redis is used for session and rate limiting

## Database Seeding

To insert initial organizations and users into your database, use the provided seed script.

### Seed with Docker Compose

After your database is up, run:

```bash
docker compose run --rm seed
```

This will wait for MySQL to be ready and then execute `scripts/seed.go` inside a container.

### Seed Locally (without Docker)

Make sure your `.env` and database are configured, then run:

```bash
go run scripts/seed.go
```

> **Note:** Only run the seed script in development or test environments. It will insert demo data.

## Example API Requests (rest-api.http)

You can use the provided `rest-api.http` file with the [REST Client extension for VS Code](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) to test the API endpoints easily. Replace `YOUR_ACCESS_TOKEN` and `YOUR_REFRESH_TOKEN` with actual values from your login response.

```http
@baseUrl = http://localhost:8080

### User Login
POST {{baseUrl}}/login
Content-Type: application/json

{
  "email": "riad@openresty.com",
  "password": "123456"
}

### Refresh Token
POST {{baseUrl}}/refresh
Content-Type: application/json

{
  "refresh_token": "YOUR_REFRESH_TOKEN"
}

### Logout
POST {{baseUrl}}/logout
Content-Type: application/json

{
  "refresh_token": "YOUR_REFRESH_TOKEN"
}

### Get Current User Info
GET {{baseUrl}}/me
Authorization: Bearer YOUR_ACCESS_TOKEN

### Switch Organization
POST {{baseUrl}}/orgs/switch
Authorization: Bearer YOUR_ACCESS_TOKEN
Content-Type: application/json

{
  "org_id": "org2"
}

### Health Check
GET {{baseUrl}}/health

### Swagger Docs
GET {{baseUrl}}/swagger/index.html
```