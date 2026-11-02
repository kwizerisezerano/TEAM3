package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	invoicemodels "github.com/kwizera/university-invoices/services/invoice/models"
	"github.com/kwizera/university-invoices/services/payment/models"
	"github.com/kwizera/university-invoices/shared/repositories"
)

type PaymentController struct {
	paymentRepo *repositories.PaymentRepository
	invoiceRepo *repositories.InvoiceRepository
	receiptRepo *repositories.ReceiptRepository
}

func NewPaymentController(
	paymentRepo *repositories.PaymentRepository,
	invoiceRepo *repositories.InvoiceRepository,
	receiptRepo *repositories.ReceiptRepository,
) *PaymentController {
	return &PaymentController{
		paymentRepo: paymentRepo,
		invoiceRepo: invoiceRepo,
		receiptRepo: receiptRepo,
	}
}

func (c *PaymentController) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify invoice exists
	invoice, err := c.invoiceRepo.GetByID(req.InvoiceID)
	if err != nil {
		http.Error(w, "Invoice not found", http.StatusNotFound)
		return
	}

	// Verify student matches
	if invoice.StudentID != req.StudentID {
		http.Error(w, "Student ID does not match invoice", http.StatusBadRequest)
		return
	}

	// Check if payment amount exceeds outstanding balance
	if req.Amount > invoice.OutstandingBalance {
		http.Error(w, "Payment amount exceeds outstanding balance", http.StatusBadRequest)
		return
	}

	// Generate payment number
	paymentNumber := fmt.Sprintf("PAY-%d-%s", req.InvoiceID, time.Now().Format("20060102150405"))

	// TODO: Get user ID from JWT token
	recordedBy := 1

	payment := &models.Payment{
		PaymentNumber:   paymentNumber,
		InvoiceID:       req.InvoiceID,
		StudentID:       req.StudentID,
		Amount:          req.Amount,
		PaymentMethod:   req.PaymentMethod,
		PaymentDate:     req.PaymentDate,
		ReferenceNumber: req.ReferenceNumber,
		Status:          models.PaymentStatusCompleted,
		Notes:           req.Notes,
		RecordedBy:      recordedBy,
	}

	if err := c.paymentRepo.Create(payment); err != nil {
		http.Error(w, "Failed to create payment", http.StatusInternalServerError)
		return
	}

	// Update invoice payment amounts
	newPaidAmount := invoice.PaidAmount + req.Amount
	newOutstandingBalance := invoice.OutstandingBalance - req.Amount

	if err := c.invoiceRepo.UpdatePaymentAmounts(req.InvoiceID, newPaidAmount, newOutstandingBalance); err != nil {
		http.Error(w, "Failed to update invoice", http.StatusInternalServerError)
		return
	}

	// Update invoice status based on balance
	if newOutstandingBalance <= 0 {
		c.invoiceRepo.UpdateInvoiceStatus(req.InvoiceID, invoicemodels.InvoiceStatusPaid)
	} else if newPaidAmount > 0 {
		c.invoiceRepo.UpdateInvoiceStatus(req.InvoiceID, invoicemodels.InvoiceStatusPartial)
	}

	// Generate receipt
	receiptNumber := fmt.Sprintf("RCT-%d-%s", payment.ID, time.Now().Format("20060102"))
	receipt := &models.Receipt{
		ReceiptNumber: receiptNumber,
		PaymentID:     payment.ID,
		InvoiceID:     req.InvoiceID,
		StudentID:     req.StudentID,
		Amount:        req.Amount,
		IssuedDate:    time.Now(),
		IssuedBy:      recordedBy,
	}

	if err := c.receiptRepo.Create(receipt); err != nil {
		// Log error but don't fail the payment
		fmt.Printf("Warning: failed to create receipt: %v\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payment)
}

func (c *PaymentController) GetPayment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}

	payment, err := c.paymentRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Payment not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payment)
}

func (c *PaymentController) GetPaymentsByInvoice(w http.ResponseWriter, r *http.Request) {
	invoiceID, err := strconv.Atoi(r.URL.Query().Get("invoice_id"))
	if err != nil {
		http.Error(w, "Invalid invoice ID", http.StatusBadRequest)
		return
	}

	payments, err := c.paymentRepo.GetByInvoiceID(invoiceID)
	if err != nil {
		http.Error(w, "Failed to get payments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payments)
}

func (c *PaymentController) GetPaymentsByStudent(w http.ResponseWriter, r *http.Request) {
	studentID, err := strconv.Atoi(r.URL.Query().Get("student_id"))
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit == 0 {
		limit = 50
	}

	payments, err := c.paymentRepo.GetByStudentID(studentID, limit, offset)
	if err != nil {
		http.Error(w, "Failed to get payments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payments)
}

func (c *PaymentController) ListPayments(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit == 0 {
		limit = 50
	}

	payments, err := c.paymentRepo.List(limit, offset)
	if err != nil {
		http.Error(w, "Failed to list payments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payments)
}

func (c *PaymentController) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}

	var req models.UpdatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	payment, err := c.paymentRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Payment not found", http.StatusNotFound)
		return
	}

	if req.Status != nil {
		payment.Status = *req.Status
	}
	if req.ReferenceNumber != nil {
		payment.ReferenceNumber = *req.ReferenceNumber
	}
	if req.Notes != nil {
		payment.Notes = *req.Notes
	}

	if err := c.paymentRepo.Update(payment); err != nil {
		http.Error(w, "Failed to update payment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payment)
}

func (c *PaymentController) DeletePayment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}

	if err := c.paymentRepo.Delete(id); err != nil {
		http.Error(w, "Failed to delete payment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
