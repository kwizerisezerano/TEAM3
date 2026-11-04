package jobs

import (
	"database/sql"
	"log"
	"time"

	"github.com/kwizera/university-invoices/internal/config"
	"github.com/kwizera/university-invoices/services/mopay"
	"github.com/kwizera/university-invoices/shared/repositories"
)

type JobScheduler struct {
	db           *sql.DB
	paymentRepo  *repositories.PaymentRepository
	invoiceRepo  *repositories.InvoiceRepository
	mopayService *mopay.Service
	stopChan     chan struct{}
}

func NewJobScheduler(db *sql.DB, cfg *config.Config) *JobScheduler {
	return &JobScheduler{
		db:           db,
		paymentRepo:  repositories.NewPaymentRepository(db),
		invoiceRepo:  repositories.NewInvoiceRepository(db),
		mopayService: mopay.NewService(cfg.MoPayAPIURL, cfg.MoPayAPIToken, cfg.MoPayWalletPhone, cfg.PaymentCallbackURL, cfg.PaymentTimeoutMS),
		stopChan:     make(chan struct{}),
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
		SELECT id, mopay_transaction_id
		FROM payments
		WHERE status = 'Pending'
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
		var mopayTransactionID string
		if err := rows.Scan(&id, &mopayTransactionID); err != nil {
			continue
		}

		// Check status with MoPay API
		statusResp, err := js.mopayService.CheckTransactionStatus(mopayTransactionID)
		if err != nil {
			log.Printf("Error checking status for payment %d: %v", id, err)
			continue
		}

		// Update payment status based on MoPay response
		var newStatus string
		if statusResp.Status == 200 {
			newStatus = "Paid"
		} else if statusResp.Status == 400 {
			newStatus = "Failed"
		} else {
			// Still pending (201, 202, etc.)
			continue
		}

		updateQuery := `UPDATE payments SET status = ?, status_updated_at = NOW() WHERE id = ?`
		if _, err := js.db.Exec(updateQuery, newStatus, id); err != nil {
			log.Printf("Error updating payment %d status: %v", id, err)
			continue
		}

		// If paid, update invoice balance
		if newStatus == "Paid" {
			js.updateInvoiceBalance(id)
		}

		count++
		log.Printf("Updated payment %d to %s", id, newStatus)
	}

	if count > 0 {
		log.Printf("Updated %d payment statuses", count)
	}
}

func (js *JobScheduler) updateInvoiceBalance(paymentID uint) {
	// Get payment details
	var invoiceID uint
	var amount int64
	query := `SELECT invoice_id, amount FROM payments WHERE id = ?`
	err := js.db.QueryRow(query, paymentID).Scan(&invoiceID, &amount)
	if err != nil {
		log.Printf("Error getting payment details: %v", err)
		return
	}

	// Update invoice paid amount
	updateQuery := `
		UPDATE invoices
		SET paid_amount = paid_amount + ?,
		    status = CASE
		        WHEN paid_amount + ? >= total_amount THEN 'Paid'
		        WHEN paid_amount + ? > 0 THEN 'Partially Paid'
		        ELSE status
		    END
		WHERE id = ?
	`
	_, err = js.db.Exec(updateQuery, amount, amount, amount, invoiceID)
	if err != nil {
		log.Printf("Error updating invoice balance: %v", err)
		return
	}

	log.Printf("Updated invoice %d balance by %d", invoiceID, amount)
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
