package pgtest

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

// Record represents a row in a table. It is just a standard go map, where each
// key/value pair corresponds to a column/value pair in the actual db row.
type Record map[string]any

// RecordSet groups records belonging to a given table.
type RecordSet struct {
	Table     string
	Records   []Record
	IsDeleted bool
}

// Helper provides methods to easily inject data in a test database and
// to check the state of that database once a test has ran.
type Helper struct {
	db *sqlx.DB
}

func NewHelper(db *sqlx.DB) *Helper {
	return &Helper{
		db: db,
	}
}

// inject inserts sets of records in the database.
func (h *Helper) inject(ctx context.Context, rss []RecordSet) error {
	for _, rs := range rss {
		if len(rs.Records) == 0 {
			return fmt.Errorf("expected records for table %s", rs.Table)
		}

		columns := make(map[string]struct{})

		for _, rec := range rs.Records {
			for column := range rec {
				columns[column] = struct{}{}
			}
		}

		var uniqColumns []string
		for column := range columns {
			uniqColumns = append(uniqColumns, column)
		}

		builder := squirrel.
			Insert(rs.Table).
			Columns(uniqColumns...).
			PlaceholderFormat(squirrel.Dollar)

		for _, rec := range rs.Records {
			var values []any
			for _, column := range uniqColumns {
				values = append(values, rec[column])
			}

			builder = builder.Values(values...)
		}

		query, args := builder.MustSql()

		_, err := h.db.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("insert records in table %s: %w", rs.Table, err)
		}
	}

	return nil
}

func (h *Helper) MustInject(ctx context.Context, t *testing.T, rss []RecordSet) {
	t.Helper()

	err := h.inject(ctx, rss)
	require.NoError(t, err)
}

// check verifies that a list of record sets are present in the database.
func (h *Helper) check(ctx context.Context, rss []RecordSet) error {
	for _, rs := range rss {
		for _, rec := range rs.Records {
			builder := squirrel.
				Select("1").
				From(rs.Table).
				PlaceholderFormat(squirrel.Dollar)

			for col, val := range rec {
				builder = builder.Where(squirrel.Eq{col: val})
			}

			query, args := builder.MustSql()

			var n int

			err := h.db.QueryRowContext(ctx, query, args...).Scan(&n)
			if err != nil {
				msg := fmt.Sprintf("select row from table %s with filter %+v", rs.Table, rec)

				if errors.Is(err, sql.ErrNoRows) {
					if rs.IsDeleted {
						continue
					}

					return fmt.Errorf("%s: no corresponding row found", msg)
				}

				return fmt.Errorf("%s: %w", msg, err)
			}

			if rs.IsDeleted {
				msg := fmt.Sprintf("select row from table %s with filter %+v", rs.Table, rec)

				return fmt.Errorf("%s: row found but should have been deleted", msg)
			}
		}
	}

	return nil
}

func (h *Helper) MustCheck(ctx context.Context, t *testing.T, rss []RecordSet) {
	t.Helper()

	err := h.check(ctx, rss)
	require.NoError(t, err, "pgtest")
}
