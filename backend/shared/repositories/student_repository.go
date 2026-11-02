package repositories

import (
	"database/sql"

	"github.com/kwizera/university-invoices/services/student/models"
)

type StudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) Create(student *models.Student) error {
	query := `INSERT INTO students (student_id, first_name, last_name, email, phone, program, 
			  semester, year, date_of_birth, address, is_active)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, student.StudentID, student.FirstName, student.LastName,
		student.Email, student.Phone, student.Program, student.Semester, student.Year,
		student.DateOfBirth, student.Address, student.IsActive)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	student.ID = int(id)
	return nil
}

func (r *StudentRepository) GetByID(id int) (*models.Student, error) {
	query := `SELECT id, student_id, first_name, last_name, email, phone, program, semester, year,
			  date_of_birth, address, is_active, created_at, updated_at
			  FROM students WHERE id = ?`

	student := &models.Student{}
	err := r.db.QueryRow(query, id).Scan(
		&student.ID, &student.StudentID, &student.FirstName, &student.LastName,
		&student.Email, &student.Phone, &student.Program, &student.Semester, &student.Year,
		&student.DateOfBirth, &student.Address, &student.IsActive, &student.CreatedAt, &student.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return student, nil
}

func (r *StudentRepository) GetByStudentID(studentID string) (*models.Student, error) {
	query := `SELECT id, student_id, first_name, last_name, email, phone, program, semester, year,
			  date_of_birth, address, is_active, created_at, updated_at
			  FROM students WHERE student_id = ?`

	student := &models.Student{}
	err := r.db.QueryRow(query, studentID).Scan(
		&student.ID, &student.StudentID, &student.FirstName, &student.LastName,
		&student.Email, &student.Phone, &student.Program, &student.Semester, &student.Year,
		&student.DateOfBirth, &student.Address, &student.IsActive, &student.CreatedAt, &student.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return student, nil
}

func (r *StudentRepository) List(limit, offset int) ([]models.Student, error) {
	query := `SELECT id, student_id, first_name, last_name, email, phone, program, semester, year,
			  date_of_birth, address, is_active, created_at, updated_at
			  FROM students LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var student models.Student
		err := rows.Scan(
			&student.ID, &student.StudentID, &student.FirstName, &student.LastName,
			&student.Email, &student.Phone, &student.Program, &student.Semester, &student.Year,
			&student.DateOfBirth, &student.Address, &student.IsActive, &student.CreatedAt, &student.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	return students, nil
}

func (r *StudentRepository) Update(student *models.Student) error {
	query := `UPDATE students SET first_name = ?, last_name = ?, email = ?, phone = ?, 
			  program = ?, semester = ?, year = ?, date_of_birth = ?, address = ?, 
			  is_active = ?, updated_at = NOW() WHERE id = ?`

	_, err := r.db.Exec(query, student.FirstName, student.LastName, student.Email, student.Phone,
		student.Program, student.Semester, student.Year, student.DateOfBirth, student.Address,
		student.IsActive, student.ID)
	return err
}

func (r *StudentRepository) Delete(id int) error {
	query := `DELETE FROM students WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
