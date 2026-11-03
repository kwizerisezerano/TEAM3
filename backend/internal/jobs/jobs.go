package jobs

import (
	"database/sql"
	"log"
	"time"

	"github.com/kwizera/university-invoices/shared/repositories"
)

type JobScheduler struct {
	db           *sql.DB
	paymentRepo  *repositories.PaymentRepository
	invoiceRepo  *repositories.InvoiceRepository
	stopChan     chan struct{}
}

func NewJobScheduler(db *sql.DB) *JobScheduler {
	return &JobScheduler{
		db:          db,
		paymentRepo: repositories.NewPaymentRepository(db),
		invoiceRepo: repositories.NewInvoiceRepository(db),
		stopChan:    make(chan struct{}),
	}
}

func (js *JobScheduler) Start() {
	log.Println("Starting background job scheduler...")

	go js.runPaymentStatusCheck()
	go js.runOverdueInvoiceCheck()
}

func (js *JobScheduler) Stop() {
	close(js.stopChan)
	log.Println("Background job scheduler stopped")
}

func (js *JobScheduler) runPaymentStatusCheck() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			js.checkPendingPayments()
		case <-js.stopChan:
			return
		}
	}
}

func (js *JobScheduler) checkPendingPayments() {
	log.Println("Checking pending payments...")

	query := `
		SELECT id, expires_at
		FROM payments
		WHERE status = 'Pending' AND expires_at < NOW()
	`

	rows, err := js.db.Query(query)
	if err != nil {
		log.Printf("Error checking pending payments: %v", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id uint
		var expiresAt time.Time
		if err := rows.Scan(&id, &expiresAt); err != nil {
			continue
		}

		updateQuery := `UPDATE payments SET status = 'Expired', status_updated_at = NOW() WHERE id = ?`
		if _, err := js.db.Exec(updateQuery, id); err != nil {
			log.Printf("Error expiring payment %d: %v", id, err)
			continue
		}
		count++
	}

	if count > 0 {
		log.Printf("Expired %d pending payments", count)
	}
}

func (js *JobScheduler) runOverdueInvoiceCheck() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			js.checkOverdueInvoices()
		case <-js.stopChan:
			return
		}
	}
}

func (js *JobScheduler) checkOverdueInvoices() {
	log.Println("Checking overdue invoices...")

	query := `
		UPDATE invoices
		SET status = 'Overdue'
		WHERE status IN ('Issued', 'Partially Paid')
		  AND due_date < CURDATE()
		  AND total_amount > paid_amount
	`

	result, err := js.db.Exec(query)
	if err != nil {
		log.Printf("Error checking overdue invoices: %v", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		log.Printf("Marked %d invoices as overdue", rowsAffected)
	}
}
