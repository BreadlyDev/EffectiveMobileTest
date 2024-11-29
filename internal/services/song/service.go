package song

import (
	"effectivemobiletesttask/internal/config"
	"effectivemobiletesttask/internal/domain/models"
	"effectivemobiletesttask/internal/services"
	"log/slog"
)

type Provider interface {
	// Song
	CreateSong(song models.SongStorage) (int64, error)
	GetSongByID(id int64) (models.SongStorage, error)
	GetSongByName(songName string) (models.SongStorage, error)
	UpdateSong(id int64, song models.SongStorage) (models.SongStorage, error)
	DeleteSong(id int64) error
	GetAllSongs(filter models.SongFilter, groupID int64, offset int, limit int) ([]models.SongStorage, error)

	// Group
	CreateGroup(groupName string) (int64, error)
	GetGroupByID(id int64) (models.Group, error)
	GetGroupByName(groupName string) (models.Group, error)
}

type Service struct {
	log      *slog.Logger
	provider Provider
	api      config.APIClient
}

func New(log *slog.Logger, provider Provider, api config.APIClient) *Service {
	return &Service{
		log:      log,
		provider: provider,
		api:      api,
	}
}

func SongToSongResp(song models.SongStorage, groupName string) models.SongResponse {
	var songResp models.SongResponse

	songResp.ID = song.ID
	songResp.Group = groupName
	songResp.Name = song.Name
	songResp.ReleaseDate = song.ReleaseDate
	songResp.Text = song.Text
	songResp.Link = song.Link

	return songResp
}

func SongRespToSongStorage(songResp models.SongResponse, groupID int64) models.SongStorage {
	var song models.SongStorage

	song.ID = songResp.ID
	song.GroupID = groupID
	song.Name = songResp.Name
	song.ReleaseDate = songResp.ReleaseDate
	song.Text = songResp.Text
	song.Link = songResp.Link

	return song
}

func SongReqAndDetsToSong(
	songReq models.SongRequest,
	songDetail models.SongDetail,
	groupID int64,
) models.SongStorage {
	var song models.SongStorage

	song.GroupID = groupID
	song.Name = songReq.Name
	song.ReleaseDate = songDetail.ReleaseDate
	song.Text = songDetail.Text
	song.Link = songDetail.Link

	return song
}

func ValidateSongDetails(songDetail models.SongDetail) (string, error) {
	if songDetail.ReleaseDate.IsZero() {
		return "releaseDate", services.ErrFieldIsRequired
	}
	if songDetail.Text == "" {
		return "text", services.ErrFieldIsRequired
	}
	if songDetail.Link == "" {
		return "link", services.ErrFieldIsRequired
	}

	return "", nil
}
