package modules

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oscarsalomon89/go-hexagonal/internal/platform/config"
	"github.com/oscarsalomon89/go-hexagonal/internal/platform/httpserver"
	db "github.com/oscarsalomon89/go-hexagonal/internal/platform/pg"
	pkgredis "github.com/oscarsalomon89/go-hexagonal/internal/platform/redis"
	"go.uber.org/fx"
)

var internalFactories = fx.Provide(
	db.NewDBConnections,
	pkgredis.NewRedisConnection,
	func(cfg config.Cache) time.Duration { return cfg.TTL },
	httpserver.NewHTTPGinServer,
	func(server *http.Server, cfg config.Configuration) *gin.RouterGroup {
		router := server.Handler.(*gin.Engine).Group("/api/" + cfg.APIVersion)

		router.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
			})
		})
		return router
	},
)

var internalModule = fx.Options(
	internalFactories,
	fx.Invoke(
		httpserver.StartServer,
	),
)
