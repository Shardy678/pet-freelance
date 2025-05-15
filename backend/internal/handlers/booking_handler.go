package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shardy678/pet-freelance/backend/internal/service"
)

type BookingHandler struct {
	svc *service.BookingService
}

func NewBookingHandler(svc *service.BookingService) *BookingHandler {
	return &BookingHandler{svc}
}

type createBookingReq struct {
	OfferID string `json:"offer_id" binding:"required,uuid"`
	SlotID  string `json:"slot_id"  binding:"required,uuid"`
}

func (h *BookingHandler) Create(c *gin.Context) {
	var req createBookingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse IDs
	offerID, err := uuid.Parse(req.OfferID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offer_id"})
		return
	}
	slotID, err := uuid.Parse(req.SlotID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid slot_id"})
		return
	}

	// Extract user (owner) ID from JWT
	uid, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	ownerID, err := uuid.Parse(uid.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	booking, err := h.svc.BookSlot(c.Request.Context(), offerID, slotID, ownerID)
	if err != nil {
		if err == service.ErrSlotAlreadyBooked {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, booking)
}

func (h *BookingHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}
	b, err := h.svc.GetBooking(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}
	c.JSON(http.StatusOK, b)
}

func (h *BookingHandler) List(c *gin.Context) {
	// List bookings for the logged-in user
	uid, _ := c.Get("uid")
	ownerID, _ := uuid.Parse(uid.(string))

	list, err := h.svc.ListByOwner(c.Request.Context(), ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}
