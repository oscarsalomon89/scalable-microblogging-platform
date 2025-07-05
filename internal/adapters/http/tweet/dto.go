package tweet

import "time"

type createTweetRequest struct {
	Content string `json:"content" validate:"required,max=280"`
}

type createTweetResponse struct {
	Message string `json:"message"`
	TweetID string `json:"tweet_id"`
}

type tweetsResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
