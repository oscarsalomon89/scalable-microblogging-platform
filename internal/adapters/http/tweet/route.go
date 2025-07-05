package tweet

import "github.com/gin-gonic/gin"

const (
	tweetPath = "/tweets"
)

type TweetHandlerRouter struct {
	hdl *handler
}

func NewRouter(hdl *handler) *TweetHandlerRouter {
	return &TweetHandlerRouter{
		hdl: hdl,
	}
}

func (r *TweetHandlerRouter) AddRoutes(router *gin.RouterGroup) {
	router.POST(tweetPath, r.hdl.CreateTweet)
	router.GET(tweetPath+"/timeline", r.hdl.GetTimeline)
}
