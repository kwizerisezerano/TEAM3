package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kwizera/university-invoices/internal/encryption"
	"github.com/kwizera/university-invoices/services/student/models"
	"github.com/kwizera/university-invoices/shared/repositories"
)

// DecryptEmail decrypts an encrypted email address
func DecryptEmail(encryptedEmail string) (string, error) {
	key := encryption.GetEncryptionKey()
	return encryption.Decrypt(encryptedEmail, key)
}

// EncryptEmail encrypts an email address
func EncryptEmail(email string) (string, error) {
	key := encryption.GetEncryptionKey()
	return encryption.Encrypt(email, key)
}

type StudentController struct {
	studentRepo *repositories.StudentRepository
}

func NewStudentController(studentRepo *repositories.StudentRepository) *StudentController {
	return &StudentController{studentRepo: studentRepo}
}

func (c *StudentController) CreateStudent(w http.ResponseWriter, r *http.Request) {
	var req models.CreateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
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

	student := &models.Student{
		StudentID:   req.StudentID,
		FirstName:   encryptedFirstName,
		LastName:    encryptedLastName,
		Email:       encryptedEmail,
		Phone:       req.Phone,
		Program:     req.Program,
		Semester:    req.Semester,
		Year:        req.Year,
		DateOfBirth: req.DateOfBirth,
		Address:     req.Address,
		IsActive:    true,
	}

	if err := c.studentRepo.Create(student); err != nil {
		http.Error(w, "Failed to create student", http.StatusInternalServerError)
		return
	}

	// Decrypt email for response
	decryptedEmail, _ := DecryptEmail(student.Email)
	student.Email = decryptedEmail

	// Decrypt first name
	decryptedFirstName, _ := DecryptEmail(student.FirstName)
	student.FirstName = decryptedFirstName

	// Decrypt last name
	decryptedLastName, _ := DecryptEmail(student.LastName)
	student.LastName = decryptedLastName

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(student)
}

func (c *StudentController) GetStudent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	student, err := c.studentRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	// Decrypt email
	decryptedEmail, err := DecryptEmail(student.Email)
	if err == nil {
		student.Email = decryptedEmail
	}

	// Decrypt first name
	decryptedFirstName, err := DecryptEmail(student.FirstName)
	if err == nil {
		student.FirstName = decryptedFirstName
	}

	// Decrypt last name
	decryptedLastName, err := DecryptEmail(student.LastName)
	if err == nil {
		student.LastName = decryptedLastName
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

func (c *StudentController) GetStudentByStudentID(w http.ResponseWriter, r *http.Request) {
	studentID := r.URL.Query().Get("student_id")
	if studentID == "" {
		http.Error(w, "Student ID is required", http.StatusBadRequest)
		return
	}

	student, err := c.studentRepo.GetByStudentID(studentID)
	if err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	// Decrypt email
	decryptedEmail, err := DecryptEmail(student.Email)
	if err == nil {
		student.Email = decryptedEmail
	}

	// Decrypt first name
	decryptedFirstName, err := DecryptEmail(student.FirstName)
	if err == nil {
		student.FirstName = decryptedFirstName
	}

	// Decrypt last name
	decryptedLastName, err := DecryptEmail(student.LastName)
	if err == nil {
		student.LastName = decryptedLastName
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

func (c *StudentController) ListStudents(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit == 0 {
		limit = 50
	}

	students, err := c.studentRepo.List(limit, offset)
	if err != nil {
		http.Error(w, "Failed to list students", http.StatusInternalServerError)
		return
	}

	// Decrypt emails and names
	for i := range students {
		decryptedEmail, err := DecryptEmail(students[i].Email)
		if err == nil {
			students[i].Email = decryptedEmail
		}
		decryptedFirstName, err := DecryptEmail(students[i].FirstName)
		if err == nil {
			students[i].FirstName = decryptedFirstName
		}
		decryptedLastName, err := DecryptEmail(students[i].LastName)
		if err == nil {
			students[i].LastName = decryptedLastName
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func (c *StudentController) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	student, err := c.studentRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	// Update fields if provided
	if req.FirstName != nil {
		encryptedFirstName, err := EncryptEmail(*req.FirstName)
		if err != nil {
			http.Error(w, "Failed to encrypt first name", http.StatusInternalServerError)
			return
		}
		student.FirstName = encryptedFirstName
	}
	if req.LastName != nil {
		encryptedLastName, err := EncryptEmail(*req.LastName)
		if err != nil {
			http.Error(w, "Failed to encrypt last name", http.StatusInternalServerError)
			return
		}
		student.LastName = encryptedLastName
	}
	if req.Email != nil {
		encryptedEmail, err := EncryptEmail(*req.Email)
		if err != nil {
			http.Error(w, "Failed to encrypt email", http.StatusInternalServerError)
			return
		}
		student.Email = encryptedEmail
	}
	if req.Phone != nil {
		student.Phone = *req.Phone
	}
	if req.Program != nil {
		student.Program = *req.Program
	}
	if req.Semester != nil {
		student.Semester = *req.Semester
	}
	if req.Year != nil {
		student.Year = *req.Year
	}
	if req.DateOfBirth != nil {
		student.DateOfBirth = req.DateOfBirth
	}
	if req.Address != nil {
		student.Address = *req.Address
	}
	if req.IsActive != nil {
		student.IsActive = *req.IsActive
	}

	if err := c.studentRepo.Update(student); err != nil {
		http.Error(w, "Failed to update student", http.StatusInternalServerError)
		return
	}

	// Decrypt email and names for response
	decryptedEmail, _ := DecryptEmail(student.Email)
	student.Email = decryptedEmail

	decryptedFirstName, _ := DecryptEmail(student.FirstName)
	student.FirstName = decryptedFirstName

	decryptedLastName, _ := DecryptEmail(student.LastName)
	student.LastName = decryptedLastName

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

func (c *StudentController) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	if err := c.studentRepo.Delete(id); err != nil {
		http.Error(w, "Failed to delete student", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
