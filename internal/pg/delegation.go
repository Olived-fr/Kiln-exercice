package pg

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"kiln-exercice/internal/model"
	"kiln-exercice/pkg/pg"
)

type DelegationRepository struct {
	db        *sqlx.DB
	batchSize int
}

func NewDelegationRepository(db *sqlx.DB, batchSize int) *DelegationRepository {
	return &DelegationRepository{
		db:        db,
		batchSize: batchSize,
	}
}

// InsertDelegations bulk inserts a list of delegations into the database.
func (r *DelegationRepository) InsertDelegations(ctx context.Context, delegations []model.Delegation) error {
	const query = `
	INSERT INTO delegation (datetime, amount, delegator, height, tx_hash)
	VALUES (:datetime, :amount, :delegator, :height, :tx_hash)
	ON CONFLICT (tx_hash) DO NOTHING
	`

	return pg.Tx(
		ctx, r.db, func(tx *sqlx.Tx) error {
			for i := 0; i < len(delegations); i += r.batchSize {
				end := i + r.batchSize
				if end > len(delegations) {
					end = len(delegations)
				}

				batch := delegations[i:end]
				_, err := tx.NamedExecContext(ctx, query, batch)
				if err != nil {
					return fmt.Errorf("batch insert: %w", err)
				}
			}
			return nil
		},
	)
}

// ListDelegations returns a list of delegations paginated.
func (r *DelegationRepository) ListDelegations(ctx context.Context, year, offset, limit int) ([]model.Delegation, error) {
	var (
		whereClauses []string
		queryArgs    []any
	)

	query := `SELECT datetime, amount, delegator, height FROM delegation`

	if year != 0 {
		whereClauses = append(whereClauses, "EXTRACT(YEAR FROM datetime) = $1")
		queryArgs = append(queryArgs, year)
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	query = pg.NewPagination(limit, offset, "datetime").Embed(query)

	var delegations []model.Delegation
	err := r.db.SelectContext(ctx, &delegations, query, queryArgs...)
	if err != nil {
		return nil, err
	}

	return delegations, nil
}
