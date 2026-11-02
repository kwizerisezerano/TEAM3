package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/kwizera/university-invoices/services/fee/models"
	"github.com/kwizera/university-invoices/shared/repositories"
)

type FeeController struct {
	feeRepo *repositories.FeeRepository
}

func NewFeeController(feeRepo *repositories.FeeRepository) *FeeController {
	return &FeeController{feeRepo: feeRepo}
}

func (c *FeeController) CreateFeeType(w http.ResponseWriter, r *http.Request) {
	var req models.CreateFeeTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	feeType := &models.FeeType{
		Name:        req.Name,
		Description: req.Description,
		Amount:      req.Amount,
		Currency:    req.Currency,
		IsRecurring: req.IsRecurring,
		IsActive:    true,
	}

	if err := c.feeRepo.Create(feeType); err != nil {
		http.Error(w, "Failed to create fee type", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(feeType)
}

func (c *FeeController) GetFeeType(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid fee type ID", http.StatusBadRequest)
		return
	}

	feeType, err := c.feeRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Fee type not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feeType)
}

func (c *FeeController) ListFeeTypes(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit == 0 {
		limit = 50
	}

	feeTypes, err := c.feeRepo.List(limit, offset)
	if err != nil {
		http.Error(w, "Failed to list fee types", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feeTypes)
}

func (c *FeeController) UpdateFeeType(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid fee type ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateFeeTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	feeType, err := c.feeRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Fee type not found", http.StatusNotFound)
		return
	}

	if req.Name != nil {
		feeType.Name = *req.Name
	}
	if req.Description != nil {
		feeType.Description = *req.Description
	}
	if req.Amount != nil {
		feeType.Amount = *req.Amount
	}
	if req.Currency != nil {
		feeType.Currency = *req.Currency
	}
	if req.IsRecurring != nil {
		feeType.IsRecurring = *req.IsRecurring
	}
	if req.IsActive != nil {
		feeType.IsActive = *req.IsActive
	}

	if err := c.feeRepo.Update(feeType); err != nil {
		http.Error(w, "Failed to update fee type", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feeType)
}

func (c *FeeController) DeleteFeeType(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid fee type ID", http.StatusBadRequest)
		return
	}

	if err := c.feeRepo.Delete(id); err != nil {
		http.Error(w, "Failed to delete fee type", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
