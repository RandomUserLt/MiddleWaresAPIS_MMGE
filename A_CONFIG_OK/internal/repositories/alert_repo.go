package repositories

import (
	"context"
	"database/sql"
	"errors"

	"middleware/example/internal/models"
)

type AlertRepository interface {
	Init(ctx context.Context) error
	List(ctx context.Context, agendaID *string) ([]models.Alert, error)
	Get(ctx context.Context, id string) (*models.Alert, error)
	Create(ctx context.Context, a models.Alert) error
	Update(ctx context.Context, a models.Alert) error
	Delete(ctx context.Context, id string) error
}

type AlertSQLite struct{ db *sql.DB }

func NewAlertSQLite(db *sql.DB) *AlertSQLite { return &AlertSQLite{db: db} }

func (r *AlertSQLite) Init(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS alerts (
            id         TEXT PRIMARY KEY,
            agenda_id  TEXT NOT NULL,
            target     TEXT NOT NULL,
            condition  TEXT NOT NULL DEFAULT 'always',
            created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
            FOREIGN KEY(agenda_id) REFERENCES agendas(id) ON DELETE CASCADE
        );
        CREATE INDEX IF NOT EXISTS idx_alerts_agenda ON alerts(agenda_id);
    `)
	return err
}

func (r *AlertSQLite) List(ctx context.Context, agendaID *string) ([]models.Alert, error) {
	if agendaID != nil && *agendaID != "" {
		rows, err := r.db.QueryContext(ctx, `SELECT id, agenda_id, target, condition FROM alerts WHERE agenda_id=? ORDER BY id`, *agendaID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var out []models.Alert
		for rows.Next() {
			var a models.Alert
			if err := rows.Scan(&a.ID, &a.AgendaID, &a.Target, &a.Condition); err != nil {
				return nil, err
			}
			out = append(out, a)
		}
		return out, rows.Err()
	}

	rows, err := r.db.QueryContext(ctx, `SELECT id, agenda_id, target, condition FROM alerts ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Alert
	for rows.Next() {
		var a models.Alert
		if err := rows.Scan(&a.ID, &a.AgendaID, &a.Target, &a.Condition); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *AlertSQLite) Get(ctx context.Context, id string) (*models.Alert, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, agenda_id, target, condition FROM alerts WHERE id=?`, id)
	var a models.Alert
	if err := row.Scan(&a.ID, &a.AgendaID, &a.Target, &a.Condition); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *AlertSQLite) Create(ctx context.Context, a models.Alert) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO alerts(id, agenda_id, target, condition) VALUES(?,?,?,?)`, a.ID, a.AgendaID, a.Target, a.Condition)
	return err
}

func (r *AlertSQLite) Update(ctx context.Context, a models.Alert) error {
	res, err := r.db.ExecContext(ctx, `UPDATE alerts SET agenda_id=?, target=?, condition=? WHERE id=?`, a.AgendaID, a.Target, a.Condition, a.ID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *AlertSQLite) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM alerts WHERE id=?`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
