package modules

import (
	"github.com/gin-gonic/gin"
	userhdl "github.com/oscarsalomon89/scalable-microblogging-platform/internal/adapters/http/user"
	userrepo "github.com/oscarsalomon89/scalable-microblogging-platform/internal/adapters/postgres/user"
	timelinerepo "github.com/oscarsalomon89/scalable-microblogging-platform/internal/adapters/redis/timeline"
	"github.com/oscarsalomon89/scalable-microblogging-platform/internal/application/user"
	"go.uber.org/fx"
)

var userFactories = fx.Provide(
	fx.Annotate(
		userrepo.NewUserRepository,
		fx.As(new(user.UserCreator)),
		fx.As(new(user.UserFinder)),
	),
	fx.Annotate(
		timelinerepo.NewCache,
		fx.As(new(user.TimelineCache)),
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
