package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shardy678/pet-freelance/backend/internal/repository"
)

type ProfileHandler struct{ users *repository.UserRepository }

func NewProfileHandler(u *repository.UserRepository) *ProfileHandler { return &ProfileHandler{u} }

func (h *ProfileHandler) Me(c *gin.Context) {
	uid := c.GetString("uid")
	user, err := h.users.FindByID(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
	})
}
