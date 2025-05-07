package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shardy678/pet-freelance/backend/internal/handlers"
	"github.com/shardy678/pet-freelance/backend/internal/models"
	"github.com/shardy678/pet-freelance/backend/internal/repository"
	"github.com/shardy678/pet-freelance/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupServiceRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	// use Gin in test mode
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// in-memory sqlite
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// migrate only the Service model
	err = db.AutoMigrate(&models.Service{})
	assert.NoError(t, err)

	// wire up
	repo := repository.NewServiceRepository(db)
	svc := service.NewServiceService(repo)
	h := handlers.NewServiceHandler(svc)

	r.POST("/services", h.Create)
	r.GET("/services", h.List)
	r.GET("/services/:id", h.Get)

	return r, db
}

func TestCreateAndGetService(t *testing.T) {
	router, _ := setupServiceRouter(t)

	// 1) Create a new Service
	payload := map[string]interface{}{
		"name":                 "Grooming",
		"description":          "Full-service pet grooming",
		"base_price":           30.5,
		"default_duration_min": 90,
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/services", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var created models.Service
	err := json.Unmarshal(w.Body.Bytes(), &created)
	assert.NoError(t, err)
	assert.Equal(t, "Grooming", created.Name)
	assert.Equal(t, 30.5, created.BasePrice)
	assert.Equal(t, 90, created.DefaultDurationMin)
	assert.NotEqual(t, uuid.Nil, created.ID)

	// 2) Get it back by ID
	url := "/services/" + created.ID.String()
	req2 := httptest.NewRequest(http.MethodGet, url, nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
	var fetched models.Service
	err = json.Unmarshal(w2.Body.Bytes(), &fetched)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
	assert.Equal(t, "Grooming", fetched.Name)
}

func TestListServices(t *testing.T) {
	router, db := setupServiceRouter(t)

	// seed two services
	services := []models.Service{
		{ID: uuid.New(), Name: "Walking", BasePrice: 10, DefaultDurationMin: 30},
		{ID: uuid.New(), Name: "Training", BasePrice: 50, DefaultDurationMin: 60},
	}
	err := db.Create(&services).Error
	assert.NoError(t, err)

	// list
	req := httptest.NewRequest(http.MethodGet, "/services", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var list []models.Service
	err = json.Unmarshal(w.Body.Bytes(), &list)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
	// since handler orders by name asc:
	assert.Equal(t, "Training", list[0].Name)
	assert.Equal(t, "Walking", list[1].Name)
}
