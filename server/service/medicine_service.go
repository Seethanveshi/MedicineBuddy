package service

import (
	"context"
	"time"

	"MedicineBuddy/model"
	"MedicineBuddy/repository"

	"github.com/google/uuid"
)

type MedicineService struct {
	medicineRepo *repository.MedicineRepository
	doseService  *DoseService
}

func NewMedicineService(
	medRepo *repository.MedicineRepository,
	doseSvc *DoseService,
) *MedicineService {
	return &MedicineService{
		medicineRepo: medRepo,
		doseService:  doseSvc,
	}
}

func (s *MedicineService) CreateMedicine(
	ctx context.Context,
	name string,
	dosage string,
	startDate time.Time,
	endDate *time.Time,
	timeOfDay time.Time,
	daysOfWeek []int,
) error {

	medicineID := uuid.New()
	medicine := &model.Medicine{
		ID:        medicineID,
		UserID:    uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Name:      name,
		Dosage:    dosage,
		StartDate: startDate,
		EndDate:   endDate,
	}

	schedule := &model.Schedule{
		ID:         uuid.New(),
		MedicineID: medicineID,
		TimeOfDay:  timeOfDay,
		DaysOfWeek: daysOfWeek,
	}

	if err := s.medicineRepo.CreateMedicineWithSchedule(ctx, medicine, schedule); err != nil {
		return err
	}

	// 🔥 Trigger dose generation
	return s.doseService.GenerateUpcomingDoses(ctx, medicine, schedule, 7)
}


