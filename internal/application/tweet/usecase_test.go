package tweet_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/oscarsalomon89/go-hexagonal/internal/application/tweet"
	"github.com/oscarsalomon89/go-hexagonal/internal/application/tweet/mocks"
	"github.com/oscarsalomon89/go-hexagonal/internal/application/user"
	twcontext "github.com/oscarsalomon89/go-hexagonal/pkg/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type dependencies struct {
	userFinder    *mocks.UserFinder
	tweetReader   *mocks.TweetReader
	tweetsCreator *mocks.TweetCreator
	cache         *mocks.TimelineCache
}

func init() {
	twcontext.NewLogger()
}

func Test_usecase_CreateTweet(t *testing.T) {
	type input struct {
		ctx   context.Context
		tweet *tweet.Tweet
	}

	type output struct {
		err error
	}

	tests := []struct {
		name         string
		input        input
		output       output
		dependencies func(in input, d *dependencies)
		assert       func(t *testing.T, expected, actual output)
	}{
		{
			name: "should return error if userFinder.ExistsByID returns error",
			input: input{
				ctx:   twcontext.NewTestContext(),
				tweet: &tweet.Tweet{UserID: "u1"},
			},
			output: output{err: fmt.Errorf("failed to check user ID: %w", assert.AnError)},
			dependencies: func(in input, d *dependencies) {
				d.userFinder.On("ExistsByID", in.ctx, in.tweet.UserID).Return(false, assert.AnError)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected.err.Error(), actual.err.Error())
			},
		},
		{
			name: "should return error if user does not exist",
			input: input{
				ctx:   twcontext.NewTestContext(),
				tweet: &tweet.Tweet{UserID: "u1"},
			},
			output: output{err: user.ErrUserNotFound},
			dependencies: func(in input, d *dependencies) {
				d.userFinder.On("ExistsByID", in.ctx, in.tweet.UserID).Return(false, nil)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "should return error if tweetsCreator.CreateTweet returns error",
			input: input{
				ctx:   twcontext.NewTestContext(),
				tweet: &tweet.Tweet{UserID: "u1"},
			},
			output: output{err: fmt.Errorf("failed to create tweet: %w", assert.AnError)},
			dependencies: func(in input, d *dependencies) {
				d.userFinder.On("ExistsByID", in.ctx, in.tweet.UserID).Return(true, nil)
				d.tweetsCreator.On("CreateTweet", in.ctx, in.tweet).Return(assert.AnError)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected.err.Error(), actual.err.Error())
			},
		},
		{
			name: "should create tweet successfully",
			input: input{
				ctx:   twcontext.NewTestContext(),
				tweet: &tweet.Tweet{UserID: "u1"},
			},
			output: output{err: nil},
			dependencies: func(in input, d *dependencies) {
				d.userFinder.On("ExistsByID", in.ctx, in.tweet.UserID).Return(true, nil)
				d.tweetsCreator.On("CreateTweet", in.ctx, in.tweet).Return(nil)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected, actual)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &dependencies{
				userFinder:    mocks.NewUserFinder(t),
				tweetReader:   mocks.NewTweetReader(t),
				tweetsCreator: mocks.NewTweetCreator(t),
				cache:         mocks.NewTimelineCache(t),
			}

			// Synchronize with the goroutine
			var wg sync.WaitGroup
			if tt.name == "should create tweet successfully" {
				followers := []string{"f1", "f2"}
				ctx := twcontext.NewDetachedWithRequestID(tt.input.ctx)
				d.userFinder.On("GetFollowers", ctx, tt.input.tweet.UserID).Return(followers, nil)
				wg.Add(len(followers))
				for _, follower := range followers {
					d.cache.On("InvalidateTimeline", ctx, follower).Return(nil).Run(func(args mock.Arguments) {
						wg.Done()
					})
				}
			}
			tt.dependencies(tt.input, d)

			uc := tweet.NewTweetUseCase(d.userFinder, d.tweetReader, d.tweetsCreator, d.cache)
			var actual output
			actual.err = uc.CreateTweet(tt.input.ctx, tt.input.tweet)

			// Wait for the goroutine to finish
			done := make(chan struct{})
			go func() {
				wg.Wait()
				close(done)
			}()
			if done != nil {
				select {
				case <-done:
				case <-time.After(2 * time.Second):
					t.Error("goroutine did not finish in time")
				}
			}

			tt.assert(t, tt.output, actual)
		})
	}
}

func Test_usecase_GetTimeline(t *testing.T) {
	type input struct {
		ctx    context.Context
		userID string
		limit  int
		offset int
	}

	type output struct {
		err    error
		tweets []tweet.Tweet
	}

	tests := []struct {
		name         string
		input        input
		output       output
		dependencies func(in input, d *dependencies)
		assert       func(t *testing.T, expected, actual output)
	}{
		{
			name: "should return tweets from cache if cache hit and not empty",
			input: input{
				ctx:    twcontext.NewTestContext(),
				userID: "u1",
				limit:  10,
				offset: 0,
			},
			output: output{tweets: []tweet.Tweet{{ID: "t1"}, {ID: "t2"}}, err: nil},
			dependencies: func(in input, d *dependencies) {
				d.cache.On("GetTimeline", in.ctx, in.userID).Return([]tweet.Tweet{{ID: "t1"}, {ID: "t2"}}, nil)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "should return empty slice if cache hit but empty",
			input: input{
				ctx:    twcontext.NewTestContext(),
				userID: "u1",
				limit:  10,
				offset: 0,
			},
			output: output{tweets: []tweet.Tweet{}, err: nil},
			dependencies: func(in input, d *dependencies) {
				d.cache.On("GetTimeline", in.ctx, in.userID).Return([]tweet.Tweet{}, nil)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "should return error if userFinder.GetFollowees returns error",
			input: input{
				ctx:    twcontext.NewTestContext(),
				userID: "u1",
				limit:  10,
				offset: 0,
			},
			output: output{tweets: nil, err: fmt.Errorf("failed to get followees: %w", assert.AnError)},
			dependencies: func(in input, d *dependencies) {
				d.cache.On("GetTimeline", in.ctx, in.userID).Return(nil, assert.AnError)
				d.userFinder.On("GetFollowees", in.ctx, in.userID).Return(nil, assert.AnError)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected.err.Error(), actual.err.Error())
			},
		},
		{
			name: "should return empty slice if user has no followees",
			input: input{
				ctx:    twcontext.NewTestContext(),
				userID: "u1",
				limit:  10,
				offset: 0,
			},
			output: output{tweets: []tweet.Tweet{}, err: nil},
			dependencies: func(in input, d *dependencies) {
				d.cache.On("GetTimeline", in.ctx, in.userID).Return(nil, assert.AnError)
				d.userFinder.On("GetFollowees", in.ctx, in.userID).Return([]string{}, nil)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "should return error if tweetReader.GetTweetsByUserIDs returns error",
			input: input{
				ctx:    twcontext.NewTestContext(),
				userID: "u1",
				limit:  10,
				offset: 0,
			},
			output: output{tweets: nil, err: fmt.Errorf("error retrieving timeline from Cassandra: %w", assert.AnError)},
			dependencies: func(in input, d *dependencies) {
				d.cache.On("GetTimeline", in.ctx, in.userID).Return(nil, assert.AnError)
				d.userFinder.On("GetFollowees", in.ctx, in.userID).Return([]string{"f1", "f2"}, nil)
				d.tweetReader.On("GetTweetsByUserIDs", in.ctx, []string{"f1", "f2"}, in.limit, in.offset).Return(nil, assert.AnError)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected.err.Error(), actual.err.Error())
			},
		},
		{
			name: "should return empty slice if DB returns empty",
			input: input{
				ctx:    twcontext.NewTestContext(),
				userID: "u1",
				limit:  10,
				offset: 0,
			},
			output: output{tweets: []tweet.Tweet{}, err: nil},
			dependencies: func(in input, d *dependencies) {
				d.cache.On("GetTimeline", in.ctx, in.userID).Return(nil, assert.AnError)
				d.userFinder.On("GetFollowees", in.ctx, in.userID).Return([]string{"f1", "f2"}, nil)
				d.tweetReader.On("GetTweetsByUserIDs", in.ctx, []string{"f1", "f2"}, in.limit, in.offset).Return([]tweet.Tweet{}, nil)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "should return tweets from DB and set cache",
			input: input{
				ctx:    twcontext.NewTestContext(),
				userID: "u1",
				limit:  10,
				offset: 0,
			},
			output: output{tweets: []tweet.Tweet{{ID: "t1"}, {ID: "t2"}}, err: nil},
			dependencies: func(in input, d *dependencies) {
				d.cache.On("GetTimeline", in.ctx, in.userID).Return(nil, assert.AnError)
				d.userFinder.On("GetFollowees", in.ctx, in.userID).Return([]string{"f1", "f2"}, nil)
				d.tweetReader.On("GetTweetsByUserIDs", in.ctx, []string{"f1", "f2"}, in.limit, in.offset).Return([]tweet.Tweet{{ID: "t1"}, {ID: "t2"}}, nil)
				d.cache.On("SetTimeline", in.ctx, in.userID, []tweet.Tweet{{ID: "t1"}, {ID: "t2"}}).Return(nil)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected, actual)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &dependencies{
				userFinder:    mocks.NewUserFinder(t),
				tweetReader:   mocks.NewTweetReader(t),
				tweetsCreator: mocks.NewTweetCreator(t),
				cache:         mocks.NewTimelineCache(t),
			}
			tt.dependencies(tt.input, d)

			uc := tweet.NewTweetUseCase(d.userFinder, d.tweetReader, d.tweetsCreator, d.cache)
			var actual output
			actual.tweets, actual.err = uc.GetTimeline(tt.input.ctx, tt.input.userID, tt.input.limit, tt.input.offset)
			tt.assert(t, tt.output, actual)
		})
	}
}
