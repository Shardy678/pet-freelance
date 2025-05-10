package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/shardy678/pet-freelance/backend/internal/models"
	"github.com/shardy678/pet-freelance/backend/internal/repository"
	"github.com/shardy678/pet-freelance/backend/internal/service"
)

type ServiceOfferHandler struct {
	repo *repository.ServiceOfferRepository
	svc  *service.ServiceOfferService
}

func NewServiceOfferHandler(
	repo *repository.ServiceOfferRepository,
	svc *service.ServiceOfferService,
) *ServiceOfferHandler {
	return &ServiceOfferHandler{repo: repo, svc: svc}
}

type createServiceOfferReq struct {
	ServiceID           string  `json:"service_id" binding:"required,uuid"`
	Title               string  `json:"title" binding:"required"`
	Description         string  `json:"description" binding:"required"`
	Price               float32 `json:"price" binding:"required,gt=0"`
	Currency            string  `json:"currency" binding:"required,len=3"`
	PriceType           string  `json:"price_type" binding:"required,oneof=hourly fixed"`
	DurationEstimateMin int     `json:"duration_estimate_min" binding:"omitempty,gt=0"`
}

// Create handles POST /offers
func (h *ServiceOfferHandler) Create(c *gin.Context) {
	var req createServiceOfferReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract freelancer ID from JWT claims
	uid, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	freelancerID, err := uuid.Parse(uid.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	// Parse service ID
	serviceID, err := uuid.Parse(req.ServiceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service_id"})
		return
	}

	offer := &models.ServiceOffer{
		ID:                  uuid.New(),
		FreelancerID:        freelancerID,
		ServiceID:           serviceID,
		Title:               req.Title,
		Description:         req.Description,
		Price:               req.Price,
		Currency:            req.Currency,
		PriceType:           req.PriceType,
		DurationEstimateMin: req.DurationEstimateMin,
	}

	if err := h.repo.Create(c.Request.Context(), offer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, offer)
}

// List handles GET /offers and supports optional ?service_id=…
func (h *ServiceOfferHandler) List(c *gin.Context) {
	// Check for an optional service_id query parameter
	svcParam := c.Query("service_id")
	var offers []models.ServiceOffer
	var err error

	if svcParam != "" {
		// parse and fetch only those offers
		svcID, parseErr := uuid.Parse(svcParam)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service_id"})
			return
		}
		offers, err = h.svc.ListByService(c.Request.Context(), svcID)
	} else {
		// no filter → list all active offers
		offers, err = h.svc.ListOffers(c.Request.Context())
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, offers)
}

// Get handles GET /offers/:id
func (h *ServiceOfferHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offer id"})
		return
	}

	offer, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "offer not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, offer)
}

// ListByService handles GET /services/:service_id/offers
func (h *ServiceOfferHandler) ListByService(c *gin.Context) {
	serviceID, err := uuid.Parse(c.Param("service_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service_id"})
		return
	}

	offers, err := h.svc.ListByService(c.Request.Context(), serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, offers)
}
