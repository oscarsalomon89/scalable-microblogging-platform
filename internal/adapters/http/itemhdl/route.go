package itemhdl

import "github.com/gin-gonic/gin"

const (
	itemPath = "/items"
)

type ItemHandlerRouter struct {
	itemhdl *itemHandler
}

func NewRouter(itemhdl *itemHandler) *ItemHandlerRouter {
	return &ItemHandlerRouter{
		itemhdl: itemhdl,
	}
}

func (r *ItemHandlerRouter) AddRoutesV1(v1 *gin.RouterGroup) {
	v1.POST(itemPath, r.itemhdl.CreateItem)
}
