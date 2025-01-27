package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/thejixer/jixifood/services/auth/internal/config"
	"github.com/thejixer/jixifood/shared/models"
)

type PostgresStore struct {
	db       *sql.DB
	AuthRepo models.AuthRepository
}

func NewPostgresStore(cfg *config.AuthServiceConfig) (*PostgresStore, error) {

	dbName := cfg.DB_NAME
	dbUser := cfg.DB_USER
	dbPassword := cfg.DB_PASSWORD
	dbHost := cfg.DB_HOST
	conString := fmt.Sprintf("user=%v dbname=%v password=%v sslmode=disable host=%v", dbUser, dbName, dbPassword, dbHost)
	db, err := sql.Open("postgres", conString)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	AuthRepo := NewAuthRepo(db)

	return &PostgresStore{
		db:       db,
		AuthRepo: AuthRepo,
	}, nil
}

func (s *PostgresStore) Init() error {

	return nil
}
