package itemhdl

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	ItemUseCase interface {
		CreateItem(name string) error
	}

	itemHandler struct {
		useCase ItemUseCase
	}
)

func NewHandler(useCase ItemUseCase) *itemHandler {
	return &itemHandler{useCase: useCase}
}

func (h *itemHandler) CreateItem(c *gin.Context) {
	var body struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.useCase.CreateItem(body.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Item criado com sucesso"})
}
