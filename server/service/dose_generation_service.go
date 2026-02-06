package service

import (
	"MedicineBuddy/model"
	"MedicineBuddy/repository"
	"context"
	"slices"
	"time"

	"github.com/google/uuid"
)

type DoseService struct {
	doseRepository repository.DoseRepository
}

func NewDoseService(repo repository.DoseRepository) *DoseService {
	return &DoseService{doseRepository: repo}
}

func (s *DoseService) GenerateUpcomingDoses(
	ctx context.Context,
	medicine *model.Medicine,
	schedule *model.Schedule,
	daysAhead int,
) error {
	now := time.Now().UTC()
	start := now.Truncate(24 * time.Hour)
	end := start.AddDate(0, 0, daysAhead)

	for day := start; day.Before(end); day = day.AddDate(0, 0, 1) {
		if day.Before(medicine.StartDate) {
			continue
		}

		if medicine.EndDate != nil && day.After(*medicine.EndDate) {
			continue
		}

		weekday := int(day.Weekday())

		if !contains(schedule.DaysOfWeek, weekday) {
			continue
		}

		scheduledLocal := time.Date(
			day.Year(),
			day.Month(),
			day.Day(),
			schedule.TimeOfDay.Hour(),
			schedule.TimeOfDay.Minute(),
			0,
			0,
			time.Local,
		)

		scheduledAt := scheduledLocal.UTC()

		if scheduledAt.Before(now) {
			continue
		}

		exists, err := s.doseRepository.Exists(ctx, medicine.ID, scheduledAt)
		if err != nil {
			return err
		}
		if exists {
			continue
		}

		dose := &model.DoseLog{
			ID:         uuid.New(),
			MedicineID: medicine.ID,
			ScheduleAt: scheduledAt,
			Status:     model.DosePending,
		}

		if err := s.doseRepository.Create(ctx, dose); err != nil {
			return err
		}
	}

	return nil
}

func (s *DoseService) MarkMissedDose(ctx context.Context) error {
	now := time.Now().UTC()
	return s.doseRepository.MarkMissedDose(ctx, now)
}

func (s *DoseService) MarkDoseTaken(ctx context.Context, doseID uuid.UUID) error {
	now := time.Now().UTC()

	return s.doseRepository.UpdateDoseStatus(
		ctx,
		doseID,
		model.DosePending,
		model.DoseTaken,
		&now,
	)
}

func (s *DoseService) MarkDoseSkipped(ctx context.Context, doseID uuid.UUID) error {
	now := time.Now().UTC()

	return s.doseRepository.UpdateDoseStatus(
		ctx,
		doseID,
		model.DosePending,
		model.DoseMissed,
		&now,
	)
}

func contains(days []int, target int) bool {
	return slices.Contains(days, target)
}