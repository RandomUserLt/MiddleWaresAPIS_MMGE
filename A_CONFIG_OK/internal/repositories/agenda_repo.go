package repositories

import (
	"context"
	"database/sql"
	"errors"

	"middleware/example/internal/models"
)

type AgendaRepository interface {
	Init(ctx context.Context) error
	List(ctx context.Context) ([]models.Agenda, error)
	Get(ctx context.Context, id string) (*models.Agenda, error)
	Create(ctx context.Context, a models.Agenda) error
	Update(ctx context.Context, a models.Agenda) error
	Delete(ctx context.Context, id string) error
}

type AgendaSQLite struct{ db *sql.DB }

func NewAgendaSQLite(db *sql.DB) *AgendaSQLite { return &AgendaSQLite{db: db} }

func (r *AgendaSQLite) Init(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS agendas(
			id   TEXT PRIMARY KEY,
			name TEXT NOT NULL
		);
	`)
	return err
}

func (r *AgendaSQLite) List(ctx context.Context) ([]models.Agenda, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name FROM agendas ORDER BY id`)
	if err != nil { return nil, err }
	defer rows.Close()

	var out []models.Agenda
	for rows.Next() {
		var a models.Agenda
		if err := rows.Scan(&a.ID, &a.Name); err != nil { return nil, err }
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *AgendaSQLite) Get(ctx context.Context, id string) (*models.Agenda, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, name FROM agendas WHERE id=?`, id)
	var a models.Agenda
	if err := row.Scan(&a.ID, &a.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) { return nil, nil }
		return nil, err
	}
	return &a, nil
}

func (r *AgendaSQLite) Create(ctx context.Context, a models.Agenda) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO agendas(id,name) VALUES(?,?)`, a.ID, a.Name)
	return err
}

func (r *AgendaSQLite) Update(ctx context.Context, a models.Agenda) error {
	res, err := r.db.ExecContext(ctx, `UPDATE agendas SET name=? WHERE id=?`, a.Name, a.ID)
	if err != nil { return err }
	if n, _ := res.RowsAffected(); n == 0 { return sql.ErrNoRows }
	return nil
}

func (r *AgendaSQLite) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM agendas WHERE id=?`, id)
	if err != nil { return err }
	if n, _ := res.RowsAffected(); n == 0 { return sql.ErrNoRows }
	return nil
}