package repositories

import (
	"database/sql"

	"github.com/kwizera/university-invoices/services/payment/models"
)

type ReceiptRepository struct {
	db *sql.DB
}

func NewReceiptRepository(db *sql.DB) *ReceiptRepository {
	return &ReceiptRepository{db: db}
}

func (r *ReceiptRepository) Create(receipt *models.Receipt) error {
	query := `INSERT INTO receipts (receipt_number, payment_id, invoice_id, student_id, amount, issued_date, issued_by)
			  VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, receipt.ReceiptNumber, receipt.PaymentID, receipt.InvoiceID,
		receipt.StudentID, receipt.Amount, receipt.IssuedDate, receipt.IssuedBy)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	receipt.ID = int(id)
	return nil
}

func (r *ReceiptRepository) GetByID(id int) (*models.Receipt, error) {
	query := `SELECT id, receipt_number, payment_id, invoice_id, student_id, amount, issued_date, issued_by, created_at
			  FROM receipts WHERE id = ?`

	receipt := &models.Receipt{}
	err := r.db.QueryRow(query, id).Scan(
		&receipt.ID, &receipt.ReceiptNumber, &receipt.PaymentID, &receipt.InvoiceID,
		&receipt.StudentID, &receipt.Amount, &receipt.IssuedDate, &receipt.IssuedBy, &receipt.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return receipt, nil
}

func (r *ReceiptRepository) GetByPaymentID(paymentID int) (*models.Receipt, error) {
	query := `SELECT id, receipt_number, payment_id, invoice_id, student_id, amount, issued_date, issued_by, created_at
			  FROM receipts WHERE payment_id = ?`

	receipt := &models.Receipt{}
	err := r.db.QueryRow(query, paymentID).Scan(
		&receipt.ID, &receipt.ReceiptNumber, &receipt.PaymentID, &receipt.InvoiceID,
		&receipt.StudentID, &receipt.Amount, &receipt.IssuedDate, &receipt.IssuedBy, &receipt.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return receipt, nil
}

func (r *ReceiptRepository) GetByStudentID(studentID int, limit, offset int) ([]models.Receipt, error) {
	query := `SELECT id, receipt_number, payment_id, invoice_id, student_id, amount, issued_date, issued_by, created_at
			  FROM receipts WHERE student_id = ? ORDER BY issued_date DESC LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, studentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var receipts []models.Receipt
	for rows.Next() {
		var receipt models.Receipt
		err := rows.Scan(
			&receipt.ID, &receipt.ReceiptNumber, &receipt.PaymentID, &receipt.InvoiceID,
			&receipt.StudentID, &receipt.Amount, &receipt.IssuedDate, &receipt.IssuedBy, &receipt.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		receipts = append(receipts, receipt)
	}

	return receipts, nil
}
