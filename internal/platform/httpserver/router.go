package httpserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func NewHTTPGinServer() *http.Server {
	port := os.Getenv("WEB_SERVER_PORT")

	if port == "" {
		port = "8080"
	}

	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: gin.Default(),
	}
}

func StartServer(lc fx.Lifecycle, srv *http.Server) *http.Server {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("Error running server: %v", err)
				}
			}()
			log.Println("Server running on port 8080")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down server...")
			shutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			return srv.Shutdown(shutCtx)
		},
	})

	return srv
}
