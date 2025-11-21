package credentials

import (
	"context"
	"database/sql"
	"engidoneauth/internal/apperror"
	"engidoneauth/internal/db"
	"fmt"
	"time"

	"github.com/google/uuid"
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
			return false, apperror.New(ErrUserNotFound, "Usuario no encontrado")
		}
		return false, apperror.New(ErrInternalError, fmt.Sprintf("Error actualizando usuario: %v", err))
	}

	return true, nil
}
