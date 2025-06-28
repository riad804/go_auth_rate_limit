# Golang Auth API with Org Switching

## Features
- Login, Refresh, Logout
- Org switch
- JWT & Redis-backed session
- Modular (Gin + Uber Fx)

## Seeded Users
- sohel@tenbyte.com / 123456
- jane@openresty.com / 123456
- riad@openresty.com / 123456
- tanvir@tenbyte.com / 123456

## Run
```bash
cp .env.example .env
go run cmd/api/main.go
```

### OR

## Docker Run
```bash
docker compose up --build
```
