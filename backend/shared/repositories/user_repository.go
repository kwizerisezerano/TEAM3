package repositories

import (
	"database/sql"

	"github.com/kwizera/university-invoices/services/user/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `INSERT INTO users (email, password, first_name, last_name, role, student_id, is_active)
			  VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, user.Email, user.Password, user.FirstName,
		user.LastName, user.Role, user.StudentID, user.IsActive)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	return nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	query := `SELECT id, email, password, first_name, last_name, role, student_id, is_active, created_at, updated_at
			  FROM users WHERE id = ?`

	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.Role, &user.StudentID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, password, first_name, last_name, role, student_id, is_active, created_at, updated_at
			  FROM users WHERE email = ?`

	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.Role, &user.StudentID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) List(limit, offset int) ([]models.User, error) {
	query := `SELECT id, email, first_name, last_name, role, student_id, is_active, created_at, updated_at
			  FROM users LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Email, &user.FirstName, &user.LastName,
			&user.Role, &user.StudentID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) Update(user *models.User) error {
	query := `UPDATE users SET email = ?, first_name = ?, last_name = ?, role = ?, 
			  student_id = ?, is_active = ?, updated_at = NOW() WHERE id = ?`

	_, err := r.db.Exec(query, user.Email, user.FirstName, user.LastName,
		user.Role, user.StudentID, user.IsActive, user.ID)
	return err
}

func (r *UserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
