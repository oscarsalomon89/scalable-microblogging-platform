package tweet

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oscarsalomon89/go-hexagonal/internal/application/tweet"
	twcontext "github.com/oscarsalomon89/go-hexagonal/pkg/context"
	"github.com/oscarsalomon89/go-hexagonal/pkg/httperrors"
	"github.com/oscarsalomon89/go-hexagonal/pkg/validator"
)

const headerUserID = "X-User-ID"

type (
	TweetUseCase interface {
		CreateTweet(ctx context.Context, tweet *tweet.Tweet) error
		GetTimeline(ctx context.Context, userID string, limit, offset int) ([]tweet.Tweet, error)
	}

	tweetHandler struct {
		useCase TweetUseCase
	}
)

func NewHandler(useCase TweetUseCase) *tweetHandler {
	return &tweetHandler{useCase: useCase}
}

func (h *tweetHandler) CreateTweet(c *gin.Context) {
	ctx := twcontext.New(c.Request)
	logger := twcontext.Logger(ctx)

	userIDStr := c.GetHeader(headerUserID)
	if userIDStr == "" {
		handleError(c, httperrors.NewSimple(httperrors.ErrBadRequest, "Missing X-User-ID header"))
		return
	}

	_, err := uuid.Parse(userIDStr)
	if err != nil {
		handleError(c, httperrors.NewSimple(httperrors.ErrBadRequest, "Invalid UUID in X-User-ID header"))
		return
	}

	var req createTweetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, httperrors.NewSimple(httperrors.ErrBadRequest, "Failed to bind JSON"))
		return
	}

	if err := validator.Validate(req); err != nil {
		logger.WithError(err).Error("Failed to validate request")
		handleError(c, httperrors.New(httperrors.ErrBadRequest, "Failed to validate request", err.Error(), nil))
		return
	}

	tweetDomain := tweet.Tweet{
		UserID:  userIDStr,
		Content: req.Content,
	}
	if err := h.useCase.CreateTweet(c.Request.Context(), &tweetDomain); err != nil {
		logger.WithError(err).Error("Failed to create tweet")
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, createTweetResponse{
		Message: "Tweet created successfully",
		TweetID: tweetDomain.ID,
	})
}

func (h *tweetHandler) GetTimeline(c *gin.Context) {
	ctx := twcontext.New(c.Request)
	logger := twcontext.Logger(ctx)

	userIDStr := c.GetHeader(headerUserID)
	if userIDStr == "" {
		handleError(c, httperrors.NewSimple(httperrors.ErrBadRequest, "Missing X-User-ID header"))
		return
	}

	limit, offset := getPaginationParams(c)

	tweets, err := h.useCase.GetTimeline(ctx, userIDStr, limit, offset)
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

func getPaginationParams(c *gin.Context) (int, int) {
	var err error
	limit := 100
	offset := 0

	if limitParam := c.Query("limit"); limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			limit = 100
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
