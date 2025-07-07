package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oscarsalomon89/scalable-microblogging-platform/internal/application/user"
	"github.com/oscarsalomon89/scalable-microblogging-platform/pkg/httperrors"
)

func handleError(c *gin.Context, err error) {
	var apiError *httperrors.APIError

	switch {
	case errors.Is(err, user.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, httperrors.NewSimple(httperrors.ErrBadRequest, "Invalid input"))
	case errors.Is(err, user.ErrCannotFollowSelf):
		c.JSON(http.StatusBadRequest, httperrors.NewSimple(httperrors.ErrBadRequest, "Cannot follow self"))
	case errors.Is(err, user.ErrUsernameExists):
		c.JSON(http.StatusBadRequest, httperrors.NewSimple(httperrors.ErrConflict, "User already exists"))
	case errors.Is(err, user.ErrUserNotFound):
		c.JSON(http.StatusNotFound, httperrors.NewSimple(httperrors.ErrNotFound, "User not found"))
	case errors.Is(err, user.ErrFolloweeNotFound):
		c.JSON(http.StatusNotFound, httperrors.NewSimple(httperrors.ErrNotFound, "Followee not found"))
	case errors.Is(err, user.ErrAlreadyFollowing):
		c.JSON(http.StatusConflict, httperrors.NewSimple(httperrors.ErrConflict, "Already following"))
	case errors.Is(err, user.ErrNotFollowing):
		c.JSON(http.StatusConflict, httperrors.NewSimple(httperrors.ErrConflict, "Not following"))
	case errors.Is(err, user.ErrCannotUnfollowSelf):
		c.JSON(http.StatusConflict, httperrors.NewSimple(httperrors.ErrConflict, "Cannot unfollow self"))
	case errors.As(err, &apiError):
		c.JSON(apiError.Code, apiError)
	default:
		c.JSON(http.StatusInternalServerError, httperrors.NewSimple(httperrors.ErrInternal, "Internal server error"))
	}
}
