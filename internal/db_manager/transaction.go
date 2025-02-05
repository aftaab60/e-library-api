package db_manager

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type key string

const (
	transactionKey key = "transaction_key"
)

// WrapInTransaction wraps operations in a database transaction.
func WrapInTransaction(ctx context.Context, db ItxDB, f func(ctx context.Context) error, onRollback func(error)) (err error) {
	if db != nil {
		tx := GetTransactionFromContext(ctx)
		if tx == nil {
			tx, err = db.Begin()
			if err != nil {
				return err
			}
			ctx = context.WithValue(ctx, transactionKey, tx)

			defer func() {
				if r := recover(); r != nil {
					RollbackTransaction(fmt.Errorf("panic error: %v", r), tx, onRollback)
					panic(r)
				}
				if err != nil {
					RollbackTransaction(err, tx, onRollback)
				} else {
					if commitErr := tx.Commit(); commitErr != nil {
						RollbackTransaction(commitErr, tx, onRollback)
						err = commitErr
					}
				}
			}()
		}
		err = f(ctx)
		return err
	} else {
		//for in-memory, it's a regular execution.
		//But we mimic transaction into map and manage it through a common transactionManager
		if err := f(ctx); err != nil {
			if onRollback != nil {
				onRollback(err)
			}
			return err
		}
		return nil
	}
}

// GetTransactionFromContext retrieves the transaction from context.
func GetTransactionFromContext(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(transactionKey).(*sql.Tx); ok {
		return tx
	}
	return nil
}

// RollbackTransaction rolls back the transaction with optional rollback handling.
func RollbackTransaction(err error, tx *sql.Tx, onRollback func(error)) {
	if rollbackErr := tx.Rollback(); rollbackErr != nil {
		log.Printf("Rollback failed: %v", rollbackErr)
	}
	if onRollback != nil {
		onRollback(err)
	}
}
