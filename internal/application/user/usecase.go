package user

import (
	"context"
	"fmt"

	twcontext "github.com/oscarsalomon89/go-hexagonal/pkg/context"
)

type (
	UserFinder interface {
		ExistsByUsername(ctx context.Context, username string) (bool, error)
		ExistsByID(ctx context.Context, id string) (bool, error)
		IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error)
	}

	UserCreator interface {
		CreateUser(ctx context.Context, user *User) error
		FollowUser(ctx context.Context, followerID, followeeID string) error
		UnfollowUser(ctx context.Context, followerID, followeeID string) error
	}

	TimelineCache interface {
		InvalidateTimeline(ctx context.Context, userID string) error
	}

	userUseCase struct {
		creator UserCreator
		finder  UserFinder
		cache   TimelineCache
	}
)

func NewUserUseCase(creator UserCreator, finder UserFinder, cache TimelineCache) *userUseCase {
	return &userUseCase{creator: creator, finder: finder, cache: cache}
}

func (uc *userUseCase) CreateUser(ctx context.Context, user *User) error {
	if user.Username == "" {
		return ErrInvalidInput
	}

	exist, err := uc.finder.ExistsByUsername(ctx, user.Username)
	if err != nil {
		return fmt.Errorf("failed to check username: %w", err)
	}

	if exist {
		return ErrUsernameExists
	}

	if err := uc.creator.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (uc *userUseCase) FollowUser(ctx context.Context, followerID, followeeID string) error {
	if followerID == "" || followeeID == "" {
		return ErrInvalidInput
	}

	if followerID == followeeID {
		return ErrCannotFollowSelf
	}

	if followerExists, err := uc.finder.ExistsByID(ctx, followerID); err != nil {
		return fmt.Errorf("failed to check follower with ID %s: %w", followerID, err)
	} else if !followerExists {
		return ErrUserNotFound
	}

	if followeeExists, err := uc.finder.ExistsByID(ctx, followeeID); err != nil {
		return fmt.Errorf("failed to check followee with ID %s: %w", followeeID, err)
	} else if !followeeExists {
		return ErrFolloweeNotFound
	}

	exists, err := uc.finder.IsFollowing(ctx, followerID, followeeID)
	if err != nil {
		return fmt.Errorf("error checking follow relationship: %w", err)
	}
	if exists {
		return ErrAlreadyFollowing
	}

	if err := uc.creator.FollowUser(ctx, followerID, followeeID); err != nil {
		return fmt.Errorf("error following user: %w", err)
	}

	go uc.invalidateTimelineAsync(twcontext.NewDetachedWithRequestID(ctx), followerID)

	return nil
}

func (uc *userUseCase) UnfollowUser(ctx context.Context, followerID, followeeID string) error {
	if followerID == "" || followeeID == "" {
		return ErrInvalidInput
	}

	if followerID == followeeID {
		return ErrCannotUnfollowSelf
	}

	if followerExists, err := uc.finder.ExistsByID(ctx, followerID); err != nil {
		return fmt.Errorf("failed to check follower with ID %s: %w", followerID, err)
	} else if !followerExists {
		return ErrUserNotFound
	}

	if followeeExists, err := uc.finder.ExistsByID(ctx, followeeID); err != nil {
		return fmt.Errorf("failed to check followee with ID %s: %w", followeeID, err)
	} else if !followeeExists {
		return ErrFolloweeNotFound
	}

	exists, err := uc.finder.IsFollowing(ctx, followerID, followeeID)
	if err != nil {
		return fmt.Errorf("error checking follow relationship: %w", err)
	}
	if !exists {
		return ErrNotFollowing
	}

	if err := uc.creator.UnfollowUser(ctx, followerID, followeeID); err != nil {
		return fmt.Errorf("error unfollowing user: %w", err)
	}

	go uc.invalidateTimelineAsync(twcontext.NewDetachedWithRequestID(ctx), followerID)

	return nil
}

func (uc *userUseCase) invalidateTimelineAsync(ctx context.Context, userID string) {
	logger := twcontext.Logger(ctx)

	if err := uc.cache.InvalidateTimeline(ctx, userID); err != nil {
		logger.WithError(err).WithField("user_id", userID).Error("failed to invalidate timeline")
	}
}
