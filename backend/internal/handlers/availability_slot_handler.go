package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shardy678/pet-freelance/backend/internal/service"
)

type AvailabilitySlotHandler struct {
	svc *service.AvailabilitySlotService
}

func NewAvailabilitySlotHandler(s *service.AvailabilitySlotService) *AvailabilitySlotHandler {
	return &AvailabilitySlotHandler{svc: s}
}

type createSlotReq struct {
	StartTime string `json:"start_time" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EndTime   string `json:"end_time"   binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

func (h *AvailabilitySlotHandler) Create(c *gin.Context) {
	offerID, err := uuid.Parse(c.Param("offer_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offer_id"})
		return
	}

	var req createSlotReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	start, _ := time.Parse(time.RFC3339, req.StartTime)
	end, _ := time.Parse(time.RFC3339, req.EndTime)

	slot, err := h.svc.CreateSlot(c.Request.Context(), offerID, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, slot)
}

func (h *AvailabilitySlotHandler) List(c *gin.Context) {
	offerID, err := uuid.Parse(c.Param("offer_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offer_id"})
		return
	}

	fromStr := c.Query("from")
	toStr := c.Query("to")
	onlyAvail := c.DefaultQuery("available", "true") == "true"

	from, _ := time.Parse(time.RFC3339, fromStr)
	to, _ := time.Parse(time.RFC3339, toStr)

	slots, err := h.svc.ListSlots(c.Request.Context(), offerID, onlyAvail, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, slots)
}

type updateSlotReq struct {
	StartTime string `json:"start_time" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	EndTime   string `json:"end_time"   binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	IsBooked  *bool  `json:"is_booked"  binding:"omitempty"`
}

func (h *AvailabilitySlotHandler) Update(c *gin.Context) {
	slotID, err := uuid.Parse(c.Param("slot_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid slot_id"})
		return
	}
	var req updateSlotReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	slot, err := h.svc.UpdateSlot(
		c.Request.Context(),
		slotID,
		parseOrDefault(req.StartTime, time.RFC3339),
		parseOrDefault(req.EndTime, time.RFC3339),
		boolOrDefault(req.IsBooked, false),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, slot)
}

func (h *AvailabilitySlotHandler) Delete(c *gin.Context) {
	slotID, err := uuid.Parse(c.Param("slot_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid slot_id"})
		return
	}
	if err := h.svc.DeleteSlot(c.Request.Context(), slotID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// helpers

func parseOrDefault(val, layout string) time.Time {
	if val == "" {
		return time.Time{}
	}
	t, _ := time.Parse(layout, val)
	return t
}

func boolOrDefault(ptr *bool, def bool) bool {
	if ptr == nil {
		return def
	}
	return *ptr
}
