package repository

import (
	"MedicineBuddy/dto"
	"MedicineBuddy/model"
	"context"

	"github.com/google/uuid"
)

type MediTakerRepository interface {
	Create(ctx context.Context, meditaker model.MediTaker) error
	ListByPatient(ctx context.Context, patientID uuid.UUID) ([]dto.MediTakerResponse, error)
	Update(ctx context.Context, id uuid.UUID, patientID uuid.UUID, medireq dto.MediTakerRequest) error
	Delete(ctx context.Context, id uuid.UUID, patientID uuid.UUID) error
}
