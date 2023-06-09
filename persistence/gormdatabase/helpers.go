package gormdatabase

import (
	"errors"

	"github.com/jackc/pgconn"
)

func IsDuplicateKeyError(err error) bool {
	// unique_violation code is equals to 23505
	return IsPostgresError(err, "23505")
}

func IsPostgresError(err error, errCode string) bool {
	var pgerr *pgconn.PgError
	if ok := errors.As(err, &pgerr); ok {
		return pgerr.Code == errCode
	}

	return false
}
