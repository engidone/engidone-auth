package jwt

import (
	"context"
	"engidoneauth/internal/db"
	"engidoneauth/log"
	"time"
	"github.com/google/uuid"
	"github.com/samber/oops"
)

func (rp sqlRepository) syncRefreshToken(userID, refreshToken string) (*db.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	item, err := rp.dbq.InsertOrUpdateRefreshToken(
		ctx, db.InsertOrUpdateRefreshTokenParams{
			UserID:       uuid.MustParse(userID),
			RefreshToken: refreshToken,
		},
	)

	if err != nil {
		log.Error("Error syncing refresh token: ", log.String("user_id", userID), log.Err(err))
		return nil, oops.
			With("user_id", userID).
			With("refresh_token", refreshToken).
			Wrap(InsertOrUpdateRefreshToken)
	}

	return &item, nil
}
