package tweet

import (
	"context"

	twcontext "github.com/oscarsalomon89/go-hexagonal/pkg/context"
)

// invalidateFollowersTimelinesAsync invalidates the timeline cache for all followers of a user asynchronously.
//
// TODO: For better scalability and decoupling, consider using a queue system (e.g., AWS SQS) to handle timeline invalidations asynchronously instead of spawning goroutines.
func (uc *usecase) invalidateFollowersTimelinesAsync(ctx context.Context, userID string) {
	logger := twcontext.Logger(ctx)

	followers, err := uc.userFinder.GetFollowers(ctx, userID)
	if err != nil {
		logger.WithError(err).Error("Failed to get followers of user")
		return
	}
	for _, followerID := range followers {
		fID := followerID
		go func() {
			if err := uc.cache.InvalidateTimeline(ctx, fID); err != nil {
				logger.WithError(err).WithField("follower_id", fID).Error("failed to invalidate timeline")
			}
		}()
	}
}
