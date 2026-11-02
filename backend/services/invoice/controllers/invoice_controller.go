package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"github.com/kwizera/university-invoices/services/invoice/models"
	"github.com/kwizera/university-invoices/shared/repositories"
)

type InvoiceController struct {
	invoiceRepo     *repositories.InvoiceRepository
	invoiceItemRepo *repositories.InvoiceItemRepository
	studentRepo     *repositories.StudentRepository
}

func NewInvoiceController(
	invoiceRepo *repositories.InvoiceRepository,
	invoiceItemRepo *repositories.InvoiceItemRepository,
	studentRepo *repositories.StudentRepository,
) *InvoiceController {
	return &InvoiceController{
		invoiceRepo:     invoiceRepo,
		invoiceItemRepo: invoiceItemRepo,
		studentRepo:     studentRepo,
	}
}

func (c *InvoiceController) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	var req models.CreateInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify student exists
	_, err := c.studentRepo.GetByID(req.StudentID)
	if err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	// Calculate totals
	var subtotal float64
	for _, item := range req.Items {
		subtotal += item.UnitPrice * float64(item.Quantity)
	}

	taxAmount := subtotal * (req.TaxRate / 100)
	totalAmount := subtotal + taxAmount - req.DiscountAmount

	// Generate invoice number
	invoiceNumber := fmt.Sprintf("INV-%d-%s", req.StudentID, time.Now().Format("20060102"))

	// TODO: Get user ID from JWT token
	createdBy := 1

	invoice := &models.Invoice{
		InvoiceNumber:       invoiceNumber,
		StudentID:            req.StudentID,
		IssueDate:            req.IssueDate,
		DueDate:              req.DueDate,
		Subtotal:             subtotal,
		TaxAmount:            taxAmount,
		DiscountAmount:       req.DiscountAmount,
		TotalAmount:          totalAmount,
		PaidAmount:           0,
		OutstandingBalance:   totalAmount,
		Status:               models.InvoiceStatusPending,
		Notes:                req.Notes,
		CreatedBy:            createdBy,
	}

	if err := c.invoiceRepo.Create(invoice); err != nil {
		http.Error(w, "Failed to create invoice", http.StatusInternalServerError)
		return
	}

	// Create invoice items
	var items []models.InvoiceItem
	for _, item := range req.Items {
		invoiceItem := &models.InvoiceItem{
			InvoiceID:   invoice.ID,
			FeeTypeID:   item.FeeTypeID,
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Amount:      item.UnitPrice * float64(item.Quantity),
		}

		if err := c.invoiceItemRepo.Create(invoiceItem); err != nil {
			http.Error(w, "Failed to create invoice items", http.StatusInternalServerError)
			return
		}

		items = append(items, *invoiceItem)
	}

	response := models.InvoiceWithItems{
		Invoice: *invoice,
		Items:   items,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (c *InvoiceController) GetInvoice(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid invoice ID", http.StatusBadRequest)
		return
	}

	invoice, err := c.invoiceRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Invoice not found", http.StatusNotFound)
		return
	}

	items, err := c.invoiceItemRepo.GetByInvoiceID(id)
	if err != nil {
		http.Error(w, "Failed to get invoice items", http.StatusInternalServerError)
		return
	}

	response := models.InvoiceWithItems{
		Invoice: *invoice,
		Items:   items,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *InvoiceController) GetInvoicesByStudent(w http.ResponseWriter, r *http.Request) {
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

	invoices, err := c.invoiceRepo.GetByStudentID(studentID, limit, offset)
	if err != nil {
		http.Error(w, "Failed to get invoices", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoices)
}

func (c *InvoiceController) ListInvoices(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit == 0 {
		limit = 50
	}

	invoices, err := c.invoiceRepo.List(limit, offset)
	if err != nil {
		http.Error(w, "Failed to list invoices", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoices)
}

func (c *InvoiceController) UpdateInvoice(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid invoice ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	invoice, err := c.invoiceRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Invoice not found", http.StatusNotFound)
		return
	}

	if req.DueDate != nil {
		invoice.DueDate = *req.DueDate
	}
	if req.DiscountAmount != nil {
		invoice.DiscountAmount = *req.DiscountAmount
	}
	if req.Notes != nil {
		invoice.Notes = *req.Notes
	}
	if req.Status != nil {
		invoice.Status = *req.Status
	}

	if err := c.invoiceRepo.Update(invoice); err != nil {
		http.Error(w, "Failed to update invoice", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoice)
}

func (c *InvoiceController) DeleteInvoice(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid invoice ID", http.StatusBadRequest)
		return
	}

	// Delete invoice items first
	if err := c.invoiceItemRepo.DeleteByInvoiceID(id); err != nil {
		http.Error(w, "Failed to delete invoice items", http.StatusInternalServerError)
		return
	}

	if err := c.invoiceRepo.Delete(id); err != nil {
		http.Error(w, "Failed to delete invoice", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
