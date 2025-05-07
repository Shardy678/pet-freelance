package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shardy678/pet-freelance/backend/internal/models"
	"github.com/shardy678/pet-freelance/backend/internal/service"
)

type ServiceHandler struct {
	svc *service.ServiceService
}

func NewServiceHandler(s *service.ServiceService) *ServiceHandler {
	return &ServiceHandler{svc: s}
}

type createServiceReq struct {
	Name               string  `json:"name" binding:"required"`
	Description        string  `json:"description"`
	BasePrice          float64 `json:"base_price" binding:"gte=0"`
	DefaultDurationMin int     `json:"default_duration_min" binding:"gte=1"`
}

// Create handles POST /services
func (h *ServiceHandler) Create(c *gin.Context) {
	var req createServiceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	svc := &models.Service{
		Name:               req.Name,
		Description:        req.Description,
		BasePrice:          req.BasePrice,
		DefaultDurationMin: req.DefaultDurationMin,
	}

	created, err := h.svc.Create(c.Request.Context(), svc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// List handles GET /services
func (h *ServiceHandler) List(c *gin.Context) {
	list, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// Get handles GET /services/:id
func (h *ServiceHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service id"})
		return
	}

	svc, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
		return
	}
	c.JSON(http.StatusOK, svc)
}
