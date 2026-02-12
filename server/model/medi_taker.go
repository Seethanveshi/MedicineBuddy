package model

import (
	"time"

	"github.com/google/uuid"
)

type MediTaker struct {
	ID           uuid.UUID
	PatientID    uuid.UUID
	Name         string
	Email        string
	Relationship string
	CreatedAt    time.Time
}
