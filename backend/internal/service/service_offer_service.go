package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/shardy678/pet-freelance/backend/internal/models"
	"github.com/shardy678/pet-freelance/backend/internal/repository"
)

type ServiceOfferService struct {
	repo *repository.ServiceOfferRepository
}

func NewServiceOfferService(r *repository.ServiceOfferRepository) *ServiceOfferService {
	return &ServiceOfferService{repo: r}
}

func (s *ServiceOfferService) CreateOffer(ctx context.Context, freelancerID uuid.UUID, inp *models.ServiceOffer) (*models.ServiceOffer, error) {
	inp.FreelancerID = freelancerID
	if err := s.repo.Create(ctx, inp); err != nil {
		return nil, err
	}
	return inp, nil
}

func (s *ServiceOfferService) GetOffer(ctx context.Context, id uuid.UUID) (*models.ServiceOffer, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ServiceOfferService) ListOffers(ctx context.Context) ([]models.ServiceOffer, error) {
	return s.repo.ListAll(ctx)
}

func (s *ServiceOfferService) ListByService(ctx context.Context, serviceID uuid.UUID) ([]models.ServiceOffer, error) {
	return s.repo.ListByService(ctx, serviceID)
}
