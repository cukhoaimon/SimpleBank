package usecase

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *SQLStore {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted, // default
	})

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rollbackError := tx.Rollback(); rollbackError != nil {
			return fmt.Errorf("tx error: %v, rollback err: %v", err, rollbackError)
		}

		return err
	}

	return tx.Commit()
}
