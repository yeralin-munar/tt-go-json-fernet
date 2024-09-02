package postgres

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type TransactionManager struct {
	log *log.Helper

	db *DB
}

func NewTransactionManager(db *DB, logger log.Logger) *TransactionManager {
	return &TransactionManager{
		log: log.NewHelper(logger),
		db:  db,
	}
}

func (t *TransactionManager) InTransaction(ctx context.Context, txFunc func(ctx context.Context) error) error {
	tx, err := t.db.GetConnection(ctx).Begin(ctx)
	if err != nil {
		return err
	}

	err = txFunc(context.WithValue(ctx, dbKey, tx))
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil {
			t.log.Errorf("Rollback transaction failed: %v", rollbackErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
