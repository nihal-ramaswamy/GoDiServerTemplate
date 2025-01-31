package fx_utils

import (
	"context"

	serverconfig "example.go_fx_gin/internal/config/server"
	middleware "example.go_fx_gin/internal/middleware/log"
	"example.go_fx_gin/internal/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func newServerEngine(lc fx.Lifecycle, config *serverconfig.Config, log *zap.Logger) *gin.Engine {
	gin.SetMode(config.GinMode)

	server := gin.Default()
	server.Use(cors.New(config.Cors))
	server.Use(middleware.DefaultStructuredLogger(log))
	server.Use(gin.Recovery())

	utils.NewRoutes(server)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Starting server on port", zap.String("port", config.Port))

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping server")
			defer log.Sync()

			return nil
		},
	})

	return server
}

var serverModule = fx.Module(
	"serverModule",
	fx.Provide(newServerEngine),
)
