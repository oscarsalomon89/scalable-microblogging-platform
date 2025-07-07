package user_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/oscarsalomon89/scalable-microblogging-platform/internal/application/user"
	"github.com/oscarsalomon89/scalable-microblogging-platform/internal/application/user/mocks"
	twcontext "github.com/oscarsalomon89/scalable-microblogging-platform/pkg/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	twcontext.NewLogger()
}

func Test_userUseCase_CreateUser(t *testing.T) {
	type input struct {
		ctx  context.Context
		user *user.User
	}

	type output struct {
		err error
	}

	type dependencies struct {
		creator *mocks.UserCreator
		finder  *mocks.UserFinder
		cache   *mocks.TimelineCache
	}

	tests := []struct {
		name         string
		input        input
		output       output
		dependencies func(in input, d *dependencies)
		assert       func(t *testing.T, expected, actual output)
	}{
		{
			name: "should return error if username is empty",
			input: input{
				ctx:  twcontext.NewTestContext(),
				user: &user.User{Username: ""},
			},
			output:       output{err: user.ErrInvalidInput},
			dependencies: func(in input, d *dependencies) {},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "should return error if finder.ExistsByUsername returns error",
			input: input{
				ctx:  twcontext.NewTestContext(),
				user: &user.User{Username: "alice"},
			},
			output: output{err: fmt.Errorf("failed to check username: %w", assert.AnError)},
			dependencies: func(in input, d *dependencies) {
				d.finder.On("ExistsByUsername", in.ctx, in.user.Username).Return(false, assert.AnError)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "should return error if username already exists",
			input: input{
				ctx:  twcontext.NewTestContext(),
				user: &user.User{Username: "bob"},
			},
			output: output{err: user.ErrUsernameExists},
			dependencies: func(in input, d *dependencies) {
				d.finder.On("ExistsByUsername", in.ctx, in.user.Username).Return(true, nil)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "should return error if creator.CreateUser returns error",
			input: input{
				ctx:  twcontext.NewTestContext(),
				user: &user.User{Username: "carol"},
			},
			output: output{err: fmt.Errorf("failed to create user: %w", assert.AnError)},
			dependencies: func(in input, d *dependencies) {
				d.finder.On("ExistsByUsername", in.ctx, in.user.Username).Return(false, nil)
				d.creator.On("CreateUser", in.ctx, in.user).Return(assert.AnError)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "should create user successfully when username is unique",
			input: input{
				ctx:  twcontext.NewTestContext(),
				user: &user.User{Username: "dave"},
			},
			output: output{err: nil},
			dependencies: func(in input, d *dependencies) {
				d.finder.On("ExistsByUsername", in.ctx, in.user.Username).Return(false, nil)
				d.creator.On("CreateUser", in.ctx, in.user).Return(nil)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected, actual)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &dependencies{
				creator: mocks.NewUserCreator(t),
				finder:  mocks.NewUserFinder(t),
				cache:   mocks.NewTimelineCache(t),
			}
			tt.dependencies(tt.input, d)

			uc := user.NewUserUseCase(d.creator, d.finder, d.cache)
			var actual output
			actual.err = uc.CreateUser(tt.input.ctx, tt.input.user)

			tt.assert(t, tt.output, actual)
		})
	}
}

func Test_userUseCase_FollowUser(t *testing.T) {
	type input struct {
		ctx        context.Context
		followerID string
		followeeID string
	}

	type output struct {
		err error
	}

	type dependencies struct {
		creator *mocks.UserCreator
		finder  *mocks.UserFinder
		cache   *mocks.TimelineCache
	}

	tests := []struct {
		name         string
		input        input
		output       output
		dependencies func(in input, d *dependencies)
		assert       func(t *testing.T, expected, actual output)
	}{
		{
			name: "should return error if followerID or followeeID is empty",
			input: input{
				ctx:        twcontext.NewTestContext(),
				followerID: "",
				followeeID: "123",
			},
			output:       output{err: user.ErrInvalidInput},
			dependencies: func(in input, d *dependencies) {},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "should return error if followerID equals followeeID",
			input: input{
				ctx:        twcontext.NewTestContext(),
				followerID: "123",
				followeeID: "123",
			},
			output:       output{err: user.ErrCannotFollowSelf},
			dependencies: func(in input, d *dependencies) {},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "should return error if finder.ExistsByID for follower returns error",
			input: input{
				ctx:        twcontext.NewTestContext(),
				followerID: "f1",
				followeeID: "f2",
			},
			output: output{err: fmt.Errorf("failed to check follower with ID %s: %w", "f1", assert.AnError)},
			dependencies: func(in input, d *dependencies) {
				d.finder.On("ExistsByID", in.ctx, in.followerID).Return(false, assert.AnError)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected.err.Error(), actual.err.Error())
			},
		},
		{
			name: "should return error if follower does not exist",
			input: input{
				ctx:        twcontext.NewTestContext(),
				followerID: "f1",
				followeeID: "f2",
			},
			output: output{err: user.ErrUserNotFound},
			dependencies: func(in input, d *dependencies) {
				d.finder.On("ExistsByID", in.ctx, in.followerID).Return(false, nil)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "should return error if finder.ExistsByID for followee returns error",
			input: input{
				ctx:        twcontext.NewTestContext(),
				followerID: "f1",
				followeeID: "f2",
			},
			output: output{err: fmt.Errorf("failed to check followee with ID %s: %w", "f2", assert.AnError)},
			dependencies: func(in input, d *dependencies) {
				d.finder.On("ExistsByID", in.ctx, in.followerID).Return(true, nil)
				d.finder.On("ExistsByID", in.ctx, in.followeeID).Return(false, assert.AnError)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected.err.Error(), actual.err.Error())
			},
		},
		{
			name: "should return error if followee does not exist",
			input: input{
				ctx:        twcontext.NewTestContext(),
				followerID: "f1",
				followeeID: "f2",
			},
			output: output{err: user.ErrFolloweeNotFound},
			dependencies: func(in input, d *dependencies) {
				d.finder.On("ExistsByID", in.ctx, in.followerID).Return(true, nil)
				d.finder.On("ExistsByID", in.ctx, in.followeeID).Return(false, nil)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "should return error if finder.IsFollowing returns error",
			input: input{
				ctx:        twcontext.NewTestContext(),
				followerID: "f1",
				followeeID: "f2",
			},
			output: output{err: fmt.Errorf("error checking follow relationship: %w", assert.AnError)},
			dependencies: func(in input, d *dependencies) {
				d.finder.On("ExistsByID", in.ctx, in.followerID).Return(true, nil)
				d.finder.On("ExistsByID", in.ctx, in.followeeID).Return(true, nil)
				d.finder.On("IsFollowing", in.ctx, in.followerID, in.followeeID).Return(false, assert.AnError)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected.err.Error(), actual.err.Error())
			},
		},
		{
			name: "should return error if already following",
			input: input{
				ctx:        twcontext.NewTestContext(),
				followerID: "f1",
				followeeID: "f2",
			},
			output: output{err: user.ErrAlreadyFollowing},
			dependencies: func(in input, d *dependencies) {
				d.finder.On("ExistsByID", in.ctx, in.followerID).Return(true, nil)
				d.finder.On("ExistsByID", in.ctx, in.followeeID).Return(true, nil)
				d.finder.On("IsFollowing", in.ctx, in.followerID, in.followeeID).Return(true, nil)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "should return error if creator.FollowUser returns error",
			input: input{
				ctx:        twcontext.NewTestContext(),
				followerID: "f1",
				followeeID: "f2",
			},
			output: output{err: fmt.Errorf("error following user: %w", assert.AnError)},
			dependencies: func(in input, d *dependencies) {
				d.finder.On("ExistsByID", in.ctx, in.followerID).Return(true, nil)
				d.finder.On("ExistsByID", in.ctx, in.followeeID).Return(true, nil)
				d.finder.On("IsFollowing", in.ctx, in.followerID, in.followeeID).Return(false, nil)
				d.creator.On("FollowUser", in.ctx, in.followerID, in.followeeID).Return(assert.AnError)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected.err.Error(), actual.err.Error())
			},
		},
		{
			name: "should follow user successfully",
			input: input{
				ctx:        twcontext.NewTestContext(),
				followerID: "f1",
				followeeID: "f2",
			},
			output: output{err: nil},
			dependencies: func(in input, d *dependencies) {
				d.finder.On("ExistsByID", in.ctx, in.followerID).Return(true, nil)
				d.finder.On("ExistsByID", in.ctx, in.followeeID).Return(true, nil)
				d.finder.On("IsFollowing", in.ctx, in.followerID, in.followeeID).Return(false, nil)
				d.creator.On("FollowUser", in.ctx, in.followerID, in.followeeID).Return(nil)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected, actual)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &dependencies{
				creator: mocks.NewUserCreator(t),
				finder:  mocks.NewUserFinder(t),
				cache:   mocks.NewTimelineCache(t),
			}

			var done chan struct{}
			// Synchronize with the goroutine
			if tt.name == "should follow user successfully" {
				done = make(chan struct{})
				d.cache.On("InvalidateTimeline", twcontext.NewDetachedWithRequestID(tt.input.ctx), tt.input.followerID).Return(nil).Run(func(args mock.Arguments) {
					close(done)
				})
			}
			tt.dependencies(tt.input, d)

			uc := user.NewUserUseCase(d.creator, d.finder, d.cache)
			var actual output
			actual.err = uc.FollowUser(tt.input.ctx, tt.input.followerID, tt.input.followeeID)

			// Wait for the goroutine to finish
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

func Test_userUseCase_UnfollowUser(t *testing.T) {
	type input struct {
		ctx        context.Context
		followerID string
		followeeID string
	}

	type output struct {
		err error
	}

	type dependencies struct {
		creator *mocks.UserCreator
		finder  *mocks.UserFinder
		cache   *mocks.TimelineCache
	}

	tests := []struct {
		name         string
		input        input
		output       output
		dependencies func(in input, d *dependencies)
		assert       func(t *testing.T, expected, actual output)
	}{
		{
			name: "should return error if followerID or followeeID is empty",
			input: input{
				ctx:        twcontext.NewTestContext(),
				followerID: "",
				followeeID: "123",
			},
			output:       output{err: user.ErrInvalidInput},
			dependencies: func(in input, d *dependencies) {},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "should return error if followerID equals followeeID",
			input: input{
				ctx:        twcontext.NewTestContext(),
				followerID: "123",
				followeeID: "123",
			},
			output:       output{err: user.ErrCannotUnfollowSelf},
			dependencies: func(in input, d *dependencies) {},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "should return error if finder.ExistsByID for follower returns error",
			input: input{
				ctx:        twcontext.NewTestContext(),
				followerID: "f1",
				followeeID: "f2",
			},
			output: output{err: fmt.Errorf("failed to check follower with ID %s: %w", "f1", assert.AnError)},
			dependencies: func(in input, d *dependencies) {
				d.finder.On("ExistsByID", in.ctx, in.followerID).Return(false, assert.AnError)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected.err.Error(), actual.err.Error())
			},
		},
		{
			name: "should return error if follower does not exist",
			input: input{
				ctx:        twcontext.NewTestContext(),
				followerID: "f1",
				followeeID: "f2",
			},
			output: output{err: user.ErrUserNotFound},
			dependencies: func(in input, d *dependencies) {
				d.finder.On("ExistsByID", in.ctx, in.followerID).Return(false, nil)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "should return error if finder.ExistsByID for followee returns error",
			input: input{
				ctx:        twcontext.NewTestContext(),
				followerID: "f1",
				followeeID: "f2",
			},
			output: output{err: fmt.Errorf("failed to check followee with ID %s: %w", "f2", assert.AnError)},
			dependencies: func(in input, d *dependencies) {
				d.finder.On("ExistsByID", in.ctx, in.followerID).Return(true, nil)
				d.finder.On("ExistsByID", in.ctx, in.followeeID).Return(false, assert.AnError)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected.err.Error(), actual.err.Error())
			},
		},
		{
			name: "should return error if followee does not exist",
			input: input{
				ctx:        twcontext.NewTestContext(),
				followerID: "f1",
				followeeID: "f2",
			},
			output: output{err: user.ErrFolloweeNotFound},
			dependencies: func(in input, d *dependencies) {
				d.finder.On("ExistsByID", in.ctx, in.followerID).Return(true, nil)
				d.finder.On("ExistsByID", in.ctx, in.followeeID).Return(false, nil)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "should return error if finder.IsFollowing returns error",
			input: input{
				ctx:        twcontext.NewTestContext(),
				followerID: "f1",
				followeeID: "f2",
			},
			output: output{err: fmt.Errorf("error checking follow relationship: %w", assert.AnError)},
			dependencies: func(in input, d *dependencies) {
				d.finder.On("ExistsByID", in.ctx, in.followerID).Return(true, nil)
				d.finder.On("ExistsByID", in.ctx, in.followeeID).Return(true, nil)
				d.finder.On("IsFollowing", in.ctx, in.followerID, in.followeeID).Return(false, assert.AnError)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected.err.Error(), actual.err.Error())
			},
		},
		{
			name: "should return error if not following",
			input: input{
				ctx:        twcontext.NewTestContext(),
				followerID: "f1",
				followeeID: "f2",
			},
			output: output{err: user.ErrNotFollowing},
			dependencies: func(in input, d *dependencies) {
				d.finder.On("ExistsByID", in.ctx, in.followerID).Return(true, nil)
				d.finder.On("ExistsByID", in.ctx, in.followeeID).Return(true, nil)
				d.finder.On("IsFollowing", in.ctx, in.followerID, in.followeeID).Return(false, nil)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "should return error if creator.UnfollowUser returns error",
			input: input{
				ctx:        twcontext.NewTestContext(),
				followerID: "f1",
				followeeID: "f2",
			},
			output: output{err: fmt.Errorf("error unfollowing user: %w", assert.AnError)},
			dependencies: func(in input, d *dependencies) {
				d.finder.On("ExistsByID", in.ctx, in.followerID).Return(true, nil)
				d.finder.On("ExistsByID", in.ctx, in.followeeID).Return(true, nil)
				d.finder.On("IsFollowing", in.ctx, in.followerID, in.followeeID).Return(true, nil)
				d.creator.On("UnfollowUser", in.ctx, in.followerID, in.followeeID).Return(assert.AnError)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected.err.Error(), actual.err.Error())
			},
		},
		{
			name: "should unfollow user successfully",
			input: input{
				ctx:        twcontext.NewTestContext(),
				followerID: "f1",
				followeeID: "f2",
			},
			output: output{err: nil},
			dependencies: func(in input, d *dependencies) {
				d.finder.On("ExistsByID", in.ctx, in.followerID).Return(true, nil)
				d.finder.On("ExistsByID", in.ctx, in.followeeID).Return(true, nil)
				d.finder.On("IsFollowing", in.ctx, in.followerID, in.followeeID).Return(true, nil)
				d.creator.On("UnfollowUser", in.ctx, in.followerID, in.followeeID).Return(nil)
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected, actual)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &dependencies{
				creator: mocks.NewUserCreator(t),
				finder:  mocks.NewUserFinder(t),
				cache:   mocks.NewTimelineCache(t),
			}

			var done chan struct{}
			// Synchronize with the goroutine
			if tt.name == "should unfollow user successfully" {
				done = make(chan struct{})
				d.cache.On("InvalidateTimeline", twcontext.NewDetachedWithRequestID(tt.input.ctx), tt.input.followerID).Return(nil).Run(func(args mock.Arguments) {
					close(done)
				})
			}

			tt.dependencies(tt.input, d)

			uc := user.NewUserUseCase(d.creator, d.finder, d.cache)
			var actual output
			actual.err = uc.UnfollowUser(tt.input.ctx, tt.input.followerID, tt.input.followeeID)

			// Wait for the goroutine to finish
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
