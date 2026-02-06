package repository

import (
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
}
