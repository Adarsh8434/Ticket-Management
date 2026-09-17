package handler

import (
	"database/sql"
	"net/http"
	"strconv"
	"ticket-system/internal/model"
	"ticket-system/internal/service"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	TicketService *service.TicketService
}

type CreateTicketRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}
type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func NewTicketHandler(ticketService *service.TicketService) *TicketHandler {
	return &TicketHandler{
		TicketService: ticketService,
	}
}

// Create creates a new ticket for the logged-in user.
func (h *TicketHandler) Create(c *gin.Context) {

	var request CreateTicketRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	// Get user ID from JWT middleware
	userID, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	userIDInt64, ok := userID.(int64)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid user id",
		})
		return
	}

	// Create ticket
	ticket, err := h.TicketService.CreateTicket(
		userIDInt64,
		request.Title,
		request.Description,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create ticket",
		})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

func (h *TicketHandler) GetAll(c *gin.Context) {
	userID, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	userIDInt64, ok := userID.(int64)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid user id",
		})
		return
	}

	tickets, err := h.TicketService.GetTicketsByUserID(userIDInt64)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch tickets",
		})
		return
	}

	if tickets == nil {
		tickets = []model.Ticket{}
	}

	c.JSON(http.StatusOK, tickets)
}
func (h *TicketHandler) GetByID(c *gin.Context) {
	ticketID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid ticket id",
		})
		return
	}

	userID, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	userIDInt64, ok := userID.(int64)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid user id",
		})
		return
	}

	ticket, err := h.TicketService.GetTicketByIDAndUserID(ticketID, userIDInt64)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "ticket not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch ticket",
		})
		return
	}

	c.JSON(http.StatusOK, ticket)
}
func (h *TicketHandler) UpdateStatus(c *gin.Context) {
	ticketID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid ticket id",
		})
		return
	}

	var request UpdateStatusRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	userID, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	userIDInt64, ok := userID.(int64)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid user id",
		})
		return
	}

	err = h.TicketService.UpdateTicketStatus(
		ticketID,
		userIDInt64,
		request.Status,
	)

	if err != nil {
		switch err.Error() {
		case "invalid status":
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		case "ticket must be in_progress before it can be closed",
			"in_progress ticket cannot be moved back to open",
			"closed ticket cannot be reopened or changed":
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		case "sql: no rows in result set":
			c.JSON(http.StatusNotFound, gin.H{
				"error": "ticket not found",
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to update ticket status",
			})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ticket status updated successfully",
	})
}
