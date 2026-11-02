package models

import "time"

type Student struct {
	ID          int        `json:"id"`
	StudentID   string     `json:"student_id"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Email       string     `json:"email"`
	Phone       string     `json:"phone"`
	Program     string     `json:"program"`
	Semester    string     `json:"semester"`
	Year        int        `json:"year"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	Address     string     `json:"address"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateStudentRequest struct {
	StudentID   string     `json:"student_id" binding:"required"`
	FirstName   string     `json:"first_name" binding:"required"`
	LastName    string     `json:"last_name" binding:"required"`
	Email       string     `json:"email" binding:"required,email"`
	Phone       string     `json:"phone" binding:"required"`
	Program     string     `json:"program" binding:"required"`
	Semester    string     `json:"semester" binding:"required"`
	Year        int        `json:"year" binding:"required"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Address     string     `json:"address"`
}

type UpdateStudentRequest struct {
	FirstName   *string    `json:"first_name"`
	LastName    *string    `json:"last_name"`
	Email       *string    `json:"email"`
	Phone       *string    `json:"phone"`
	Program     *string    `json:"program"`
	Semester    *string    `json:"semester"`
	Year        *int       `json:"year"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Address     *string    `json:"address"`
	IsActive    *bool      `json:"is_active"`
}
