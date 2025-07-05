package tweet

import "github.com/gin-gonic/gin"

const (
	tweetPath = "/tweets"
)

type TweetHandlerRouter struct {
	tweethdl *tweetHandler
}

func NewRouter(tweethdl *tweetHandler) *TweetHandlerRouter {
	return &TweetHandlerRouter{
		tweethdl: tweethdl,
	}
}

func (r *TweetHandlerRouter) AddRoutes(router *gin.RouterGroup) {
	router.POST(tweetPath, r.tweethdl.CreateTweet)
	router.GET(tweetPath+"/timeline", r.tweethdl.GetTimeline)
}
