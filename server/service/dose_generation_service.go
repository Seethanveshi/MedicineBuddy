package service

import (
	"MedicineBuddy/dto"
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
		model.DoseSkipped,
		&now,
	)
}

func contains(days []int, target int) bool {
	return slices.Contains(days, target)
}

func (s *DoseService) GetTodayDoses(ctx context.Context) ([]dto.DoseLogResponse, error) {

	nowLocal := time.Now().In(time.Local)

	startLocal := time.Date(
		nowLocal.Year(),
		nowLocal.Month(),
		nowLocal.Day(),
		0, 0, 0, 0,
		time.Local,
	)

	endLocal := startLocal.Add(24 * time.Hour)

	return s.doseRepository.GetDosesBetween(
		ctx,
		startLocal.UTC(),
		endLocal.UTC(),
	)
}

func (s *DoseService) GetDosesByDate(
	ctx context.Context,
	date time.Time,
) ([]dto.DoseLogResponse, error) {
	return s.doseRepository.GetDosesByDate(ctx, date)
}

func (s *DoseService) GetUpcomingDoses(
	ctx context.Context,
	days int,
) ([]dto.DoseLogResponse, error) {

	now := time.Now().UTC()
	end := now.AddDate(0, 0, days)

	return s.doseRepository.GetDosesBetween(ctx, now, end)
}

func (s *DoseService) GetDoseHistory(
	ctx context.Context,
	limit int,
) ([]dto.DoseLogResponse, error) {

	return s.doseRepository.GetDoseHistory(ctx, limit)
}
