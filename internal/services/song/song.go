package song

import (
	"effectivemobiletesttask/internal/client/song"
	"effectivemobiletesttask/internal/domain/models"
	"effectivemobiletesttask/internal/storage"
	lg "effectivemobiletesttask/internal/utils/logger"
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

func (s *Service) CreateSong(songReq models.SongRequest) (int64, error) {
	const op = "services.song.CreateSong"
	s.log.With(slog.String("operation", op))
	s.log.Debug("start song creation", slog.String("songName", songReq.Name), slog.String("group", songReq.Group))

	groupID, err := s.createOrGetGroup(songReq.Group)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	s.log.Debug("group retrieved or created", slog.Int64("groupID", groupID))

	songDetail, err := song.NewClient(s.log, s.api).GetSongDetail(songReq)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	s.log.Debug("fetched song details", slog.String("songReleaseDate", songDetail.ReleaseDate.Format("2006-01-02")))

	field, err := ValidateSongDetails(songDetail)
	if err != nil {
		return 0, fmt.Errorf("%s: '%s' %w", op, field, err)
	}

	song := SongReqAndDetsToSong(songReq, songDetail, groupID)
	id, err := s.provider.CreateSong(song)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	s.log.Info("song created successfully", slog.Int64("songID", id), slog.String("songReleaseDate", songDetail.ReleaseDate.Format("2006-01-02")))
	return id, nil
}

func (s *Service) GetSongByID(id int64) (models.SongResponse, error) {
	const op = "services.song.GetSongByID"
	s.log.With(slog.String("operation", op))
	s.log.Debug("start fetching song by ID", slog.Int64("songID", id))

	songResp, _, err := s.getSongAndGroup(id)
	if err != nil {
		s.log.Error("error fetching song", lg.Err(err))
		return models.SongResponse{}, err
	}

	s.log.Debug("fetched song successfully", slog.Int64("songID", id), slog.String("songName", songResp.Name))
	return songResp, nil
}

func (s *Service) GetSongByName(songName string) (models.SongResponse, error) {
	const op = "services.song.GetSongByName"
	s.log.With(slog.String("operation", op))
	s.log.Debug("start fetching song by name", slog.String("songName", songName))

	song, err := s.provider.GetSongByName(songName)
	if err != nil {
		s.log.Error("error fetching song by name", lg.Err(err))
		return models.SongResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	return s.GetSongByID(song.ID)
}

func (s *Service) GetSongTextByID(id int64, verse int) (string, error) {
	const op = "services.song.GetSongTextByID"
	s.log.With(slog.String("operation", op))
	s.log.Debug("start fetching song text by ID", slog.Int64("songID", id), slog.Int("verse", verse))

	songResp, _, err := s.getSongAndGroup(id)
	if err != nil {
		return "", err
	}

	songVerses := strings.Split(songResp.Text, "\n\n")
	if verse > len(songVerses)-1 {
		verse = len(songVerses) - 1
	}
	if verse < 0 {
		verse = 0
	}

	s.log.Debug("fetched song text", slog.String("verseText", songVerses[verse]))
	return songVerses[verse], nil
}

func (s *Service) GetSongTextByName(songName string, verse int) (string, error) {
	const op = "services.song.GetSongTextByName"
	s.log.With(slog.String("operation", op))
	s.log.Debug("start fetching song text by name", slog.String("songName", songName), slog.Int("verse", verse))

	songResp, err := s.GetSongByName(songName)
	if err != nil {
		return "", err
	}

	songVerses := strings.Split(songResp.Text, "\n\n")
	if verse > len(songVerses)-1 {
		verse = len(songVerses) - 1
	}
	if verse < 0 {
		verse = 0
	}

	s.log.Debug("fetched song text", slog.String("verseText", songVerses[verse]))
	return songVerses[verse], nil
}

func (s *Service) UpdateSong(id int64, newSong models.SongResponse) (models.SongResponse, error) {
	const op = "services.song.UpdateSong"
	s.log.With(slog.String("operation", op))

	s.log.Debug("start updating song", slog.Int64("songID", id), slog.String("songName", newSong.Name))

	groupID, err := s.createOrGetGroup(newSong.Group)
	if err != nil {
		return models.SongResponse{}, fmt.Errorf("%s: %w", op, err)
	}
	s.log.Debug("group retrieved or created for update", slog.Int64("groupID", groupID))

	_, err = s.provider.UpdateSong(id, SongRespToSongStorage(newSong, groupID))
	if err != nil {
		return models.SongResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	s.log.Debug("updated song successfully", slog.Int64("songID", id), slog.String("songName", newSong.Name))
	return newSong, nil
}

func (s *Service) DeleteSong(id int64) error {
	const op = "services.song.DeleteSong"
	s.log.With(slog.String("operation", op))

	s.log.Debug("start deleting song", slog.Int64("songID", id))
	if err := s.provider.DeleteSong(id); err != nil {
		s.log.Error("error deleting song", lg.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	s.log.Debug("deleted song successfully", slog.Int64("songID", id))

	return nil
}

func (s *Service) GetAllSongs(filter models.SongFilter, offset int, limit int) ([]models.SongResponse, error) {
	const op = "services.song.GetAllSongs"

	s.log.With(slog.String("operation", op))
	s.log.Debug("start fetching songs with filters", slog.Any("filter", filter), slog.Int("offset", offset), slog.Int("limit", limit))

	var group models.Group
	if filter.Group != "" {
		group, err := s.GetGroupByName(filter.Group)
		if err != nil {
			if errors.Is(err, storage.ErrGroupNotFound) {
				return nil, fmt.Errorf("%s: %w", op, err)
			}

			s.log.Error("error during fetching group: ", lg.Err(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		s.log.Debug("fetched group for filtering", slog.Any("group", group))
	}

	songs, err := s.provider.GetAllSongs(filter, group.ID, offset, limit)
	if err != nil {
		s.log.Error("error fetching songs", lg.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	groups, err := s.fetchGroups(songs)
	if err != nil {
		s.log.Error("error fetching groups for songs", lg.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var songResps []models.SongResponse
	for i := 0; i < len(songs); i++ {
		songResps = append(songResps, SongToSongResp(songs[i], groups[i].Name))
	}

	s.log.Debug("fetched songs successfully", slog.Int("totalSongs", len(songResps)))
	s.logSongsWithoutText(songs)

	return songResps, nil
}
