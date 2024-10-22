package pg

import (
	"context"

	"github.com/jmoiron/sqlx"

	"kiln-exercice/internal/model"
)

type PollingRepository struct {
	db *sqlx.DB
}

func NewPollingRepository(db *sqlx.DB) *PollingRepository {
	return &PollingRepository{
		db: db,
	}
}

// UpsertPolling inserts or updates a polling record.
func (r *PollingRepository) UpsertPolling(ctx context.Context, polling model.Polling) error {
	const query = `
	INSERT INTO polling (id, last_polled_at) 
	VALUES ($1, $2) 
	ON CONFLICT (id)
	DO UPDATE SET last_polled_at = $2`

	_, err := r.db.ExecContext(ctx, query, polling.ID, polling.LastPolledAt)
	if err != nil {
		return err
	}

	return nil
}

// GetLastPolling retrieves the last polling record.
func (r *PollingRepository) GetLastPolling(ctx context.Context) (model.Polling, error) {
	const query = `SELECT id, last_polled_at FROM polling ORDER BY last_polled_at DESC LIMIT 1`

	var polling model.Polling
	err := r.db.GetContext(ctx, &polling, query)
	if err != nil {
		return model.Polling{}, err
	}

	return polling, nil
}
