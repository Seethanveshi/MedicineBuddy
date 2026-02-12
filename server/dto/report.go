package dto

import "time"

type DoseDetail struct {
	MedicineName string     `json:"medicine_name"`
	Dosage       string     `json:"dosage"`
	ScheduledAt  time.Time  `json:"scheduled_at"`
	Status       string     `json:"status"`
	TakenAt      *time.Time `json:"taken_at"`
}

type DayReport struct {
	Date  string       `json:"date"`
	Doses []DoseDetail `json:"doses"`
}

type WeeklyDetailedReport struct {
	PatientName string      `json:"patient_name"`
	Days        []DayReport `json:"days"`
}
