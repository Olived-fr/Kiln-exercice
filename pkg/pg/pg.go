package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // register postgres driver
)

type Parameters struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	DBName   string `env:"POSTGRES_DB"`
	Username string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
}

type TxxBeginner interface {
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
}

func New(params Parameters) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable", // sslmode=disable is used for simplicity of dev.
		params.Username,
		params.Password,
		params.Host,
		params.Port,
		params.DBName,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return sqlx.NewDb(db.DB, "postgres"), nil
}

// Tx allows a user function to run multiple SQL queries inside a single transaction.
func Tx(ctx context.Context, tb TxxBeginner, fn func(*sqlx.Tx) error) error {
	if tx, err := tb.BeginTxx(ctx, nil); err != nil {
		return fmt.Errorf("begin txx: %w", err)
	} else if err = fn(tx); err != nil {
		return fmt.Errorf("fn: %w", errors.Join(err, tx.Rollback()))
	} else if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}
