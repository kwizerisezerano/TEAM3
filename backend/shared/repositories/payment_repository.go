package repositories

import (
	"database/sql"

	"github.com/kwizera/university-invoices/services/payment/models"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(payment *models.Payment) error {
	query := `INSERT INTO payments (payment_number, invoice_id, student_id, amount, payment_method,
			  payment_date, reference_number, status, notes, recorded_by)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, payment.PaymentNumber, payment.InvoiceID, payment.StudentID,
		payment.Amount, payment.PaymentMethod, payment.PaymentDate, payment.ReferenceNumber,
		payment.Status, payment.Notes, payment.RecordedBy)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	payment.ID = int(id)
	return nil
}

func (r *PaymentRepository) GetByID(id int) (*models.Payment, error) {
	query := `SELECT id, payment_number, invoice_id, student_id, amount, payment_method,
			  payment_date, reference_number, status, notes, recorded_by, created_at, updated_at
			  FROM payments WHERE id = ?`

	payment := &models.Payment{}
	err := r.db.QueryRow(query, id).Scan(
		&payment.ID, &payment.PaymentNumber, &payment.InvoiceID, &payment.StudentID,
		&payment.Amount, &payment.PaymentMethod, &payment.PaymentDate, &payment.ReferenceNumber,
		&payment.Status, &payment.Notes, &payment.RecordedBy, &payment.CreatedAt, &payment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (r *PaymentRepository) GetByInvoiceID(invoiceID int) ([]models.Payment, error) {
	query := `SELECT id, payment_number, invoice_id, student_id, amount, payment_method,
			  payment_date, reference_number, status, notes, recorded_by, created_at, updated_at
			  FROM payments WHERE invoice_id = ? ORDER BY payment_date DESC`

	rows, err := r.db.Query(query, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		err := rows.Scan(
			&payment.ID, &payment.PaymentNumber, &payment.InvoiceID, &payment.StudentID,
			&payment.Amount, &payment.PaymentMethod, &payment.PaymentDate, &payment.ReferenceNumber,
			&payment.Status, &payment.Notes, &payment.RecordedBy, &payment.CreatedAt, &payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

func (r *PaymentRepository) GetByStudentID(studentID int, limit, offset int) ([]models.Payment, error) {
	query := `SELECT id, payment_number, invoice_id, student_id, amount, payment_method,
			  payment_date, reference_number, status, notes, recorded_by, created_at, updated_at
			  FROM payments WHERE student_id = ? ORDER BY payment_date DESC LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, studentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		err := rows.Scan(
			&payment.ID, &payment.PaymentNumber, &payment.InvoiceID, &payment.StudentID,
			&payment.Amount, &payment.PaymentMethod, &payment.PaymentDate, &payment.ReferenceNumber,
			&payment.Status, &payment.Notes, &payment.RecordedBy, &payment.CreatedAt, &payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

func (r *PaymentRepository) List(limit, offset int) ([]models.Payment, error) {
	query := `SELECT id, payment_number, invoice_id, student_id, amount, payment_method,
			  payment_date, reference_number, status, notes, recorded_by, created_at, updated_at
			  FROM payments ORDER BY payment_date DESC LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		err := rows.Scan(
			&payment.ID, &payment.PaymentNumber, &payment.InvoiceID, &payment.StudentID,
			&payment.Amount, &payment.PaymentMethod, &payment.PaymentDate, &payment.ReferenceNumber,
			&payment.Status, &payment.Notes, &payment.RecordedBy, &payment.CreatedAt, &payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

func (r *PaymentRepository) Update(payment *models.Payment) error {
	query := `UPDATE payments SET status = ?, reference_number = ?, notes = ?, 
			  updated_at = NOW() WHERE id = ?`

	_, err := r.db.Exec(query, payment.Status, payment.ReferenceNumber, payment.Notes, payment.ID)
	return err
}

func (r *PaymentRepository) Delete(id int) error {
	query := `DELETE FROM payments WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
