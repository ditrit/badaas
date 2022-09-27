package gormdatabase

import "github.com/jackc/pgconn"

func IsDuplicateKeyError(err error) bool {
	// unique_violation code is equals to 23505
	return isPostgresError(err, "23505")
}

func isPostgresError(err error, errCode string) bool {
	pgErr, ok := err.(*pgconn.PgError)
	if ok {
		return pgErr.Code == errCode

	}
	return false
}
