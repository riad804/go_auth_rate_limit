package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/riad804/go_auth/internal/app"
	"github.com/riad804/go_auth/internal/config"
	"github.com/riad804/go_auth/internal/handlers"
	"github.com/riad804/go_auth/internal/repository"
	"github.com/riad804/go_auth/internal/service"
	"github.com/riad804/go_auth/internal/utils"
	"github.com/riad804/go_auth/pkg/database"
	r "github.com/riad804/go_auth/pkg/redis"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func main() {
	fx.New(
		fx.Provide(
			config.LoadConfig,
			database.NewMySQLDB,
			r.NewRedisClient,
			func(r *r.RedisClient) *redis.Client {
				return r.Client
			},
			func(db *gorm.DB) repository.UserRepository {
				return repository.NewUserRepository(db)
			},
			repository.NewTokenRepository,
			service.NewAuthService,
			service.NewUserService,
			handlers.NewAuthHandler,
			utils.NewJWTWrapper,
			app.NewServer,
		),
		fx.Invoke(app.RegisterHooks),
	).Run()
}
