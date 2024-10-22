package pgtest

import (
	"context"
	"fmt"
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // register postgres driver
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	dbUser   = "root"
	dbPass   = "pass"
	dbName   = "test_db"
	initPath = "./scripts/init.sql"
)

// PostgresContainer wraps the logic of starting a Postgres testcontainer (i.e. docker)
// instance, initializing the database schema and establishing a working connection.
type PostgresContainer struct {
	container testcontainers.Container
	db        *sqlx.DB
}

type Option func(*PostgresContainer) error

func NewPostgresContainer(ctx context.Context, t *testing.T) *PostgresContainer {
	t.Helper()

	c, err := createPostgresContainer(ctx)
	require.NoError(t, err)

	t.Cleanup(c.terminate)

	return c
}

func createPostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	var (
		pc  PostgresContainer
		err error
	)

	path, err := filepath.Abs(initPath)

	pc.container, err = postgres.Run(
		ctx,
		"postgres:17-alpine",
		postgres.WithInitScripts(path),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(15*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("start container: %w", err)
	}

	host, err := pc.container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("get container host: %w", err)
	}

	port, err := pc.container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return nil, fmt.Errorf("get container mapped port: %w", err)
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		dbUser,
		dbPass,
		net.JoinHostPort(host, port.Port()),
		dbName,
	)

	if pc.db, err = sqlx.Connect("postgres", dsn); err != nil {
		return nil, fmt.Errorf("connect and ping db: %w", err)
	}

	return &pc, nil
}

// GetDB returns a connection to the database wrapped within the testcontainer instance.
func (c *PostgresContainer) GetDB() *sqlx.DB {
	return c.db
}

// terminate shutdowns the testcontainer instance.
func (c *PostgresContainer) terminate() {
	_ = c.container.Terminate(context.Background())
}
