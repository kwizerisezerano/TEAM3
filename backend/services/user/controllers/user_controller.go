package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kwizera/university-invoices/services/user/models"
	"github.com/kwizera/university-invoices/shared/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	userRepo *repositories.UserRepository
}

func NewUserController(userRepo *repositories.UserRepository) *UserController {
	return &UserController{userRepo: userRepo}
}

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	existingUser, _ := c.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		http.Error(w, "User with this email already exists", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Encrypt email
	encryptedEmail, err := EncryptEmail(req.Email)
	if err != nil {
		http.Error(w, "Failed to encrypt email", http.StatusInternalServerError)
		return
	}

	// Encrypt first name
	encryptedFirstName, err := EncryptEmail(req.FirstName)
	if err != nil {
		http.Error(w, "Failed to encrypt first name", http.StatusInternalServerError)
		return
	}

	// Encrypt last name
	encryptedLastName, err := EncryptEmail(req.LastName)
	if err != nil {
		http.Error(w, "Failed to encrypt last name", http.StatusInternalServerError)
		return
	}

	user := &models.User{
		Email:     encryptedEmail,
		Password:  string(hashedPassword),
		FirstName: encryptedFirstName,
		LastName:  encryptedLastName,
		Role:      req.Role,
		StudentID: req.StudentID,
		IsActive:  true,
	}

	if err := c.userRepo.Create(user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Don't return password
	user.Password = ""
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Encrypt email for database lookup
	encryptedEmail, err := EncryptEmail(req.Email)
	if err != nil {
		http.Error(w, "Failed to process request", http.StatusInternalServerError)
		return
	}

	user, err := c.userRepo.GetByEmail(encryptedEmail)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !user.IsActive {
		http.Error(w, "Account is inactive", http.StatusUnauthorized)
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Decrypt email for response
	decryptedEmail, err := DecryptEmail(user.Email)
	if err != nil {
		http.Error(w, "Failed to process request", http.StatusInternalServerError)
		return
	}
	user.Email = decryptedEmail

	// Decrypt first name
	decryptedFirstName, err := DecryptEmail(user.FirstName)
	if err == nil {
		user.FirstName = decryptedFirstName
	}

	// Decrypt last name
	decryptedLastName, err := DecryptEmail(user.LastName)
	if err == nil {
		user.LastName = decryptedLastName
	}

	// TODO: Generate JWT token here
	token := "jwt_token_placeholder"

	user.Password = ""
	response := models.LoginResponse{
		Token: token,
		User:  *user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := c.userRepo.GetByID(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Decrypt email
	decryptedEmail, err := DecryptEmail(user.Email)
	if err != nil {
		http.Error(w, "Failed to process request", http.StatusInternalServerError)
		return
	}
	user.Email = decryptedEmail

	// Decrypt first name
	decryptedFirstName, err := DecryptEmail(user.FirstName)
	if err == nil {
		user.FirstName = decryptedFirstName
	}

	// Decrypt last name
	decryptedLastName, err := DecryptEmail(user.LastName)
	if err == nil {
		user.LastName = decryptedLastName
	}

	user.Password = ""
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (c *UserController) ListUsers(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit == 0 {
		limit = 50
	}

	users, err := c.userRepo.List(limit, offset)
	if err != nil {
		http.Error(w, "Failed to list users", http.StatusInternalServerError)
		return
	}

	// Remove passwords and decrypt emails and names from response
	for i := range users {
		users[i].Password = ""
		decryptedEmail, err := DecryptEmail(users[i].Email)
		if err == nil {
			users[i].Email = decryptedEmail
		}
		decryptedFirstName, err := DecryptEmail(users[i].FirstName)
		if err == nil {
			users[i].FirstName = decryptedFirstName
		}
		decryptedLastName, err := DecryptEmail(users[i].LastName)
		if err == nil {
			users[i].LastName = decryptedLastName
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := c.userRepo.GetByID(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Update fields if provided
	if req.FirstName != nil {
		encryptedFirstName, err := EncryptEmail(*req.FirstName)
		if err != nil {
			http.Error(w, "Failed to encrypt first name", http.StatusInternalServerError)
			return
		}
		user.FirstName = encryptedFirstName
	}
	if req.LastName != nil {
		encryptedLastName, err := EncryptEmail(*req.LastName)
		if err != nil {
			http.Error(w, "Failed to encrypt last name", http.StatusInternalServerError)
			return
		}
		user.LastName = encryptedLastName
	}
	if req.Email != nil {
		encryptedEmail, err := EncryptEmail(*req.Email)
		if err != nil {
			http.Error(w, "Failed to encrypt email", http.StatusInternalServerError)
			return
		}
		user.Email = encryptedEmail
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.StudentID != nil {
		user.StudentID = req.StudentID
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := c.userRepo.Update(user); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	// Decrypt email and names for response
	decryptedEmail, _ := DecryptEmail(user.Email)
	user.Email = decryptedEmail

	decryptedFirstName, _ := DecryptEmail(user.FirstName)
	user.FirstName = decryptedFirstName

	decryptedLastName, _ := DecryptEmail(user.LastName)
	user.LastName = decryptedLastName

	user.Password = ""
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := c.userRepo.Delete(id); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
