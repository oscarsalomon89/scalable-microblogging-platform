package tweet

import (
	"time"

	"github.com/google/uuid"
	"github.com/oscarsalomon89/scalable-microblogging-platform/internal/application/tweet"
	"gorm.io/gorm"
)

type Tweet struct {
	ID        uuid.UUID      `gorm:"primaryKey;column:id"`
	UserID    string         `gorm:"column:user_id;not null"`
	Content   string         `gorm:"column:content;not null;type:text;size:280"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at"`
}

func (t *Tweet) toDomain() tweet.Tweet {
	return tweet.Tweet{
		ID:        t.ID.String(),
		UserID:    t.UserID,
		Content:   t.Content,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func fromDomain(t *tweet.Tweet) *Tweet {
	return &Tweet{
		ID:      uuid.New(),
		UserID:  t.UserID,
		Content: t.Content,
	}
}
