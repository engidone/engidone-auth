package jwt

import (
	"context"
	"time"
)

func (rp sqlRepository) getUserByRefreshToken(refreshToken string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userID, err := rp.dbq.GetUserByRefreshToken(ctx, refreshToken)

	return userID.String(), err
}
