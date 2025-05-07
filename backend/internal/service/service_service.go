package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/shardy678/pet-freelance/backend/internal/models"
	"github.com/shardy678/pet-freelance/backend/internal/repository"
)

var ErrServiceNotFound = errors.New("service not found")

type ServiceService struct {
	repo *repository.ServiceRepository
}

func NewServiceService(r *repository.ServiceRepository) *ServiceService {
	return &ServiceService{repo: r}
}

// Create creates a new Service.
func (s *ServiceService) Create(ctx context.Context, svc *models.Service) (*models.Service, error) {
	// assign UUID if not set
	if svc.ID == uuid.Nil {
		svc.ID = uuid.New()
	}
	if err := s.repo.Create(ctx, svc); err != nil {
		return nil, err
	}
	return svc, nil
}

// Get fetches a Service by its ID.
func (s *ServiceService) Get(ctx context.Context, id uuid.UUID) (*models.Service, error) {
	svc, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrServiceNotFound
	}
	return svc, nil
}

// List returns all non-deleted Services.
func (s *ServiceService) List(ctx context.Context) ([]models.Service, error) {
	return s.repo.ListAll(ctx)
}
