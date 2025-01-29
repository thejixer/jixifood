package repository

import (
	"context"
	"database/sql"
	"time"

	apperrors "github.com/thejixer/jixifood/shared/errors"
	"github.com/thejixer/jixifood/shared/models"
)

func (s *PostgresStore) CreateTypes() {
	s.db.Query(`CREATE TYPE valid_user_status AS ENUM ('incomplete', 'complete');`)
}

func (s *PostgresStore) createAuthTables() error {

	query := `
		CREATE TABLE if not exists roles (
			id SERIAL PRIMARY KEY,
			name VARCHAR(30) UNIQUE NOT NULL,
			description TEXT
		);
	
		CREATE TABLE if not exists permissions (
			id SERIAL PRIMARY KEY,
			name VARCHAR(30) UNIQUE NOT NULL,
			description TEXT
		);
	
		CREATE TABLE if not exists role_permissions (
			id SERIAL PRIMARY KEY,
			role_id INT NOT NULL,
			permission_id INT NOT NULL,
			FOREIGN KEY (role_id) REFERENCES roles(id),
			FOREIGN KEY (permission_id) REFERENCES permissions(id),
			UNIQUE(role_id, permission_id)
		);

		CREATE TABLE IF NOT EXISTS seeded (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);

		create table if not exists users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100),
			phone_number VARCHAR(20) UNIQUE,
			status valid_user_status,
			role_id int,
			createdAt TIMESTAMPTZ
		);

	`

	_, err := s.db.Exec(query)

	return err
}

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (r *AuthRepo) CreateUser(ctx context.Context, phoneNumber string, roleID uint64) (*models.UserEntity, error) {

	NewUser := &models.UserEntity{
		Name:        "",
		PhoneNumber: phoneNumber,
		Status:      "incomplete",
		RoleID:      roleID,
		CreatedAt:   time.Now().UTC(),
	}

	// the roleId 0 is intended to be a customer
	if NewUser.RoleID == 0 {
		err := r.db.QueryRowContext(ctx, `SELECT id FROM ROLES WHERE NAME = 'customer'`).Scan(&NewUser.RoleID)
		if err != nil {
			return nil, err
		}
	}

	query := `
		INSERT INTO users (name, phone_number, status, role_id, createdAt)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`

	lastInsertId := 0

	insertErr := r.db.QueryRowContext(
		ctx,
		query,
		NewUser.Name,
		NewUser.PhoneNumber,
		NewUser.Status,
		NewUser.RoleID,
		NewUser.CreatedAt,
	).Scan(&lastInsertId)

	if insertErr != nil {
		return nil, insertErr
	}

	NewUser.ID = uint64(lastInsertId)

	return NewUser, nil
}

func (r *AuthRepo) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*models.UserEntity, error) {

	rows, err := r.db.Query("SELECT * FROM USERS WHERE phone_number = $1", phoneNumber)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return ScanIntoUserEntity(rows)
	}

	return nil, apperrors.ErrNotFound
}

func (r *AuthRepo) GetUserByID(ctx context.Context, id uint64) (*models.UserEntity, error) {

	rows, err := r.db.Query("SELECT * FROM USERS WHERE id = $1", id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return ScanIntoUserEntity(rows)
	}

	return nil, apperrors.ErrNotFound
}

func (r *AuthRepo) GetRoleByID(ctx context.Context, id uint64) (*models.Role, error) {
	rows, err := r.db.Query("SELECT * FROM roles WHERE id = $1", id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return ScanIntoRole(rows)
	}

	return nil, apperrors.ErrNotFound
}

func ScanIntoUserEntity(rows *sql.Rows) (*models.UserEntity, error) {
	u := new(models.UserEntity)
	if err := rows.Scan(
		&u.ID,
		&u.Name,
		&u.PhoneNumber,
		&u.Status,
		&u.RoleID,
		&u.CreatedAt,
	); err != nil {
		return nil, err
	}
	return u, nil
}

func ScanIntoRole(rows *sql.Rows) (*models.Role, error) {
	role := new(models.Role)
	if err := rows.Scan(
		&role.ID,
		&role.Name,
		&role.Description,
	); err != nil {
		return nil, err
	}
	return role, nil
}
