package repositories

import (
	"database/sql"

	"github.com/kwizera/university-invoices/services/user/models"
)

type AuditRepository struct {
	db *sql.DB
}

func NewAuditRepository(db *sql.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Create(audit *models.AuditLog) error {
	query := `INSERT INTO audit_logs (user_id, action, entity, entity_id, old_value, new_value, ip_address, user_agent)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, audit.UserID, audit.Action, audit.Entity, audit.EntityID,
		audit.OldValue, audit.NewValue, audit.IPAddress, audit.UserAgent)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	audit.ID = int(id)
	return nil
}

func (r *AuditRepository) GetByEntity(entity string, entityID int, limit, offset int) ([]models.AuditLog, error) {
	query := `SELECT id, user_id, action, entity, entity_id, old_value, new_value, ip_address, user_agent, created_at
			  FROM audit_logs WHERE entity = ? AND entity_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, entity, entityID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var audits []models.AuditLog
	for rows.Next() {
		var audit models.AuditLog
		err := rows.Scan(
			&audit.ID, &audit.UserID, &audit.Action, &audit.Entity, &audit.EntityID,
			&audit.OldValue, &audit.NewValue, &audit.IPAddress, &audit.UserAgent, &audit.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		audits = append(audits, audit)
	}

	return audits, nil
}

func (r *AuditRepository) GetByUserID(userID int, limit, offset int) ([]models.AuditLog, error) {
	query := `SELECT id, user_id, action, entity, entity_id, old_value, new_value, ip_address, user_agent, created_at
			  FROM audit_logs WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var audits []models.AuditLog
	for rows.Next() {
		var audit models.AuditLog
		err := rows.Scan(
			&audit.ID, &audit.UserID, &audit.Action, &audit.Entity, &audit.EntityID,
			&audit.OldValue, &audit.NewValue, &audit.IPAddress, &audit.UserAgent, &audit.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		audits = append(audits, audit)
	}

	return audits, nil
}

func (r *AuditRepository) List(limit, offset int) ([]models.AuditLog, error) {
	query := `SELECT id, user_id, action, entity, entity_id, old_value, new_value, ip_address, user_agent, created_at
			  FROM audit_logs ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var audits []models.AuditLog
	for rows.Next() {
		var audit models.AuditLog
		err := rows.Scan(
			&audit.ID, &audit.UserID, &audit.Action, &audit.Entity, &audit.EntityID,
			&audit.OldValue, &audit.NewValue, &audit.IPAddress, &audit.UserAgent, &audit.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		audits = append(audits, audit)
	}

	return audits, nil
}
