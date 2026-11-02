package models

import "time"

type PaymentMethod string

const (
	PaymentMethodCash         PaymentMethod = "cash"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodMobileMoney  PaymentMethod = "mobile_money"
	PaymentMethodCard         PaymentMethod = "card"
	PaymentMethodCheck        PaymentMethod = "check"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

type Payment struct {
	ID              int           `json:"id"`
	PaymentNumber   string        `json:"payment_number"`
	InvoiceID       int           `json:"invoice_id"`
	StudentID       int           `json:"student_id"`
	Amount          float64       `json:"amount"`
	PaymentMethod   PaymentMethod `json:"payment_method"`
	PaymentDate     time.Time     `json:"payment_date"`
	ReferenceNumber string        `json:"reference_number"`
	Status          PaymentStatus `json:"status"`
	Notes           string        `json:"notes"`
	RecordedBy      int           `json:"recorded_by"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

type CreatePaymentRequest struct {
	InvoiceID       int           `json:"invoice_id" binding:"required"`
	StudentID       int           `json:"student_id" binding:"required"`
	Amount          float64       `json:"amount" binding:"required,gt=0"`
	PaymentMethod   PaymentMethod `json:"payment_method" binding:"required"`
	PaymentDate     time.Time     `json:"payment_date" binding:"required"`
	ReferenceNumber string        `json:"reference_number"`
	Notes           string        `json:"notes"`
}

type UpdatePaymentRequest struct {
	Status          *PaymentStatus `json:"status"`
	ReferenceNumber *string        `json:"reference_number"`
	Notes           *string        `json:"notes"`
}
