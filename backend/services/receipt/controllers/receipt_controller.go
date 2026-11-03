package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/signintech/gopdf"

	"github.com/kwizera/university-invoices/shared/repositories"
)

type ReceiptController struct {
	receiptRepo *repositories.ReceiptRepository
}

func NewReceiptController(db *sql.DB) *ReceiptController {
	return &ReceiptController{
		receiptRepo: repositories.NewReceiptRepository(db),
	}
}

func (c *ReceiptController) GetReceipt(w http.ResponseWriter, r *http.Request) {
	receiptID, _ := strconv.ParseUint(chi.URLParam(r, "receiptId"), 10, 32)

	receipt, err := c.receiptRepo.GetByID(uint(receiptID))
	if err != nil {
		http.Error(w, `{"code":"NOT_FOUND","message":"Receipt not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(receipt)
}

func (c *ReceiptController) GetReceiptPDF(w http.ResponseWriter, r *http.Request) {
	receiptID, _ := strconv.ParseUint(chi.URLParam(r, "receiptId"), 10, 32)

	receipt, err := c.receiptRepo.GetByID(uint(receiptID))
	if err != nil {
		http.Error(w, `{"code":"NOT_FOUND","message":"Receipt not found"}`, http.StatusNotFound)
		return
	}

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	err = pdf.AddTTFFont("arial", "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf")
	if err != nil {
		pdf.AddTTFFont("arial", "/System/Library/Fonts/Helvetica.ttc")
	}

	pdf.SetFont("arial", "", 14)
	pdf.SetX(50)
	pdf.SetY(50)
	pdf.Cell(nil, "UNIVERSITY OF KIGALI")
	pdf.SetY(70)
	pdf.Cell(nil, "RECEIPT")
	pdf.SetY(100)
	pdf.Cell(nil, "Receipt Number: "+receipt["receipt_number"].(string))
	pdf.SetY(120)
	pdf.Cell(nil, "Invoice Number: "+receipt["invoice_number"].(string))
	pdf.SetY(140)
	pdf.Cell(nil, "Student: "+receipt["student_name"].(string))
	pdf.SetY(160)
	pdf.Cell(nil, "Amount Paid: "+strconv.FormatInt(receipt["amount_paid"].(int64), 10)+" RWF")
	pdf.SetY(180)
	pdf.Cell(nil, "New Balance: "+strconv.FormatInt(receipt["new_balance"].(int64), 10)+" RWF")

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=receipt_"+receipt["receipt_number"].(string)+".pdf")
	pdf.Write(w)
}
