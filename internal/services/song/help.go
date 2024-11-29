package song

import (
	"effectivemobiletesttask/internal/domain/models"
	"effectivemobiletesttask/internal/storage"
	lg "effectivemobiletesttask/internal/utils/logger"
	"errors"
	"fmt"
	"log/slog"
)

func (s *Service) createOrGetGroup(groupName string) (int64, error) {
	group, err := s.GetGroupByName(groupName)
	if err != nil && !errors.Is(err, storage.ErrGroupNotFound) {
		return 0, fmt.Errorf("error fetching group: %w", err)
	}

	if group.ID == 0 {
		groupID, err := s.CreateGroup(groupName)

		if err != nil {
			return 0, fmt.Errorf("error creating group: %w", err)
		}
		return groupID, nil
	}
	return group.ID, nil
}

func (s *Service) getSongAndGroup(id int64) (models.SongResponse, models.Group, error) {
	song, err := s.provider.GetSongByID(id)
	if err != nil {
		return models.SongResponse{}, models.Group{}, fmt.Errorf("error fetching song: %w", err)
	}

	group, err := s.provider.GetGroupByID(song.GroupID)
	if err != nil {
		return models.SongResponse{}, models.Group{}, fmt.Errorf("error fetching group: %w", err)
	}

	return SongToSongResp(song, group.Name), group, nil
}

func (s *Service) logSongsWithoutText(songs []models.SongStorage) []models.SongStorage {
	var songsWithoutText []models.SongStorage
	for _, song := range songs {
		songCopy := song
		songCopy.Text = ""
		songCopy.Link = ""
		songsWithoutText = append(songsWithoutText, songCopy)
	}
	s.log.Debug("fetched songs: ", slog.Any("songs", songsWithoutText))
	return songsWithoutText
}

func (s *Service) fetchGroups(songs []models.SongStorage) ([]models.Group, error) {
	var groups []models.Group
	for _, song := range songs {
		group, err := s.provider.GetGroupByID(song.GroupID)
		if err != nil {
			s.log.Error("error during the fetching group", lg.Err(err))
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}
