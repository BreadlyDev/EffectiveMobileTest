package postgres

import (
	"database/sql"
	"effectivemobiletesttask/internal/domain/models"
	"effectivemobiletesttask/internal/storage"
	"errors"
	"fmt"
)

func (s *Storage) CreateGroup(groupName string) (int64, error) {
	const op = "storage.postgres.CreateGroup"

	stmt, err := s.db.Prepare("INSERT INTO groups(name) VALUES ($1) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var id int64

	err = stmt.QueryRow(groupName).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetGroupByID(id int64) (models.Group, error) {
	const op = "storage.postgres.GetGroupByID"

	stmt, err := s.db.Prepare("SELECT * FROM groups WHERE id = $1")
	if err != nil {
		return models.Group{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)

	var group models.Group
	err = row.Scan(&group.ID, &group.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Group{}, fmt.Errorf("%s: %w", op, storage.ErrGroupNotFound)
		}
		return models.Group{}, fmt.Errorf("%s: %w", op, err)
	}

	return group, nil
}

func (s *Storage) GetGroupByName(groupName string) (models.Group, error) {
	const op = "storage.postgres.GetGroupByName"

	stmt, err := s.db.Prepare("SELECT * FROM groups WHERE name = $1")
	if err != nil {
		return models.Group{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(groupName)

	var group models.Group

	err = row.Scan(&group.ID, &group.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Group{}, fmt.Errorf("%s: %w", op, storage.ErrGroupNotFound)
		}
		return models.Group{}, fmt.Errorf("%s: %w", op, err)
	}

	return group, nil
}
