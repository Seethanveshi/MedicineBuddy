package service

import (
	"context"
	"database/sql"
	"time"

	"MedicineBuddy/dto"
	"MedicineBuddy/model"
	"MedicineBuddy/repository"

	"github.com/google/uuid"
)

type MedicineService struct {
	db                 *sql.DB
	medicineRepository *repository.MedicineRepository
	doseService        *DoseService
	doseRepository     *repository.DoseRepository
}

func NewMedicineService(db *sql.DB,
	medRepo *repository.MedicineRepository,
	doseSvc *DoseService,
	doseRepository *repository.DoseRepository,

) *MedicineService {
	return &MedicineService{
		db:                 db,
		medicineRepository: medRepo,
		doseService:        doseSvc,
		doseRepository:     doseRepository,
	}
}

func (s *MedicineService) GetByID(
	ctx context.Context,
	id uuid.UUID,
	patientID uuid.UUID,
) (dto.MedicineDetailResponse, error) {
	return s.medicineRepository.GetByID(ctx, id, patientID)
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

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()
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

	if err := s.medicineRepository.CreateMedicineWithSchedule(ctx, medicine, schedule); err != nil {
		return err
	}

	return s.doseService.GenerateUpcomingDosesTx(ctx, tx, medicine, schedule, 7)
}

func (s *MedicineService) Update(
	ctx context.Context,
	medicineID uuid.UUID,
	patientID uuid.UUID,
	req dto.CreateMedicineRequest,
) error {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// 1️⃣ update medicine
	err = s.medicineRepository.UpdateTx(ctx, tx, medicineID, patientID, req)
	if err != nil {
		return err
	}

	// 2️⃣ update schedule
	err = s.medicineRepository.UpdateScheduleByMedicineTx(
		ctx,
		tx,
		medicineID,
		req.Schedule,
	)
	if err != nil {
		return err
	}

	// 3️⃣ delete future pending doses
	err = s.medicineRepository.DeleteFutureByMedicine(
		ctx,
		tx,
		medicineID,
		time.Now(),
	)
	if err != nil {
		return err
	}

	// 4️⃣ fetch updated medicine & schedule (important)
	medicine, err := s.medicineRepository.GetModelByIDTx(ctx, tx, medicineID)
	if err != nil {
		return err
	}

	schedule, err := s.medicineRepository.GetScheduleByMedicineTx(ctx, tx, medicineID)
	if err != nil {
		return err
	}

	err = s.doseService.GenerateUpcomingDosesTx(
		ctx,
		tx,
		medicine,
		schedule,
		7,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}