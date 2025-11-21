package jwt

import (
	"context"
	"engidoneauth/internal/db"
	"time"

	"github.com/google/uuid"
)

func (rp sqlRepository) existsRefreshToken(userID, refreshToken string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	exists, err := rp.dbq.ExistsRefreshToken(
		ctx, db.ExistsRefreshTokenParams{
			UserID:       uuid.MustParse(userID),
			RefreshToken: refreshToken,
		},
	)

	return exists, err
}
