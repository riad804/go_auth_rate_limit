# Golang Auth API with Org Switching

## Features
- Login, Refresh, Logout
- Org switch
- JWT & Redis-backed session
- Modular (Gin + Uber Fx)

## Seeded Users
- sohel@tenbyte.com / 123456
- admin@openresty.com / adminpass

## Run
```bash
cp .env.example .env
go run cmd/api/main.go
```
