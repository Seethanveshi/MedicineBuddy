package repository

import (
	"MedicineBuddy/model"
	"context"
	"database/sql"
)

type MedicineRepository struct {
	db *sql.DB
}

func NewMedicineRepository(db *sql.DB) *MedicineRepository {
	return &MedicineRepository{db: db}
}

func (r *MedicineRepository) CreateMedicineWithSchedule(ctx context.Context, medicine *model.Medicine, schedule *model.Schedule) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := `
		INSERT INTO medicines(id, user_id, name, dosage, start_date, end_date) 
		VALUES($1, $2, $3, $4, $5, $6)	
	`

	_, err = tx.ExecContext(
		ctx,
		query,
		medicine.ID,
		medicine.UserID,
		medicine.Name,
		medicine.Dosage,
		medicine.StartDate,
		medicine.EndDate,
	)
	if err != nil {
		return err
	}

	query1 := `
		INSERT INTO schedules(id, medicine_id, time_of_day, days_of_week)
		VALUES($1, $2, $3, $4)
	`

	_, err = tx.ExecContext(
		ctx,
		query1,
		schedule.ID,
		schedule.MedicineID,
		schedule.TimeOfDay,
		schedule.DaysOfWeek,
	)
	if err != nil {
		return err
	}

	return tx.Commit()

}
