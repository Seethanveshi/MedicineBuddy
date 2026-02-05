package model

import (
	"time"

	"github.com/google/uuid"
)

type Medicine struct {
	ID uuid.UUID
	UserID uuid.UUID
	Name string
	Dosage string
	StartDate time.Time
	EndDate *time.Time
	CreatedAt time.Time
}