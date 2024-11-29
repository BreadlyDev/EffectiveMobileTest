package postgres

import (
	"database/sql"
	"effectivemobiletesttask/internal/config"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func New(cfg config.DBStorage) (*Storage, error) {
	const op = "storage.postgres.New"

	DBUrl := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
		cfg.User, cfg.Pass, cfg.DBName, cfg.Host, cfg.Port, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", DBUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}
