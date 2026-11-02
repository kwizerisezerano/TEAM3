package models

import "time"

type InvoiceStatus string

const (
	InvoiceStatusDraft     InvoiceStatus = "draft"
	InvoiceStatusPending   InvoiceStatus = "pending"
	InvoiceStatusPartial   InvoiceStatus = "partial"
	InvoiceStatusPaid      InvoiceStatus = "paid"
	InvoiceStatusOverdue   InvoiceStatus = "overdue"
	InvoiceStatusCancelled InvoiceStatus = "cancelled"
)

type Invoice struct {
	ID                int           `json:"id"`
	InvoiceNumber     string        `json:"invoice_number"`
	StudentID         int           `json:"student_id"`
	IssueDate         time.Time     `json:"issue_date"`
	DueDate           time.Time     `json:"due_date"`
	Subtotal          float64       `json:"subtotal"`
	TaxAmount         float64       `json:"tax_amount"`
	DiscountAmount    float64       `json:"discount_amount"`
	TotalAmount       float64       `json:"total_amount"`
	PaidAmount        float64       `json:"paid_amount"`
	OutstandingBalance float64      `json:"outstanding_balance"`
	Status            InvoiceStatus `json:"status"`
	Notes             string        `json:"notes"`
	CreatedBy         int           `json:"created_by"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
}

type InvoiceItem struct {
	ID          int       `json:"id"`
	InvoiceID   int       `json:"invoice_id"`
	FeeTypeID   int       `json:"fee_type_id"`
	Description string    `json:"description"`
	Quantity    int       `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	Amount      float64   `json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateInvoiceRequest struct {
	StudentID      int                 `json:"student_id" binding:"required"`
	IssueDate      time.Time           `json:"issue_date" binding:"required"`
	DueDate        time.Time           `json:"due_date" binding:"required"`
	Items          []CreateInvoiceItem `json:"items" binding:"required,min=1"`
	TaxRate        float64             `json:"tax_rate"`
	DiscountAmount float64             `json:"discount_amount"`
	Notes          string              `json:"notes"`
}

type CreateInvoiceItem struct {
	FeeTypeID   int     `json:"fee_type_id" binding:"required"`
	Description string `json:"description" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required,min=1"`
	UnitPrice   float64 `json:"unit_price" binding:"required,gt=0"`
}

type UpdateInvoiceRequest struct {
	DueDate        *time.Time    `json:"due_date"`
	DiscountAmount *float64      `json:"discount_amount"`
	Notes          *string       `json:"notes"`
	Status         *InvoiceStatus `json:"status"`
}

type InvoiceWithItems struct {
	Invoice
	Items []InvoiceItem `json:"items"`
}
