package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/thejixer/jixifood/shared/constants"
	apperrors "github.com/thejixer/jixifood/shared/errors"
	"github.com/thejixer/jixifood/shared/models"
)

func (s *PostgresStore) CreateTypes() {
	query := fmt.Sprintf(`CREATE TYPE valid_user_status AS ENUM ('%s', '%s');`,
		constants.UserStatusIncomplete,
		constants.UserStatusComplete,
	)

	s.db.Query(query)

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
			role_id INT NOT NULL,
			FOREIGN KEY (role_id) REFERENCES roles(id),
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

func (r *AuthRepo) CreateUser(ctx context.Context, phoneNumber, name string, roleID uint64) (*models.UserEntity, error) {

	NewUser := &models.UserEntity{
		Name:        name,
		PhoneNumber: phoneNumber,
		Status: func(name string) string {
			if name == "" {
				return constants.UserStatusIncomplete
			} else {
				return constants.UserStatusComplete
			}
		}(name),
		RoleID:    roleID,
		CreatedAt: time.Now().UTC(),
	}

	// the roleId 0 is intended to be a customer
	if NewUser.RoleID == 0 {
		err := r.db.QueryRowContext(ctx, `SELECT id FROM ROLES WHERE NAME = $1`, constants.RoleCustomer).Scan(&NewUser.RoleID)
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
		var pqErr *pq.Error
		if errors.As(insertErr, &pqErr) {
			switch pqErr.Code {
			case constants.PGDuplicateKeyErrorCode:
				return nil, fmt.Errorf("error in authRepo.createUser: %w: %v", apperrors.ErrDuplicatePhone, insertErr)
			case constants.PGForeignKeyViolationCode:
				return nil, fmt.Errorf("error in authRepo.createUser: bad roleID: %w: %v", apperrors.ErrInputRequirements, insertErr)
			}

		}

		return nil, fmt.Errorf("error in authRepo.createUser: %w", insertErr)

	}

	NewUser.ID = uint64(lastInsertId)

	return NewUser, nil
}

func (r *AuthRepo) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*models.UserEntity, error) {

	rows, err := r.db.QueryContext(ctx, "SELECT * FROM USERS WHERE phone_number = $1", phoneNumber)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return ScanIntoUserEntity(rows)
	}

	return nil, apperrors.ErrNotFound
}

func (r *AuthRepo) GetUserByID(ctx context.Context, id uint64) (*models.UserEntity, error) {

	rows, err := r.db.QueryContext(ctx, "SELECT * FROM USERS WHERE id = $1", id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return ScanIntoUserEntity(rows)
	}

	return nil, apperrors.ErrNotFound
}

func (r *AuthRepo) GetRoleByID(ctx context.Context, id uint64) (*models.Role, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM roles WHERE id = $1", id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return ScanIntoRole(rows)
	}

	return nil, apperrors.ErrNotFound
}

func (r *AuthRepo) CheckPermission(ctx context.Context, roleID uint64, permissionName string) (bool, error) {
	query := `
	SELECT EXISTS (
		SELECT 1
		FROM role_permissions rp
		JOIN permissions p
		ON rp.permission_id = p.id
		WHERE rp.role_id = $1 AND p.name = $2
	);
`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, roleID, permissionName).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *AuthRepo) ChangeUserRole(ctx context.Context, userID, roleID uint64) (*models.UserEntity, error) {
	query := `
		UPDATE users
		SET role_id = $1
		WHERE id = $2
		RETURNING *;
	`
	rows, err := r.db.QueryContext(ctx, query, roleID, userID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case constants.PGForeignKeyViolationCode:
				return nil, fmt.Errorf("error in authRepo.changeUserRole: bad roleID: %w: %v", apperrors.ErrInputRequirements, err)
			}

		}
		return nil, fmt.Errorf("error in authRepo.changeUserRole: %w", err)

	}
	for rows.Next() {
		return ScanIntoUserEntity(rows)
	}
	return nil, apperrors.ErrNotFound

}

func (r *AuthRepo) EditProfile(ctx context.Context, userID uint64, name string) (*models.UserEntity, error) {
	query := `
		UPDATE users
		SET name = $2, status = $3
		WHERE id = $1
		RETURNING *;
	`
	rows, err := r.db.QueryContext(ctx, query, userID, name, constants.UserStatusComplete)
	if err != nil {
		return nil, fmt.Errorf("error in authRepo.changeUserRole: %w", err)
	}
	for rows.Next() {
		return ScanIntoUserEntity(rows)
	}
	return nil, fmt.Errorf("error in authRepo.editProfile: %w: %v", apperrors.ErrNotFound, err)

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
