package userhdl

import "github.com/oscarsalomon89/go-hexagonal/internal/application/user"

type CreateUser struct {
	Username string `json:"username,omitempty" validate:"required"`
	Email    string `json:"email,omitempty" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6"`
}

func (c *CreateUser) ToDomain() *user.User {
	return &user.User{
		Name:     c.Username,
		Email:    c.Email,
		Password: c.Password,
	}
}
