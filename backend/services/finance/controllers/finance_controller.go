package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kwizera/university-invoices/shared/repositories"
)

type FinanceController struct {
	invoiceRepo   *repositories.InvoiceRepository
	paymentRepo   *repositories.PaymentRepository
	studentRepo   *repositories.StudentRepository
	programmeRepo *repositories.ProgrammeRepository
	semesterRepo  *repositories.SemesterRepository
}

func NewFinanceController(db *sql.DB) *FinanceController {
	return &FinanceController{
		invoiceRepo:   repositories.NewInvoiceRepository(db),
		paymentRepo:   repositories.NewPaymentRepository(db),
		studentRepo:   repositories.NewStudentRepository(db),
		programmeRepo: repositories.NewProgrammeRepository(db),
		semesterRepo:  repositories.NewSemesterRepository(db),
	}
}

func (c *FinanceController) GetDashboard(w http.ResponseWriter, r *http.Request) {
	invoices, _, _ := c.invoiceRepo.List(0, "", "", "", "", "", "", 1, 1000)
	payments, _, _ := c.paymentRepo.List(0, "", 1, 1000)

	totalInvoiced := int64(0)
	totalCollected := int64(0)
	overdueAmount := int64(0)
	pendingPayments := 0
	failedPayments := 0

	for _, inv := range invoices {
		totalInvoiced += inv["total_amount"].(int64)
		if inv["status"].(string) == "Overdue" {
			overdueAmount += inv["balance"].(int64)
		}
	}

	for _, pay := range payments {
		if pay["status"].(string) == "Paid" {
			totalCollected += pay["amount"].(int64)
		} else if pay["status"].(string) == "Pending" {
			pendingPayments++
		} else if pay["status"].(string) == "Failed" {
			failedPayments++
		}
	}

	outstandingBalance := totalInvoiced - totalCollected

	dashboard := map[string]interface{}{
		"total_invoiced":      totalInvoiced,
		"total_collected":     totalCollected,
		"outstanding_balance": outstandingBalance,
		"overdue_amount":      overdueAmount,
		"pending_payments":    pendingPayments,
		"failed_payments":     failedPayments,
		"total_students":      0,
		"active_invoices":     len(invoices),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dashboard)
}

func (c *FinanceController) GetReportSummary(w http.ResponseWriter, r *http.Request) {
	invoices, _, _ := c.invoiceRepo.List(0, "", "", "", "", "", "", 1, 10000)
	payments, _, _ := c.paymentRepo.List(0, "", 1, 10000)

	totalInvoiced := int64(0)
	totalCollected := int64(0)

	for _, inv := range invoices {
		totalInvoiced += inv["total_amount"].(int64)
	}

	for _, pay := range payments {
		if pay["status"].(string) == "Paid" {
			totalCollected += pay["amount"].(int64)
		}
	}

	outstandingBalance := totalInvoiced - totalCollected
	collectionRate := float64(0)
	if totalInvoiced > 0 {
		collectionRate = float64(totalCollected) / float64(totalInvoiced) * 100
	}

	summary := map[string]interface{}{
		"total_invoiced":      totalInvoiced,
		"total_collected":     totalCollected,
		"outstanding_balance": outstandingBalance,
		"collection_rate":     collectionRate,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

func (c *FinanceController) GetAccountsReport(w http.ResponseWriter, r *http.Request) {
	invoices, _, _ := c.invoiceRepo.List(0, "", "", "", "", "", "", 1, 10000)

	accounts := make(map[uint]map[string]interface{})
	for _, inv := range invoices {
		student := inv["student"].(map[string]interface{})
		studentID := student["id"].(uint)

		if _, exists := accounts[studentID]; !exists {
			accounts[studentID] = map[string]interface{}{
				"student_id":          studentID,
				"student_number":      student["student_number"],
				"student_name":        student["full_name"],
				"programme":           "",
				"total_invoiced":      int64(0),
				"total_paid":          int64(0),
				"outstanding_balance": int64(0),
				"invoice_count":       0,
			}
		}

		accounts[studentID]["total_invoiced"] = accounts[studentID]["total_invoiced"].(int64) + inv["total_amount"].(int64)
		accounts[studentID]["total_paid"] = accounts[studentID]["total_paid"].(int64) + inv["paid_amount"].(int64)
		accounts[studentID]["outstanding_balance"] = accounts[studentID]["outstanding_balance"].(int64) + inv["balance"].(int64)
		accounts[studentID]["invoice_count"] = accounts[studentID]["invoice_count"].(int) + 1
	}

	var result []map[string]interface{}
	for _, acc := range accounts {
		result = append(result, acc)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (c *FinanceController) SearchStudents(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 20
	}

	students, err := c.studentRepo.SearchStudents(query, limit)
	if err != nil {
		http.Error(w, `{"code":"INTERNAL_ERROR","message":"Failed to search students"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func (c *FinanceController) ListProgrammes(w http.ResponseWriter, r *http.Request) {
	programmes, err := c.programmeRepo.ListAll()
	if err != nil {
		http.Error(w, `{"code":"INTERNAL_ERROR","message":"Failed to fetch programmes"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(programmes)
}

func (c *FinanceController) ListSemesters(w http.ResponseWriter, r *http.Request) {
	semesters, err := c.semesterRepo.ListAll()
	if err != nil {
		http.Error(w, `{"code":"INTERNAL_ERROR","message":"Failed to fetch semesters"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(semesters)
}
