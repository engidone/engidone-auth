package credentials

import (
	"context"
	"database/sql"
	"engidoneauth/internal/apperror"
	"engidoneauth/internal/db"
	"engidoneauth/log"
	"fmt"
	"time"

	"github.com/google/uuid"
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
			return nil, apperror.New(ErrInvalidCredentials, "Invalid credentials")
		}
		log.Error("Error verifying credentials", log.String("user_id", userID), log.Err(err))
		return nil, apperror.New(ErrInternalError, fmt.Sprintf("Error verifying credentials: %v", err))
	}

	return &row, err

}
