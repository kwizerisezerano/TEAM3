package repositories

import (
	"database/sql"

	"github.com/kwizera/university-invoices/services/invoice/models"
)

type InvoiceRepository struct {
	db *sql.DB
}

func NewInvoiceRepository(db *sql.DB) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

func (r *InvoiceRepository) Create(invoice *models.Invoice) error {
	query := `INSERT INTO invoices (invoice_number, student_id, issue_date, due_date, subtotal, 
			  tax_amount, discount_amount, total_amount, paid_amount, outstanding_balance, 
			  status, notes, created_by)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, invoice.InvoiceNumber, invoice.StudentID, invoice.IssueDate,
		invoice.DueDate, invoice.Subtotal, invoice.TaxAmount, invoice.DiscountAmount,
		invoice.TotalAmount, invoice.PaidAmount, invoice.OutstandingBalance, invoice.Status,
		invoice.Notes, invoice.CreatedBy)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	invoice.ID = int(id)
	return nil
}

func (r *InvoiceRepository) GetByID(id int) (*models.Invoice, error) {
	query := `SELECT id, invoice_number, student_id, issue_date, due_date, subtotal, tax_amount,
			  discount_amount, total_amount, paid_amount, outstanding_balance, status, notes,
			  created_by, created_at, updated_at
			  FROM invoices WHERE id = ?`

	invoice := &models.Invoice{}
	err := r.db.QueryRow(query, id).Scan(
		&invoice.ID, &invoice.InvoiceNumber, &invoice.StudentID, &invoice.IssueDate,
		&invoice.DueDate, &invoice.Subtotal, &invoice.TaxAmount, &invoice.DiscountAmount,
		&invoice.TotalAmount, &invoice.PaidAmount, &invoice.OutstandingBalance, &invoice.Status,
		&invoice.Notes, &invoice.CreatedBy, &invoice.CreatedAt, &invoice.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func (r *InvoiceRepository) GetByInvoiceNumber(invoiceNumber string) (*models.Invoice, error) {
	query := `SELECT id, invoice_number, student_id, issue_date, due_date, subtotal, tax_amount,
			  discount_amount, total_amount, paid_amount, outstanding_balance, status, notes,
			  created_by, created_at, updated_at
			  FROM invoices WHERE invoice_number = ?`

	invoice := &models.Invoice{}
	err := r.db.QueryRow(query, invoiceNumber).Scan(
		&invoice.ID, &invoice.InvoiceNumber, &invoice.StudentID, &invoice.IssueDate,
		&invoice.DueDate, &invoice.Subtotal, &invoice.TaxAmount, &invoice.DiscountAmount,
		&invoice.TotalAmount, &invoice.PaidAmount, &invoice.OutstandingBalance, &invoice.Status,
		&invoice.Notes, &invoice.CreatedBy, &invoice.CreatedAt, &invoice.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func (r *InvoiceRepository) GetByStudentID(studentID int, limit, offset int) ([]models.Invoice, error) {
	query := `SELECT id, invoice_number, student_id, issue_date, due_date, subtotal, tax_amount,
			  discount_amount, total_amount, paid_amount, outstanding_balance, status, notes,
			  created_by, created_at, updated_at
			  FROM invoices WHERE student_id = ? ORDER BY issue_date DESC LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, studentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []models.Invoice
	for rows.Next() {
		var invoice models.Invoice
		err := rows.Scan(
			&invoice.ID, &invoice.InvoiceNumber, &invoice.StudentID, &invoice.IssueDate,
			&invoice.DueDate, &invoice.Subtotal, &invoice.TaxAmount, &invoice.DiscountAmount,
			&invoice.TotalAmount, &invoice.PaidAmount, &invoice.OutstandingBalance, &invoice.Status,
			&invoice.Notes, &invoice.CreatedBy, &invoice.CreatedAt, &invoice.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}

	return invoices, nil
}

func (r *InvoiceRepository) List(limit, offset int) ([]models.Invoice, error) {
	query := `SELECT id, invoice_number, student_id, issue_date, due_date, subtotal, tax_amount,
			  discount_amount, total_amount, paid_amount, outstanding_balance, status, notes,
			  created_by, created_at, updated_at
			  FROM invoices ORDER BY issue_date DESC LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []models.Invoice
	for rows.Next() {
		var invoice models.Invoice
		err := rows.Scan(
			&invoice.ID, &invoice.InvoiceNumber, &invoice.StudentID, &invoice.IssueDate,
			&invoice.DueDate, &invoice.Subtotal, &invoice.TaxAmount, &invoice.DiscountAmount,
			&invoice.TotalAmount, &invoice.PaidAmount, &invoice.OutstandingBalance, &invoice.Status,
			&invoice.Notes, &invoice.CreatedBy, &invoice.CreatedAt, &invoice.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}

	return invoices, nil
}

func (r *InvoiceRepository) Update(invoice *models.Invoice) error {
	query := `UPDATE invoices SET due_date = ?, discount_amount = ?, status = ?, 
			  notes = ?, updated_at = NOW() WHERE id = ?`

	_, err := r.db.Exec(query, invoice.DueDate, invoice.DiscountAmount, invoice.Status,
		invoice.Notes, invoice.ID)
	return err
}

func (r *InvoiceRepository) UpdatePaymentAmounts(invoiceID int, paidAmount, outstandingBalance float64) error {
	query := `UPDATE invoices SET paid_amount = ?, outstanding_balance = ?, 
			  updated_at = NOW() WHERE id = ?`

	_, err := r.db.Exec(query, paidAmount, outstandingBalance, invoiceID)
	return err
}

func (r *InvoiceRepository) UpdateInvoiceStatus(invoiceID int, status models.InvoiceStatus) error {
	query := `UPDATE invoices SET status = ?, updated_at = NOW() WHERE id = ?`

	_, err := r.db.Exec(query, status, invoiceID)
	return err
}

func (r *InvoiceRepository) Delete(id int) error {
	query := `DELETE FROM invoices WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
