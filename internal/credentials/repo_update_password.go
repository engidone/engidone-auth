package credentials

import (
	"context"
	"database/sql"
	"engidoneauth/internal/db"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
)

func (rp sqlRepository) updatePassword(userID string, newPassword string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := rp.dbq.UpdatePassword(ctx,
		db.UpdatePasswordParams{
			Password: newPassword,
			UserID:   uuid.MustParse(userID),
		},
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, oops.
				With("user_id", userID).
				Wrap(UserNotFound)
		}
		return false, oops.
			With("user_id", userID).
			Wrap(InternalError)
	}

	return true, nil
}
