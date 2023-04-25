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

// Executes the function "toExec" inside a transaction ("toExec" does not return any value)
// The transaction is automatically rolled back in case "toExec" returns an error
// opts can be used to pass arguments to the transaction
func ExecWithTransactionNoResponse(
	db *gorm.DB,
	toExec func(*gorm.DB) error,
	opts ...*sql.TxOptions,
) error {
	tx := db.Begin(opts...)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := toExec(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
