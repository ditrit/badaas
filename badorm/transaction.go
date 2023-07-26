package badorm

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"

	"github.com/ditrit/badaas/badorm/logger"
)

// type TransactionExecutor interface {
// 	Exec[RT any](
// 		db *gorm.DB,
// 		toExec func(*gorm.DB) (RT, error),
// 		opts ...*sql.TxOptions,
// 	) (RT, error)
// }
// quizas que retorne any y el metodo haga el casteo y que este sea una variable global (no tengo claro como seria la concurrencia en ese caso)

// Executes the function "toExec" inside a transaction
// The transaction is automatically rolled back in case "toExec" returns an error
// opts can be used to pass arguments to the transaction
func Transaction[RT any](
	logger logger.Interface,
	db *gorm.DB,
	toExec func(*gorm.DB) (RT, error),
	opts ...*sql.TxOptions,
) (RT, error) {
	begin := time.Now()

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

	logger.TraceTransaction(context.Background(), begin)

	return returnValue, nil
}

// TODO transaction no return
// TODO quizas podria ser un objeto en lugar de una funcion
// TODO seria bueno que fuera configurable si queres esta medicion de tiempo o no
// TODO warnings para transacciones lentas
