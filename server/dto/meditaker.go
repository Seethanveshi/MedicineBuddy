package dto

import "github.com/google/uuid"

type MediTakerRequest struct {
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Relationship string `json:"relationship"`
}

type MediTakerResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Relationship string    `json:"relationship"`
}
