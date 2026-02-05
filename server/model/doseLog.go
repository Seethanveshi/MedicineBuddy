package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	DosePending = "pending"
	DoseTaken   = "taken"
	DoseMissed  = "missed"
)

type DoseLog struct {
	ID         uuid.UUID
	MedicineID uuid.UUID
	ScheduleAt time.Time
	Status     string
	TakenAt    time.Time
}
