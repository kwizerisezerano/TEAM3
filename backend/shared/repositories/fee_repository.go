package repositories

import (
	"database/sql"

	"github.com/kwizera/university-invoices/services/fee/models"
)

type FeeRepository struct {
	db *sql.DB
}

func NewFeeRepository(db *sql.DB) *FeeRepository {
	return &FeeRepository{db: db}
}

func (r *FeeRepository) Create(feeType *models.FeeType) error {
	query := `INSERT INTO fee_types (name, description, amount, currency, is_recurring, is_active)
			  VALUES (?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, feeType.Name, feeType.Description, feeType.Amount,
		feeType.Currency, feeType.IsRecurring, feeType.IsActive)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	feeType.ID = int(id)
	return nil
}

func (r *FeeRepository) GetByID(id int) (*models.FeeType, error) {
	query := `SELECT id, name, description, amount, currency, is_recurring, is_active, created_at, updated_at
			  FROM fee_types WHERE id = ?`

	feeType := &models.FeeType{}
	err := r.db.QueryRow(query, id).Scan(
		&feeType.ID, &feeType.Name, &feeType.Description, &feeType.Amount,
		&feeType.Currency, &feeType.IsRecurring, &feeType.IsActive, &feeType.CreatedAt, &feeType.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return feeType, nil
}

func (r *FeeRepository) List(limit, offset int) ([]models.FeeType, error) {
	query := `SELECT id, name, description, amount, currency, is_recurring, is_active, created_at, updated_at
			  FROM fee_types WHERE is_active = true ORDER BY name LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeTypes []models.FeeType
	for rows.Next() {
		var feeType models.FeeType
		err := rows.Scan(
			&feeType.ID, &feeType.Name, &feeType.Description, &feeType.Amount,
			&feeType.Currency, &feeType.IsRecurring, &feeType.IsActive, &feeType.CreatedAt, &feeType.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		feeTypes = append(feeTypes, feeType)
	}

	return feeTypes, nil
}

func (r *FeeRepository) Update(feeType *models.FeeType) error {
	query := `UPDATE fee_types SET name = ?, description = ?, amount = ?, currency = ?, 
			  is_recurring = ?, is_active = ?, updated_at = NOW() WHERE id = ?`

	_, err := r.db.Exec(query, feeType.Name, feeType.Description, feeType.Amount,
		feeType.Currency, feeType.IsRecurring, feeType.IsActive, feeType.ID)
	return err
}

func (r *FeeRepository) Delete(id int) error {
	query := `DELETE FROM fee_types WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
