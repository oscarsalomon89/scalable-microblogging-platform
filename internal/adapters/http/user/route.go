package user

import "github.com/gin-gonic/gin"

const (
	userPath = "/users"
)

type UserHandlerRouter struct {
	hdl *handler
}

func NewRouter(hdl *handler) *UserHandlerRouter {
	return &UserHandlerRouter{
		hdl: hdl,
	}
}

func (r *UserHandlerRouter) AddRoutesV1(v1 *gin.RouterGroup) {
	v1.POST(userPath, r.hdl.CreateUser)
	v1.POST(userPath+"/follow", r.hdl.FollowUser)
	v1.DELETE(userPath+"/unfollow/:followeeID", r.hdl.UnfollowUser)
}
