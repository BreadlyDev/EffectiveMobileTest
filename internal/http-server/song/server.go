package song

import (
	"effectivemobiletesttask/internal/domain/models"
	"log/slog"
	"net/http"
)

type Service interface {
	CreateSong(songReq models.SongRequest) (int64, error)
	GetSongByID(id int64) (models.SongResponse, error)
	GetSongByName(songName string) (models.SongResponse, error)
	GetSongTextByID(id int64, verse int) (string, error)
	GetSongTextByName(songName string, verse int) (string, error)
	UpdateSong(id int64, song models.SongResponse) (models.SongResponse, error)
	DeleteSong(id int64) error
	GetAllSongs(filter models.SongFilter, offset int, limit int) ([]models.SongResponse, error)
}

type Server struct {
	log      *slog.Logger
	pageSize int
	service  Service
}

func New(log *slog.Logger, pageSize int, service Service) *Server {
	return &Server{
		log:      log,
		pageSize: pageSize,
		service:  service,
	}
}

func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /song/create", s.CreateSong)
	mux.HandleFunc("GET /song/{id}", s.GetSongByID)
	mux.HandleFunc("GET /song", s.GetSongByName)
	mux.HandleFunc("GET /song/{id}/text", s.GetSongTextByID)
	mux.HandleFunc("GET /song/name/text", s.GetSongTextByName)
	mux.HandleFunc("PATCH /song/{id}", s.UpdateSong)
	mux.HandleFunc("PUT /song/{id}", s.UpdateSong)
	mux.HandleFunc("DELETE /song/{id}", s.DeleteSong)
	mux.HandleFunc("GET /song/all", s.GetAllSongs)
}
