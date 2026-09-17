package service

import (
	"errors"

	"ticket-system/internal/model"
	"ticket-system/internal/repository"
)

type TicketService struct {
	TicketRepository *repository.TicketRepository
}

func NewTicketService(
	ticketRepository *repository.TicketRepository,
) *TicketService {
	return &TicketService{
		TicketRepository: ticketRepository,
	}
}

// CreateTicket creates a ticket for the logged-in user.
func (s *TicketService) CreateTicket(
	userID int64,
	title string,
	description string,
) (*model.Ticket, error) {

	return s.TicketRepository.CreateTicket(
		userID,
		title,
		description,
	)
}
func (s *TicketService) GetTicketsByUserID(userID int64) ([]model.Ticket, error) {
	return s.TicketRepository.GetTicketsByUserID(userID)
}
func (s *TicketService) GetTicketByIDAndUserID(ticketID int64, userID int64) (*model.Ticket, error) {
	return s.TicketRepository.GetTicketByIDAndUserID(ticketID, userID)
}
func (s *TicketService) UpdateTicketStatus(ticketID int64, userID int64, newStatus string) error {
	if newStatus != "open" &&
		newStatus != "in_progress" &&
		newStatus != "closed" {
		return errors.New("invalid status")
	}

	ticket, err := s.TicketRepository.GetTicketByIDAndUserID(ticketID, userID)

	if err != nil {
		return err
	}

	if ticket.Status == "closed" {
		return errors.New("closed ticket cannot be reopened or changed")
	}

	if ticket.Status == "open" && newStatus == "closed" {
		return errors.New("ticket must be in_progress before it can be closed")
	}

	if ticket.Status == "in_progress" && newStatus == "open" {
		return errors.New("in_progress ticket cannot be moved back to open")
	}

	return s.TicketRepository.UpdateTicketStatus(ticketID, userID, newStatus)
}
