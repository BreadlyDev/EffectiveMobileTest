package song

import (
	"database/sql"
	"effectivemobiletesttask/internal/domain/models"
	"effectivemobiletesttask/internal/storage"
	lg "effectivemobiletesttask/internal/utils/logger"
	"errors"
	"fmt"
	"log/slog"
)

func (s *Service) CreateGroup(groupName string) (int64, error) {
	const op = "services.song.CreateGroup"

	s.log.Debug("start creating group")
	id, err := s.provider.CreateGroup(groupName)
	if err != nil {
		s.log.Error("error during creating group", lg.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	s.log.Debug("created group", slog.Any("id", id))

	return id, nil
}

func (s *Service) GetGroupByID(id int64) (models.Group, error) {
	const op = "services.song.GetGroupByID"

	s.log.Debug("start fetching group")
	group, err := s.provider.GetGroupByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Error("group was not found")

			return models.Group{}, storage.ErrGroupNotFound
		}

		s.log.Error("error during fetching group", lg.Err(err))
		return models.Group{}, fmt.Errorf("%s: %w", op, err)
	}
	s.log.Debug("fetched group", slog.Any("group", group))

	return group, nil
}

func (s *Service) GetGroupByName(groupName string) (models.Group, error) {
	const op = "services.song.GetGroupByName"

	s.log.Debug("start fetching group")
	group, err := s.provider.GetGroupByName(groupName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Error("group was not found")

			return models.Group{}, storage.ErrGroupNotFound
		}

		s.log.Error("error during fetching group", lg.Err(err))
		return models.Group{}, fmt.Errorf("%s: %w", op, err)
	}
	s.log.Debug("fetched group", slog.Any("group", group))

	return group, nil
}
