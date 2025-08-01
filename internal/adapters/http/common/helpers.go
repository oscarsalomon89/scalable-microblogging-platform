package common

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oscarsalomon89/scalable-microblogging-platform/pkg/httperrors"
	"github.com/oscarsalomon89/scalable-microblogging-platform/pkg/validator"
)

const headerUserID = "X-User-ID"

func BindAndValidate[T any](c *gin.Context) (T, error) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		return req, httperrors.New(httperrors.ErrBadRequest, "Failed to bind JSON", err.Error(), nil)
	}

	if err := Validate(req); err != nil {
		return req, err
	}

	return req, nil
}

func Validate[T any](req T) error {
	if err := validator.Validate(req); err != nil {
		return httperrors.New(httperrors.ErrBadRequest, "Failed to validate request", err.Error(), nil)
	}

	return nil
}

func ValidateUserID(c *gin.Context) (string, error) {
	userIDStr := c.GetHeader(headerUserID)
	if userIDStr == "" {
		return "", httperrors.NewSimple(httperrors.ErrBadRequest, "Missing X-User-ID header")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", httperrors.NewSimple(httperrors.ErrBadRequest, "Invalid UUID in X-User-ID header")
	}

	return userID.String(), nil
}
