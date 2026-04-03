package domain

import "time"

type Lead struct {
	ID        string
	Name      string
	Phone     string
	Email     string
	Source    string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type LeadRepository interface {
	CreateLead(lead *Lead) error
	GetLead(id string) (*Lead, error)
	UpdateLeadStatus(id string, status string) error
	ListLeads(status string, source string, limit int) ([]Lead, error)
}

type LeadService interface {
	CreateLead(lead *Lead) error
	GetLead(id string) (*Lead, error)
	UpdateLeadStatus(id string, status string) error
	ListLeads(status string, source string, limit int) ([]Lead, error)
}
