package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/oscarsalomon89/go-hexagonal/internal/application/user"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"primaryKey;column:id"`
	Username  string         `gorm:"column:username;unique;not null"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at"`
}

type Follow struct {
	FollowerID uuid.UUID `gorm:"type:uuid;primaryKey"`
	FolloweeID uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt  time.Time `gorm:"type:timestamp with time zone;not null;default:now()"`
}

func fromDomain(u *user.User) *User {
	return &User{
		ID:       uuid.New(),
		Username: u.Username,
	}
}
