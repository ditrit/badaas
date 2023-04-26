package services

import (
	"database/sql"

	"gorm.io/gorm"
)

// Executes the function "toExec" inside a transaction
// The transaction is automatically rolled back in case "toExec" returns an error
// opts can be used to pass arguments to the transaction
func ExecWithTransaction[RT any](
	db *gorm.DB,
	toExec func(*gorm.DB) (RT, error),
	opts ...*sql.TxOptions,
) (RT, error) {
	nilValue := *new(RT)
	tx := db.Begin(opts...)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nilValue, err
	}

	returnValue, err := toExec(tx)
	if err != nil {
		tx.Rollback()
		return nilValue, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nilValue, err
	}

	return returnValue, nil
}
