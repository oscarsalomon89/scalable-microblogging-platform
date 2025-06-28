package userrepo

import (
	"context"

	"github.com/oscarsalomon89/go-hexagonal/internal/application/user"
	"github.com/oscarsalomon89/go-hexagonal/internal/platform/db"
)

type userRepository struct {
	db db.Connections
}

func NewUserRepository(db db.Connections) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *user.User) error {
	return r.db.MasterConn.Create(user).Error
}
