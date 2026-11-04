package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/kwizera/university-invoices/shared/repositories"
)

type WebhookController struct {
	paymentRepo *repositories.PaymentRepository
	invoiceRepo *repositories.InvoiceRepository
	receiptRepo *repositories.ReceiptRepository
	auditRepo   *repositories.AuditRepository
}

func NewWebhookController(db *sql.DB) *WebhookController {
	return &WebhookController{
		paymentRepo: repositories.NewPaymentRepository(db),
		invoiceRepo: repositories.NewInvoiceRepository(db),
		receiptRepo: repositories.NewReceiptRepository(db),
		auditRepo:   repositories.NewAuditRepository(db),
	}
}

func (c *WebhookController) HandleMoPayWebhook(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		TxRef  string `json:"tx_ref"`
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"code":"VALIDATION_ERROR","message":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	tx, err := c.paymentRepo.GetDB().Begin()
	if err != nil {
		http.Error(w, `{"code":"INTERNAL_ERROR","message":"Database error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	payment, err := c.paymentRepo.GetByTransactionIDInTx(tx, payload.TxRef)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	currentStatus := payment["status"].(string)
	var newStatus string
	if payload.Status == "success" || payload.Status == "paid" {
		newStatus = "Paid"
	} else if payload.Status == "failed" || payload.Status == "cancelled" {
		newStatus = "Failed"
	} else {
		newStatus = payload.Status
	}

	if currentStatus == newStatus {
		w.WriteHeader(http.StatusOK)
		return
	}

	var paidAt *time.Time
	if newStatus == "Paid" {
		now := time.Now()
		paidAt = &now
	}

	err = c.paymentRepo.UpdateStatusInTx(tx, payment["id"].(uint), newStatus, paidAt)
	if err != nil {
		http.Error(w, `{"code":"INTERNAL_ERROR","message":"Failed to update payment"}`, http.StatusInternalServerError)
		return
	}

	if newStatus == "Paid" {
		amount := payment["amount"].(int64)
		invoiceID := payment["invoice_id"].(uint)

		err = c.invoiceRepo.UpdateBalanceInTx(tx, invoiceID, amount)
		if err != nil {
			http.Error(w, `{"code":"INTERNAL_ERROR","message":"Failed to update invoice"}`, http.StatusInternalServerError)
			return
		}

		invoice, _ := c.invoiceRepo.GetByIDInTx(tx, invoiceID)
		balance := invoice["balance"].(int64)

		_, err := c.receiptRepo.CreateInTx(tx, payment["id"].(uint), invoiceID, amount, balance)
		if err != nil {
			http.Error(w, `{"code":"INTERNAL_ERROR","message":"Failed to create receipt"}`, http.StatusInternalServerError)
			return
		}

		c.auditRepo.LogInTx(tx, "gateway", nil, "Payment status updated", "payment", payment["mopay_reference"].(string), "Pending", "Paid")
		c.auditRepo.LogInTx(tx, "system", nil, "Receipt generated", "receipt", "", "", "")
		c.auditRepo.LogInTx(tx, "system", nil, "Invoice balance updated", "invoice", invoice["invoice_number"].(string), "", "")
	} else {
		c.auditRepo.LogInTx(tx, "gateway", nil, "Payment status updated", "payment", payment["mopay_reference"].(string), "Pending", newStatus)
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, `{"code":"INTERNAL_ERROR","message":"Failed to commit transaction"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
