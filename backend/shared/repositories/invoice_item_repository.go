package repositories

import (
	"database/sql"

	"github.com/kwizera/university-invoices/services/invoice/models"
)

type InvoiceItemRepository struct {
	db *sql.DB
}

func NewInvoiceItemRepository(db *sql.DB) *InvoiceItemRepository {
	return &InvoiceItemRepository{db: db}
}

func (r *InvoiceItemRepository) Create(item *models.InvoiceItem) error {
	query := `INSERT INTO invoice_items (invoice_id, fee_type_id, description, quantity, unit_price, amount)
			  VALUES (?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, item.InvoiceID, item.FeeTypeID, item.Description,
		item.Quantity, item.UnitPrice, item.Amount)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	item.ID = int(id)
	return nil
}

func (r *InvoiceItemRepository) GetByInvoiceID(invoiceID int) ([]models.InvoiceItem, error) {
	query := `SELECT id, invoice_id, fee_type_id, description, quantity, unit_price, amount, created_at
			  FROM invoice_items WHERE invoice_id = ?`

	rows, err := r.db.Query(query, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.InvoiceItem
	for rows.Next() {
		var item models.InvoiceItem
		err := rows.Scan(
			&item.ID, &item.InvoiceID, &item.FeeTypeID, &item.Description,
			&item.Quantity, &item.UnitPrice, &item.Amount, &item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *InvoiceItemRepository) DeleteByInvoiceID(invoiceID int) error {
	query := `DELETE FROM invoice_items WHERE invoice_id = ?`
	_, err := r.db.Exec(query, invoiceID)
	return err
}
