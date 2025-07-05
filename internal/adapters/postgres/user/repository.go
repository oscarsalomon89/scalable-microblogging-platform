package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/oscarsalomon89/go-hexagonal/internal/application/user"
	db "github.com/oscarsalomon89/go-hexagonal/internal/platform/pg"
)

type userRepository struct {
	db db.Connections
}

func NewUserRepository(db db.Connections) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *user.User) error {
	userModel := fromDomain(user)

	err := r.db.MasterConn.
		WithContext(ctx).
		Create(userModel).Error
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	user.ID = userModel.ID.String()
	user.CreatedAt = userModel.CreatedAt
	user.UpdatedAt = userModel.UpdatedAt

	return nil
}

func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := r.db.MasterConn.
		WithContext(ctx).
		Model(&User{}).
		Where("username = ?", username).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to find user: %w", err)
	}

	return count > 0, nil
}

func (r *userRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	var count int64
	if err := r.db.MasterConn.
		WithContext(ctx).
		Model(&User{}).
		Where("id = ?", id).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to find user: %w", err)
	}

	return count > 0, nil
}

func (r *userRepository) IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error) {
	var count int64
	if err := r.db.MasterConn.
		WithContext(ctx).
		Model(&Follow{}).
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to find follow: %w", err)
	}

	return count > 0, nil
}

func (r *userRepository) FollowUser(ctx context.Context, followerID, followeeID string) error {
	if followerID == "" || followeeID == "" {
		return fmt.Errorf("followerID or followeeID is empty")
	}

	followerUUID, err := uuid.Parse(followerID)
	if err != nil {
		return fmt.Errorf("invalid followerID: %w", err)
	}

	followeeUUID, err := uuid.Parse(followeeID)
	if err != nil {
		return fmt.Errorf("invalid followeeID: %w", err)
	}

	if err := r.db.MasterConn.
		WithContext(ctx).
		Create(&Follow{
			FollowerID: followerUUID,
			FolloweeID: followeeUUID,
		}).Error; err != nil {
		return fmt.Errorf("error creating follow relationship: %w", err)
	}

	return nil
}

func (r *userRepository) GetFollowers(ctx context.Context, id string) ([]string, error) {
	var followers []string
	if err := r.db.MasterConn.
		WithContext(ctx).
		Model(&Follow{}).
		Where("followee_id = ?", id).
		Pluck("follower_id", &followers).Error; err != nil {
		return nil, fmt.Errorf("failed to find followers: %w", err)
	}

	return followers, nil
}

func (r *userRepository) GetFollowees(ctx context.Context, id string) ([]string, error) {
	var followees []string
	if err := r.db.MasterConn.
		WithContext(ctx).
		Model(&Follow{}).
		Where("follower_id = ?", id).
		Pluck("followee_id", &followees).Error; err != nil {
		return nil, fmt.Errorf("failed to find followees: %w", err)
	}

	return followees, nil
}
