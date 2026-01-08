package recovery

import (
	"context"
	"database/sql"
	"time"

	"github.com/samber/oops"
)

// findRecoveryCode remains for backward compatibility but delegates to validateRecoveryCode
func (rp *sqlRepository) findRecoveryCode(code string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	storedCode, err := rp.dbq.GetRecoveryCode(ctx, code)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", RecoveryCodeNotFound
		}

		return "", oops.With("error", err.Error()).Wrap(InternalError)
	}

	return storedCode, nil
}
