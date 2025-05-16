package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shardy678/pet-freelance/backend/internal/service"
)

type ActivityHandler struct {
	svc *service.ActivityService
}

func NewActivityHandler(svc *service.ActivityService) *ActivityHandler {
	return &ActivityHandler{svc}
}

// GET /api/activities?limit=10
func (h *ActivityHandler) List(c *gin.Context) {
	uid, _ := c.Get("uid")
	userID, err := uuid.Parse(uid.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}
	limit := 10
	if l := c.Query("limit"); l != "" {
		fmt.Sscan(l, &limit)
	}
	acts, err := h.svc.List(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, acts)
}
