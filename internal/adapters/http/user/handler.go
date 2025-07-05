package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oscarsalomon89/go-hexagonal/internal/application/user"
	twcontext "github.com/oscarsalomon89/go-hexagonal/pkg/context"
	"github.com/oscarsalomon89/go-hexagonal/pkg/httperrors"
	"github.com/oscarsalomon89/go-hexagonal/pkg/validator"
)

const headerUserID = "X-User-ID"

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

	req, err := bindAndValidate[createUserRequest](c)
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

	userIDStr := c.GetHeader(headerUserID)
	if userIDStr == "" {
		handleError(c, httperrors.NewSimple(httperrors.ErrBadRequest, "Missing X-User-ID header"))
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		handleError(c, httperrors.NewSimple(httperrors.ErrBadRequest, "Invalid UUID in X-User-ID header"))
		return
	}

	req, err := bindAndValidate[followUserRequest](c)
	if err != nil {
		logger.WithError(err).Error("Failed to bind JSON")
		handleError(c, err)
		return
	}

	if err := h.usecase.FollowUser(ctx, userID.String(), req.FolloweeID); err != nil {
		logger.WithError(err).Error("Failed to follow user")
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, followUserResponse{
		Message: "User followed successfully",
	})
}

func bindAndValidate[T any](c *gin.Context) (T, error) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		return req, httperrors.New(httperrors.ErrBadRequest, "Failed to bind JSON", err.Error(), nil)
	}

	if err := validator.Validate(req); err != nil {
		return req, httperrors.New(httperrors.ErrBadRequest, "Failed to validate request", err.Error(), nil)
	}

	return req, nil
}
