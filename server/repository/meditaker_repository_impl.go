package repository

import (
	"MedicineBuddy/dto"
	"MedicineBuddy/model"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type MediTakerImpl struct {
	db *sql.DB
}

func NewMediTakerRepository(db *sql.DB) MediTakerRepository {
	return &MediTakerImpl{db: db}
}

func (r *MediTakerImpl) Create(ctx context.Context, meditaker model.MediTaker) error {
	query := `
	INSERT INTO meditakers (
		id,
		patient_id,
		name,
		email,
		relationship,
		created_at
	)
	VALUES ($1, $2, $3, $4, $5, Now())
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		meditaker.ID,
		meditaker.PatientID,
		meditaker.Name,
		meditaker.Email,
		meditaker.Relationship,
	)

	return err
}

func (r *MediTakerImpl) ListByPatient(ctx context.Context, patientID uuid.UUID) ([]dto.MediTakerResponse, error) {

	query := `
		SELECT id, name, email, relationship
		FROM meditakers
		WHERE patient_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.MediTakerResponse

	for rows.Next() {
		var m dto.MediTakerResponse

		if err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.Email,
			&m.Relationship,
		); err != nil {
			return nil, err
		}

		result = append(result, m)
	}

	return result, nil
}

func (r *MediTakerImpl) Delete(ctx context.Context, id uuid.UUID, patientID uuid.UUID) error {

	query := `
		DELETE FROM meditakers
		WHERE id = $1
		AND patient_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, id, patientID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("not found")
	}

	return nil
}

func (r *MediTakerImpl) Update(
	ctx context.Context,
	id uuid.UUID,
	patientID uuid.UUID,
	req dto.MediTakerRequest,
) error {

	query := `
		UPDATE meditakers
		SET name = $1,
		    email = $2,
		    relationship = $3
		WHERE id = $4
		AND patient_id = $5
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		req.Name,
		req.Email,
		req.Relationship,
		id,
		patientID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("not found")
	}

	return nil
}
