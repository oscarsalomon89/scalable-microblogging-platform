package tweet

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oscarsalomon89/go-hexagonal/internal/adapters/http/common"
	"github.com/oscarsalomon89/go-hexagonal/internal/application/tweet"
	twcontext "github.com/oscarsalomon89/go-hexagonal/pkg/context"
)

type (
	TweetUseCase interface {
		CreateTweet(ctx context.Context, tweet *tweet.Tweet) error
		GetTimeline(ctx context.Context, userID string, limit, offset int) ([]tweet.Tweet, error)
	}

	handler struct {
		usecase TweetUseCase
	}
)

func NewHandler(useCase TweetUseCase) *handler {
	return &handler{usecase: useCase}
}

func (h *handler) CreateTweet(c *gin.Context) {
	ctx := twcontext.New(c.Request)
	logger := twcontext.Logger(ctx)

	userID, err := common.ValidateUserID(c)
	if err != nil {
		logger.WithError(err).Error("Failed to validate user ID")
		handleError(c, err)
		return
	}

	req, err := common.BindAndValidate[createTweetRequest](c)
	if err != nil {
		logger.WithError(err).Error("Failed to bind JSON")
		handleError(c, err)
		return
	}

	tweetDomain := tweet.Tweet{
		UserID:  userID,
		Content: req.Content,
	}
	if err := h.usecase.CreateTweet(ctx, &tweetDomain); err != nil {
		logger.WithError(err).Error("Failed to create tweet")
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, createTweetResponse{
		Message: "Tweet created successfully",
		TweetID: tweetDomain.ID,
	})
}

func (h *handler) GetTimeline(c *gin.Context) {
	ctx := twcontext.New(c.Request)
	logger := twcontext.Logger(ctx)

	userID, err := common.ValidateUserID(c)
	if err != nil {
		logger.WithError(err).Error("Failed to validate user ID")
		handleError(c, err)
		return
	}

	limit, offset := parsePaginationParams(c)

	tweets, err := h.usecase.GetTimeline(ctx, userID, limit, offset)
	if err != nil {
		logger.WithError(err).Error("Failed to get timeline")
		handleError(c, err)
		return
	}

	response := make([]tweetsResponse, len(tweets))
	for i, tweet := range tweets {
		response[i] = tweetsResponse{
			ID:        tweet.ID,
			UserID:    tweet.UserID,
			Content:   tweet.Content,
			CreatedAt: tweet.CreatedAt,
			UpdatedAt: tweet.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

const defaultLimit = 100

func parsePaginationParams(c *gin.Context) (int, int) {
	var err error
	limit := defaultLimit
	offset := 0

	if limitParam := c.Query("limit"); limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			limit = defaultLimit
		}
	}

	if offsetParam := c.Query("offset"); offsetParam != "" {
		offset, err = strconv.Atoi(offsetParam)
		if err != nil {
			offset = 0
		}
	}

	return limit, offset
}
