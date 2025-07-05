package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oscarsalomon89/go-hexagonal/internal/adapters/http/common"
	"github.com/oscarsalomon89/go-hexagonal/internal/application/user"
	twcontext "github.com/oscarsalomon89/go-hexagonal/pkg/context"
)

type (
	UserUseCase interface {
		CreateUser(ctx context.Context, user *user.User) error
		FollowUser(ctx context.Context, followerID, followeeID string) error
	}

	handler struct {
		usecase UserUseCase
	}
)

func NewHandler(usecase UserUseCase) *handler {
	return &handler{usecase: usecase}
}

func (h *handler) CreateUser(c *gin.Context) {
	ctx := twcontext.New(c.Request)
	logger := twcontext.Logger(ctx)

	req, err := common.BindAndValidate[createUserRequest](c)
	if err != nil {
		logger.WithError(err).Error("Failed to bind JSON")
		handleError(c, err)
		return
	}

	userDomain := req.ToDomain()
	if err := h.usecase.CreateUser(ctx, userDomain); err != nil {
		logger.WithError(err).Error("Failed to create user")
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, createUserResponse{
		Message: "User created successfully",
		UserID:  userDomain.ID,
	})
}

func (h *handler) FollowUser(c *gin.Context) {
	ctx := twcontext.New(c.Request)
	logger := twcontext.Logger(ctx)

	userID, err := common.ValidateUserID(c)
	if err != nil {
		logger.WithError(err).Error("Failed to validate user ID")
		handleError(c, err)
		return
	}

	req, err := common.BindAndValidate[followUserRequest](c)
	if err != nil {
		logger.WithError(err).Error("Failed to bind JSON")
		handleError(c, err)
		return
	}

	if err := h.usecase.FollowUser(ctx, userID, req.FolloweeID); err != nil {
		logger.WithError(err).Error("Failed to follow user")
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, followUserResponse{
		Message: "User followed successfully",
	})
}
