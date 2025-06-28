package userhdl

import "github.com/gin-gonic/gin"

const (
	userPath = "/users"
)

type UserHandlerRouter struct {
	userhdl *handler
}

func NewRouter(userhdl *handler) *UserHandlerRouter {
	return &UserHandlerRouter{
		userhdl: userhdl,
	}
}

func (r *UserHandlerRouter) AddRoutesV1(v1 *gin.RouterGroup) {
	v1.POST(userPath, r.userhdl.CreateUser)
}
