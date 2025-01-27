package repository

import "database/sql"

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
