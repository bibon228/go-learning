package grps

import (
	"context"
	"leads/internal/domain"
	"leads/pb"
)

type LeadHandler struct {
	pb.UnimplementedLeadServiceServer
	leadService domain.LeadService
}

func NewLeadHandler(leadService domain.LeadService) *LeadHandler {
	return &LeadHandler{leadService: leadService}
}

func (h *LeadHandler) CreateLead(ctx context.Context, req *pb.CreateLeadRequest) (*pb.CreateLeadResponse, error) {
	lead := &domain.Lead{
		Name:   req.Name,
		Phone:  req.Phone,
		Email:  req.Email,
		Source: req.Source,
	}
	err := h.leadService.CreateLead(lead)
	if err != nil {
		return nil, err
	}
	return &pb.CreateLeadResponse{Id: lead.ID}, nil
}

func (h *LeadHandler) GetLead(ctx context.Context, req *pb.GetLeadRequest) (*pb.GetLeadResponse, error) {
	lead, err := h.leadService.GetLead(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetLeadResponse{Lead: &pb.Lead{
		Id:     lead.ID,
		Name:   lead.Name,
		Phone:  lead.Phone,
		Email:  lead.Email,
		Source: lead.Source,
		Status: lead.Status,
	}}, nil
}

func (h *LeadHandler) UpdateLeadStatus(ctx context.Context, req *pb.UpdateLeadStatusRequest) (*pb.UpdateLeadStatusResponse, error) {
	err := h.leadService.UpdateLeadStatus(req.Id, req.Status)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateLeadStatusResponse{}, nil
}

func (h *LeadHandler) ListLeads(ctx context.Context, req *pb.ListLeadsRequest) (*pb.ListLeadsResponse, error) {
	leads, err := h.leadService.ListLeads(req.Status, req.Source, int(req.Limit))
	if err != nil {
		return nil, err
	}
	var pbLeads []*pb.Lead
	for _, lead := range leads {
		pbLeads = append(pbLeads, &pb.Lead{
			Id:     lead.ID,
			Name:   lead.Name,
			Phone:  lead.Phone,
			Email:  lead.Email,
			Source: lead.Source,
			Status: lead.Status,
		})
	}
	return &pb.ListLeadsResponse{Leads: pbLeads}, nil
}
