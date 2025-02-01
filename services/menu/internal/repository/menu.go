package repository

import "database/sql"

func (s *PostgresStore) createMenuTables() error {
	return nil
}

type MenuRepo struct {
	db *sql.DB
}

func NewMenuRepo(db *sql.DB) *MenuRepo {
	return &MenuRepo{
		db: db,
	}
}
