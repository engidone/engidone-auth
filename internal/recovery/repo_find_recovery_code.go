package recovery

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"engidoneauth/internal/apperror"
)



func (rp *sqlRepository) findRecoveryCode(code string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	recoveryCode, err := rp.dbq.GetRecoveryCode(ctx, code)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return "", apperror.New(ErrInvalidRecoveryCode, "Recovery code not found")
		}
		return "", apperror.New(ErrInternalError, fmt.Sprintf("Error finding recovery: %v", err))
	}

	return recoveryCode, nil

}
