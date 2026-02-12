package repository

import (
	"MedicineBuddy/dto"
	"MedicineBuddy/model"
	"context"
	"time"

	"github.com/google/uuid"
)

type DoseRepository interface {
	Exists(ctx context.Context, medicineID uuid.UUID, scheduleAt time.Time) (bool, error)
	Create(ctx context.Context, dose *model.DoseLog) error
	MarkMissedDose(ctx context.Context, now time.Time) error
	UpdateDoseStatus(ctx context.Context, doseID uuid.UUID, fromStatus string, toStatus string, takenAt *time.Time) error
	GetDosesBetween(ctx context.Context, from time.Time, to time.Time) ([]dto.DoseLogResponse, error)
	GetDoseHistory(ctx context.Context, limit int) ([]dto.DoseLogResponse, error)
	GetDosesByDate(ctx context.Context, date time.Time) ([]dto.DoseLogResponse, error)
	WeeklyDetailed(ctx context.Context, patientID uuid.UUID, start time.Time, end time.Time) (string, []dto.DoseDetail, error)
}
