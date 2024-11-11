package safesql

import (
	"context"
	"database/sql"

	"github.com/andriusm/def-prog-exercises/safesql/internal/raw"
)

// has to be private
type compileTimeConstant string

// has to be public
type TrustedSQL struct {
	s string
}

func New(text compileTimeConstant) TrustedSQL {
	return TrustedSQL{string(text)}
}

type DB struct {
	db *sql.DB
}

func (db *DB) Query(ctx context.Context, query TrustedSQL, args ...any) (*sql.Rows, error) {
	r, err := db.db.QueryContext(ctx, query.s, args...)
	return r, err
}

type Rows = sql.Rows
type Result = sql.Result

// ---
// wrap sql.db.ExecContext and sql.Open

func (db *DB) Close() error {
	return db.db.Close()
}

func Open(driverName, dataSourceName string) (*DB, error) {
	d, err := sql.Open(driverName, dataSourceName)
	return &DB{d}, err
}

func (db *DB) QueryContext(ctx context.Context,
	query TrustedSQL, args ...any) (*Rows, error) {
	return db.db.QueryContext(ctx, query.s, args...)
}
func (db *DB) ExecContext(ctx context.Context,
	query TrustedSQL, args ...any) (Result, error) {
	return db.db.ExecContext(ctx, query.s, args...)
}

// ---

// Guaranteed to run before any package that imports safesql
func init() {
	raw.TrustedSQLCtor =
		func(unsafe string) TrustedSQL {
			return TrustedSQL{unsafe}
		}
}
