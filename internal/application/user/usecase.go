package user

import (
	"context"
	"fmt"
)

type (
	UsersFinder interface {
		ExistsByUsername(ctx context.Context, username string) (bool, error)
		ExistsByID(ctx context.Context, id string) (bool, error)
		IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error)
	}

	UsersCreator interface {
		CreateUser(ctx context.Context, user *User) error
		FollowUser(ctx context.Context, followerID, followeeID string) error
	}

	userUseCase struct {
		creator UsersCreator
		finder  UsersFinder
	}
)

func NewUserUseCase(creator UsersCreator, finder UsersFinder) *userUseCase {
	return &userUseCase{creator: creator, finder: finder}
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

	return nil
}
