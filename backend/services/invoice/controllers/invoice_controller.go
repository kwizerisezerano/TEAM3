package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/kwizera/university-invoices/services/invoice/models"
	"github.com/kwizera/university-invoices/shared/repositories"
)

type InvoiceController struct {
	invoiceRepo     *repositories.InvoiceRepository
	invoiceItemRepo *repositories.InvoiceItemRepository
	studentRepo     *repositories.StudentRepository
	validator       *validator.Validate
}

func NewInvoiceController(db *sql.DB) *InvoiceController {
	return &InvoiceController{
		invoiceRepo:     repositories.NewInvoiceRepository(db),
		invoiceItemRepo: repositories.NewInvoiceItemRepository(db),
		studentRepo:     repositories.NewStudentRepository(db),
		validator:       validator.New(),
	}
}

func (c *InvoiceController) ListInvoices(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	role := r.Context().Value("userRole").(string)

	var studentID uint
	if role == "student" {
		sid, _, _, err := c.studentRepo.GetStudentByUserID(userID)
		if err != nil {
			http.Error(w, `{"code":"NOT_FOUND","message":"Student not found"}`, http.StatusNotFound)
			return
		}
		studentID = sid
	}

	status := r.URL.Query().Get("status")
	semesterID := r.URL.Query().Get("semester_id")
	programmeID := r.URL.Query().Get("programme_id")
	queryStr := r.URL.Query().Get("q")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize == 0 {
		pageSize = 20
	}

	invoices, total, err := c.invoiceRepo.List(studentID, status, semesterID, programmeID, queryStr, from, to, page, pageSize)
	if err != nil {
		http.Error(w, `{"code":"INTERNAL_ERROR","message":"Failed to fetch invoices"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":      invoices,
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *InvoiceController) GetInvoice(w http.ResponseWriter, r *http.Request) {
	invoiceID, _ := strconv.ParseUint(chi.URLParam(r, "invoiceId"), 10, 32)
	userID := r.Context().Value("userID").(uint)
	role := r.Context().Value("userRole").(string)

	invoice, err := c.invoiceRepo.GetByID(uint(invoiceID))
	if err != nil {
		http.Error(w, `{"code":"NOT_FOUND","message":"Invoice not found"}`, http.StatusNotFound)
		return
	}

	if role == "student" {
		sid, _, _, _ := c.studentRepo.GetStudentByUserID(userID)
		if sid != invoice["student_id"].(uint) {
			http.Error(w, `{"code":"FORBIDDEN","message":"You can only view your own invoices"}`, http.StatusForbidden)
			return
		}
	}

	items, err := c.invoiceItemRepo.GetByInvoiceID(uint(invoiceID))
	if err != nil {
		http.Error(w, `{"code":"INTERNAL_ERROR","message":"Failed to fetch invoice items"}`, http.StatusInternalServerError)
		return
	}

	invoice["items"] = items

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoice)
}

func (c *InvoiceController) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)

	var req models.InvoiceInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"code":"VALIDATION_ERROR","message":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if err := c.validator.Struct(&req); err != nil {
		http.Error(w, `{"code":"VALIDATION_ERROR","message":"One or more fields are invalid"}`, http.StatusBadRequest)
		return
	}

	totalAmount := int64(0)
	for _, item := range req.Items {
		totalAmount += item.Amount
	}

	invoiceID, err := c.invoiceRepo.Create(req.StudentID, req.ProgrammeID, req.SemesterID, userID, req.DueDate, totalAmount)
	if err != nil {
		http.Error(w, `{"code":"INTERNAL_ERROR","message":"Failed to create invoice"}`, http.StatusInternalServerError)
		return
	}

	for _, item := range req.Items {
		c.invoiceItemRepo.Create(invoiceID, item.FeeTypeID, "Fee", item.Amount)
	}

	if req.Action == "issue" {
		c.invoiceRepo.Issue(invoiceID)
	}

	invoice, _ := c.invoiceRepo.GetByID(invoiceID)
	items, _ := c.invoiceItemRepo.GetByInvoiceID(invoiceID)
	invoice["items"] = items

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(invoice)
}

func (c *InvoiceController) UpdateInvoice(w http.ResponseWriter, r *http.Request) {
	invoiceID, _ := strconv.ParseUint(chi.URLParam(r, "invoiceId"), 10, 32)

	var req models.InvoiceInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"code":"VALIDATION_ERROR","message":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if err := c.validator.Struct(&req); err != nil {
		http.Error(w, `{"code":"VALIDATION_ERROR","message":"One or more fields are invalid"}`, http.StatusBadRequest)
		return
	}

	totalAmount := int64(0)
	for _, item := range req.Items {
		totalAmount += item.Amount
	}

	err := c.invoiceRepo.Update(uint(invoiceID), req.DueDate, totalAmount)
	if err != nil {
		if err.Error() == "invoice not found or not in draft status" {
			http.Error(w, `{"code":"CONFLICT","message":"Invoice is not a draft"}`, http.StatusConflict)
			return
		}
		http.Error(w, `{"code":"INTERNAL_ERROR","message":"Failed to update invoice"}`, http.StatusInternalServerError)
		return
	}

	c.invoiceItemRepo.DeleteByInvoiceID(uint(invoiceID))
	for _, item := range req.Items {
		c.invoiceItemRepo.Create(uint(invoiceID), item.FeeTypeID, "Fee", item.Amount)
	}

	invoice, _ := c.invoiceRepo.GetByID(uint(invoiceID))
	items, _ := c.invoiceItemRepo.GetByInvoiceID(uint(invoiceID))
	invoice["items"] = items

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoice)
}

func (c *InvoiceController) IssueInvoice(w http.ResponseWriter, r *http.Request) {
	invoiceID, _ := strconv.ParseUint(chi.URLParam(r, "invoiceId"), 10, 32)

	err := c.invoiceRepo.Issue(uint(invoiceID))
	if err != nil {
		http.Error(w, `{"code":"INTERNAL_ERROR","message":"Failed to issue invoice"}`, http.StatusInternalServerError)
		return
	}

	invoice, _ := c.invoiceRepo.GetByID(uint(invoiceID))
	items, _ := c.invoiceItemRepo.GetByInvoiceID(uint(invoiceID))
	invoice["items"] = items

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoice)
}
