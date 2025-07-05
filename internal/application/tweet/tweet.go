package tweet

import (
	"time"
)

type Tweet struct {
	ID        string
	UserID    string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
