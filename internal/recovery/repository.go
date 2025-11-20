package recovery

import "engidoneauth/internal/db"

type sqlRepository struct {
	dbq *db.Queries
}

func NewSQLRepository(dbq *db.Queries) *sqlRepository {
	return &sqlRepository{
		dbq: dbq,
	}
}