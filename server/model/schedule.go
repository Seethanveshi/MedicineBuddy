package model

import (
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	ID         uuid.UUID
	MedicineID uuid.UUID
	TimeOfDay  time.Time
	DaysOfWeek []int
}
