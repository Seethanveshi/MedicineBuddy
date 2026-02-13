package service

import (
	"MedicineBuddy/dto"
	"MedicineBuddy/model"
	"MedicineBuddy/repository"
	"context"

	"github.com/google/uuid"
)

type MediTakerService struct {
	meditakerRepository repository.MediTakerRepository
}

func NewMediTakerService(meditakerRepository repository.MediTakerRepository) *MediTakerService {
	return &MediTakerService{
		meditakerRepository: meditakerRepository,
	}
}

func (s *MediTakerService) CreateMediTaker(
	ctx context.Context,
	patientID uuid.UUID,
	req dto.MediTakerRequest,
) error {

	m := model.MediTaker{
		ID:           uuid.New(),
		PatientID:    patientID,
		Name:         req.Name,
		Email:        req.Email,
		Relationship: req.Relationship,
	}

	return s.meditakerRepository.Create(ctx, m)
}

func (s *MediTakerService) ListByPatient(ctx context.Context, patientID uuid.UUID) ([]dto.MediTakerResponse, error) {
	return s.meditakerRepository.ListByPatient(ctx, patientID)
}

func (s *MediTakerService) Delete(ctx context.Context, id uuid.UUID, patientID uuid.UUID,
) error {
	return s.meditakerRepository.Delete(ctx, id, patientID)
}

func (s *MediTakerService) Update(ctx context.Context, id uuid.UUID, patientID uuid.UUID, medireq dto.MediTakerRequest) error {
	return s.meditakerRepository.Update(ctx, id, patientID, medireq)
}
