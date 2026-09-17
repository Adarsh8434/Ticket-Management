package repository

import (
	"database/sql"

	"ticket-system/internal/model"
)

type TicketRepository struct {
	DB *sql.DB
}

func NewTicketRepository(db *sql.DB) *TicketRepository {
	return &TicketRepository{
		DB: db,
	}
}

// CreateTicket creates a new ticket for a user.
func (r *TicketRepository) CreateTicket(
	userID int64,
	title string,
	description string,
) (*model.Ticket, error) {

	result, err := r.DB.Exec(`
		INSERT INTO tickets (user_id, title, description, status)
		VALUES (?, ?, ?, 'open')
	`, userID, title, description)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &model.Ticket{
		ID:          id,
		UserID:      userID,
		Title:       title,
		Description: description,
		Status:      "open",
	}, nil
}
func (r *TicketRepository) GetTicketsByUserID(userID int64) ([]model.Ticket, error) {
	rows, err := r.DB.Query(`
		SELECT id, user_id, title, description, status
		FROM tickets
		WHERE user_id = ?
		ORDER BY id DESC
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []model.Ticket

	for rows.Next() {
		var ticket model.Ticket

		err := rows.Scan(
			&ticket.ID,
			&ticket.UserID,
			&ticket.Title,
			&ticket.Description,
			&ticket.Status,
		)

		if err != nil {
			return nil, err
		}

		tickets = append(tickets, ticket)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tickets, nil
}
func (r *TicketRepository) GetTicketByIDAndUserID(ticketID int64, userID int64) (*model.Ticket, error) {
	var ticket model.Ticket

	err := r.DB.QueryRow(`
		SELECT id, user_id, title, description, status
		FROM tickets
		WHERE id = ? AND user_id = ?
	`, ticketID, userID).Scan(
		&ticket.ID,
		&ticket.UserID,
		&ticket.Title,
		&ticket.Description,
		&ticket.Status,
	)

	if err != nil {
		return nil, err
	}

	return &ticket, nil
}
func (r *TicketRepository) UpdateTicketStatus(ticketID int64, userID int64, status string) error {
	result, err := r.DB.Exec(`
		UPDATE tickets
		SET status = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND user_id = ?
	`, status, ticketID, userID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
