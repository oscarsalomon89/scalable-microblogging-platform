package user

import (
	"strings"

	"github.com/oscarsalomon89/go-hexagonal/internal/application/user"
)

type createUserRequest struct {
	Username string `json:"username,omitempty" validate:"required"`
}

func (c *createUserRequest) ToDomain() *user.User {
	return &user.User{
		Username: strings.TrimSpace(c.Username),
	}
}

type createUserResponse struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}

type followUserRequest struct {
	FolloweeID string `json:"followee_id" validate:"required,validUUIDFormat"`
}

type followUserResponse struct {
	Message string `json:"message"`
}
