package service

import (
	"leads/internal/domain"
	"time"

	"github.com/google/uuid"
)

type LeadService struct {
	leadRepo domain.LeadRepository
}

func NewLeadService(repo domain.LeadRepository) domain.LeadService {
	return &LeadService{leadRepo: repo}
}

func (s *LeadService) CreateLead(lead *domain.Lead) error {
	lead.ID = uuid.New().String()
	lead.Status = "New"
	lead.CreatedAt = time.Now()
	lead.UpdatedAt = time.Now()
	return s.leadRepo.CreateLead(lead)
}

func (s *LeadService) GetLead(id string) (*domain.Lead, error) {
	return s.leadRepo.GetLead(id)
}

func (s *LeadService) UpdateLeadStatus(id string, status string) error {
	return s.leadRepo.UpdateLeadStatus(id, status)
}

func (s *LeadService) ListLeads(status string, source string, limit int) ([]domain.Lead, error) {
	return s.leadRepo.ListLeads(status, source, limit)
}
