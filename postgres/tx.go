package postgres

import (
	"context"
	"database/sql"
	
	"github.com/jmoiron/sqlx"
)

// NewTx creates a new database transaction with default isolation levels and
// passes a context to roll back a transaction in the case of a
// context cancellation. Should clean up any zombie connections that could be
// if we did not control for this.
func NewTx(ctx context.Context, db *sqlx.DB) (*sqlx.Tx, error) {
	return db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
}

// Commit will attept to commit the provided resource. Upon a failure to commit, the
// transaction will be rolled back.
func Commit(tx *sqlx.Tx, resource string) error {
	if err := tx.Commit(); err != nil {
		return Rollback(tx, resource, Err(resource, err))
	}

	return nil
}

// Rollback undoes the supplied transaction.
func Rollback(tx *sqlx.Tx, resource string, err error) error {
	if commitErr := tx.Rollback(); commitErr != nil {
		return Err(resource, commitErr)
	}
	return err
}
