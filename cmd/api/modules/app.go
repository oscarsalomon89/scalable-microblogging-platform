package modules

import (
	"github.com/gin-gonic/gin"
	"github.com/oscarsalomon89/go-hexagonal/internal/adapters/http/userhdl"
	"github.com/oscarsalomon89/go-hexagonal/internal/adapters/postgres/userrepo"
	"github.com/oscarsalomon89/go-hexagonal/internal/application/user"
	"github.com/oscarsalomon89/go-hexagonal/internal/platform/config"
	"github.com/oscarsalomon89/go-hexagonal/internal/platform/db"
	"github.com/oscarsalomon89/go-hexagonal/internal/platform/environment"
	"github.com/oscarsalomon89/go-hexagonal/internal/platform/httpserver"
	clonctx "github.com/oscarsalomon89/go-hexagonal/pkg/context"
	"github.com/oscarsalomon89/go-hexagonal/pkg/validator"
	"go.uber.org/fx"
)

func NewApp() *fx.App {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	clonctx.NewLogger()

	options := []fx.Option{
		fx.Provide(func() config.Database { return cfg.Database }),
		fx.Provide(func() config.Configuration { return cfg }),
		fx.Provide(
			db.NewDBConnections,
			newRouter,
			fx.Annotate(
				userrepo.NewUserRepository,
				fx.As(new(user.Repository)),
			),
			fx.Annotate(
				user.NewUserUseCase,
				fx.As(new(userhdl.UserUseCase)),
			),
			userhdl.NewHandler,
			userhdl.NewRouter,
			httpserver.NewHTTPGinServer,
		),
	}

	return fx.New(
		fx.Options(options...),
		fx.Invoke(validator.RegisterValidation),
		fx.Invoke(httpserver.StartServer),
		fx.Invoke(runMigrations),
		fx.Invoke(registerRouter),
	)
}

func newRouter() *gin.Engine {
	return gin.Default()
}

func registerRouter(router *gin.Engine, hdl *userhdl.UserHandlerRouter) error {
	routerGroup := router.Group("/v1/api")
	hdl.AddRoutesV1(routerGroup)
	return nil
}

func runMigrations(cfg config.Configuration) error {
	if environment.GetFromString(cfg.Scope) != environment.Local {
		return nil
	}

	return db.RunMigrations(cfg.Database)
}
