package models

import "time"

type FeeType struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	IsRecurring bool      `json:"is_recurring"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateFeeTypeRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Currency    string  `json:"currency" binding:"required"`
	IsRecurring bool    `json:"is_recurring"`
}

type UpdateFeeTypeRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Amount      *float64 `json:"amount"`
	Currency    *string  `json:"currency"`
	IsRecurring *bool    `json:"is_recurring"`
	IsActive    *bool    `json:"is_active"`
}
