package modules

import (
	"github.com/gin-gonic/gin"
	tweethdl "github.com/oscarsalomon89/scalable-microblogging-platform/internal/adapters/http/tweet"
	tweetrepo "github.com/oscarsalomon89/scalable-microblogging-platform/internal/adapters/postgres/tweet"
	userrepo "github.com/oscarsalomon89/scalable-microblogging-platform/internal/adapters/postgres/user"
	timelinerepo "github.com/oscarsalomon89/scalable-microblogging-platform/internal/adapters/redis/timeline"
	"github.com/oscarsalomon89/scalable-microblogging-platform/internal/application/tweet"
	"go.uber.org/fx"
)

var tweetFactories = fx.Provide(
	fx.Annotate(
		tweetrepo.NewTweetRepository,
		fx.As(new(tweet.TweetCreator)),
		fx.As(new(tweet.TweetReader)),
	),
	fx.Annotate(
		userrepo.NewUserRepository,
		fx.As(new(tweet.UserFinder)),
	),
	fx.Annotate(
		timelinerepo.NewCache,
		fx.As(new(tweet.TimelineCache)),
	),
	fx.Annotate(
		tweet.NewTweetUseCase,
		fx.As(new(tweethdl.TweetUseCase)),
	),
	tweethdl.NewHandler,
	tweethdl.NewRouter,
)

func registerTweetEndpoints(router *gin.RouterGroup, handler *tweethdl.TweetHandlerRouter) {
	handler.AddRoutes(router)
}

var tweetModule = fx.Options(
	fx.Invoke(
		registerTweetEndpoints,
	),
	tweetFactories,
)
