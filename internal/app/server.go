package app

import (
	"context"
	"net/http"
	"time"

	"github.com/riad804/go_auth/internal/config"
	"github.com/riad804/go_auth/internal/handlers"
	middleware "github.com/riad804/go_auth/internal/middlewares"
	"github.com/riad804/go_auth/internal/service"
	"github.com/riad804/go_auth/internal/utils"
	"github.com/riad804/go_auth/pkg/redis"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	router      *gin.Engine
	httpServer  *http.Server
	authHandler *handlers.AuthHandler
	rateLimit   gin.HandlerFunc
	auth        gin.HandlerFunc
}

func NewServer(
	lc fx.Lifecycle,
	cfg *config.Config,
	db *gorm.DB,
	redisClient *redis.RedisClient,
	authHandler *handlers.AuthHandler,
	authService *service.AuthService,
	jwtUtil *utils.JWTWrapper,
) *Server {
	router := gin.Default()

	rateLimit := middleware.RateLimitMiddleware(redisClient.Client, 5, time.Minute)
	authMiddleware := middleware.AuthMiddleware(jwtUtil)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	server := &Server{
		router:      router,
		httpServer:  srv,
		authHandler: authHandler,
		rateLimit:   rateLimit,
		auth:        authMiddleware,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			return srv.Shutdown(shutdownCtx)
		},
	})

	return server
}

func (s *Server) RegisterRoutes() {
	// Public routes
	public := s.router.Group("/")
	{
		public.POST("/login", s.rateLimit, s.authHandler.Login)
		public.POST("/refresh", s.authHandler.Refresh)
		public.POST("/logout", s.authHandler.Logout)
	}

	// Auth routes
	auth := s.router.Group("/")
	auth.Use(s.auth)
	{
		auth.GET("/me", s.authHandler.Me)
		auth.POST("/orgs/switch", s.authHandler.SwitchOrg)
	}

	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func RegisterHooks(
	lc fx.Lifecycle,
	server *Server,
) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			server.RegisterRoutes()
			return nil
		},
	})
}
