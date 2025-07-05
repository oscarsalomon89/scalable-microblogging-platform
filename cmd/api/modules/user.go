package modules

import (
	"github.com/gin-gonic/gin"
	userhdl "github.com/oscarsalomon89/go-hexagonal/internal/adapters/http/user"
	userrepo "github.com/oscarsalomon89/go-hexagonal/internal/adapters/postgres/user"
	"github.com/oscarsalomon89/go-hexagonal/internal/application/user"
	"go.uber.org/fx"
)

var userFactories = fx.Provide(
	fx.Annotate(
		userrepo.NewUserRepository,
		fx.As(new(user.UsersCreator)),
		fx.As(new(user.UsersFinder)),
	),
	fx.Annotate(
		user.NewUserUseCase,
		fx.As(new(userhdl.UserUseCase)),
	),
	userhdl.NewHandler,
	userhdl.NewRouter,
)

func registerUserEndpoints(router *gin.RouterGroup, handler *userhdl.UserHandlerRouter) {
	handler.AddRoutesV1(router)
}

var userModule = fx.Options(
	fx.Invoke(
		registerUserEndpoints,
	),
	userFactories,
)
