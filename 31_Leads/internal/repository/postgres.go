package repository

import (
	"database/sql"
	"leads/internal/domain"
	"time"
)

type PostgresLeadRepo struct {
	db *sql.DB
}

func NewPostgresLeadRepo(db *sql.DB) domain.LeadRepository {
	return &PostgresLeadRepo{db: db}
}

func (r *PostgresLeadRepo) CreateLead(lead *domain.Lead) error {
	query := `INSERT INTO leads (id, name, phone, email, source, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(query, lead.ID, lead.Name, lead.Phone, lead.Email, lead.Source, lead.Status, lead.CreatedAt, lead.UpdatedAt)
	return err
}
func (r *PostgresLeadRepo) GetLead(id string) (*domain.Lead, error) {
	query := "SELECT id, name, phone, email, source, status, created_at, updated_at FROM leads WHERE id = $1"
	row := r.db.QueryRow(query, id)
	var lead domain.Lead
	err := row.Scan(&lead.ID, &lead.Name, &lead.Phone, &lead.Email, &lead.Source, &lead.Status, &lead.CreatedAt, &lead.UpdatedAt)
	return &lead, err
}

func (r *PostgresLeadRepo) UpdateLeadStatus(id string, status string) error {
	query := "UPDATE leads SET status = $1, updated_at = $2 WHERE id = $3"
	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}

func (r *PostgresLeadRepo) ListLeads(status string, source string, limit int) ([]domain.Lead, error) {
	query := "SELECT id, name, phone, email, source, status, created_at, updated_at FROM leads"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var leads []domain.Lead
	for rows.Next() {
		var lead domain.Lead
		err := rows.Scan(&lead.ID, &lead.Name, &lead.Phone, &lead.Email, &lead.Source, &lead.Status, &lead.CreatedAt, &lead.UpdatedAt)
		if err != nil {
			return nil, err
		}
		leads = append(leads, lead)
	}
	return leads, nil
}
