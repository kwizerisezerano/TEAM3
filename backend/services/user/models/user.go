package models

import "time"

type Role string

const (
	RoleStudent   Role = "student"
	RoleStaff     Role = "staff"
	RoleAdmin     Role = "admin"
	RoleFinance   Role = "finance"
	RoleRegistrar Role = "registrar"
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Never expose in JSON
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      Role      `json:"role"`
	StudentID *string   `json:"student_id,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Email     string  `json:"email" binding:"required,email"`
	Password  string  `json:"password" binding:"required,min=8"`
	FirstName string  `json:"first_name" binding:"required"`
	LastName  string  `json:"last_name" binding:"required"`
	Role      Role    `json:"role" binding:"required"`
	StudentID *string `json:"student_id,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type UpdateUserRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email"`
	Role      *Role   `json:"role"`
	StudentID *string `json:"student_id"`
	IsActive  *bool   `json:"is_active"`
}
