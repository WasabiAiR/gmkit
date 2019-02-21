package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// CRUD is a ready to go type that implements most of the basic methods
// we use for 90%+ of our database calls.
type CRUD interface {
	Getter
	Selecter
	NameBinder
	Rebinder
}

// Execer provides the exec behavior.
type Execer interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// ExecRebinder is an interface that is agnostic for database transactions for the
// sql type, an execer can be a sqlx.DB, transactions for something custom.
type ExecRebinder interface {
	Execer
	Rebinder
}

// NamedExecer preforms an operating that returns an sql.Result and error.
type NamedExecer interface {
	NamedExecContext(ctx context.Context, query string, args interface{}) (sql.Result, error)
}

// RowQueryBinder preforms an operating that returns a pointer sql.Result.
type RowQueryBinder interface {
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	Rebinder
}

// QueryBinder preforms an operating that returns a pointer sql.Result.
type QueryBinder interface {
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	Rebinder
}

// NamedRowQuerier allows you to use the named query arguments with a row query.
type NamedRowQuerier interface {
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	NameBinder
}

// NamedQuerier allows you to use the named query arugments with a query.
type NamedQuerier interface {
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	NameBinder
}

// NamedGetBinder allows you to use the GetContext query.
type NamedGetBinder interface {
	Getter
	Rebinder
	NameBinder
}

// Getter provides get functionality.
type Getter interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

// GetRebinder provides the get and rebinding functionality.
type GetRebinder interface {
	Getter
	Rebinder
}

// NameBinder preforms an operating that returns a pointer sql.Result.
type NameBinder interface {
	BindNamed(query string, v interface{}) (bindedQuery string, args []interface{}, err error)
}

// Rebinder preforms a strings altering operation.
type Rebinder interface {
	Rebind(string) string
}

// Selecter performs select behavior with contexts.
type Selecter interface {
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

// SelectBinder preforms a query with a context.
type SelectBinder interface {
	Selecter
	Rebinder
}
