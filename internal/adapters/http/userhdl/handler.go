package userhdl

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oscarsalomon89/go-hexagonal/internal/application/user"
	pkgerrors "github.com/oscarsalomon89/go-hexagonal/pkg/errors"
	"github.com/oscarsalomon89/go-hexagonal/pkg/validator"
)

type (
	UserUseCase interface {
		CreateUser(ctx context.Context, user *user.User) error
	}

	handler struct {
		usecase UserUseCase
	}
)

func NewHandler(usecase UserUseCase) *handler {
	return &handler{usecase: usecase}
}

func (h *handler) CreateUser(c *gin.Context) {
	var req CreateUser
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, pkgerrors.NewError(pkgerrors.ErrValidation, "Invalid input", err))
		return
	}

	if err := validator.Validate(req); err != nil {
		handleError(c, err)
		return
	}

	err := h.usecase.CreateUser(c, req.ToDomain())
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Usuário criado com sucesso"})
}
