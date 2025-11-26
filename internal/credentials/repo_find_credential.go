package credentials

import (
	"context"
	"database/sql"
	"engidoneauth/internal/db"
	"engidoneauth/log"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
)

func (rp sqlRepository) findCredential(userID string, password string) (*db.Credential, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row, err := rp.dbq.GetCredential(ctx, db.GetCredentialParams{
		UserID:   uuid.MustParse(userID),
		Password: password,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, oops.
				With("user_id", userID).
				Wrap(InvalidCredentials)
		}
		log.Error("Error verifying credentials", log.String("user_id", userID), log.Err(err))
		return nil, oops.
			With("user_id", userID).
			Wrap(InternalError)
	}

	return &row, nil
}
