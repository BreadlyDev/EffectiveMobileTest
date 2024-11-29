package song

import (
	"effectivemobiletesttask/internal/domain/models"
	srv "effectivemobiletesttask/internal/http-server"
	"effectivemobiletesttask/internal/storage"
	jsn "effectivemobiletesttask/internal/utils/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// CreateSong adds a new song to the library.
// @Summary Add a new song
// @Description Add a new song to the music library
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.SongRequest true "Song Request"
// @Success 201 {object} httpserver.Response
// @Failure 400 {object} httpserver.Response
// @Failure 500 {object} httpserver.Response
// @Router /song/create [post]
func (s *Server) CreateSong(w http.ResponseWriter, r *http.Request) {
	var songReq models.SongRequest
	var resp srv.Response

	if err := jsn.ReadRequestBody(r, &songReq); err != nil {
		resp = srv.NewErrResponse(srv.ErrBadRequest.Error(), http.StatusBadRequest)

		jsn.WriteResponseBody(w, resp, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	field, err := srv.ValidateSongRequest(songReq)
	if err != nil {
		resp = srv.NewErrResponse(fmt.Sprintf("'%s' %s", field, err.Error()), http.StatusBadRequest)

		jsn.WriteResponseBody(w, resp, http.StatusBadRequest)
		return
	}

	id, err := s.service.CreateSong(songReq)
	if err != nil {
		resp = srv.NewErrResponse(srv.ErrInternalServer.Error(), http.StatusInternalServerError)

		jsn.WriteResponseBody(w, resp, http.StatusInternalServerError)
		return
	}

	resp = srv.NewResponse("Added new song", http.StatusCreated, id)

	jsn.WriteResponseBody(w, resp, http.StatusCreated)
}

// GetSongByID retrieves a song by its ID.
// @Summary Get song by ID
// @Description Get a specific song by its unique ID
// @Tags songs
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} httpserver.Response
// @Failure 400 {object} httpserver.Response
// @Failure 404 {object} httpserver.Response
// @Failure 500 {object} httpserver.Response
// @Router /song/{id} [get]
func (s *Server) GetSongByID(w http.ResponseWriter, r *http.Request) {
	idStr := srv.GetPathParameter(r, 2)
	id, err := strconv.ParseInt(idStr, 10, 64)

	var resp srv.Response

	if err != nil {
		resp = srv.NewErrResponse(srv.ErrWrongPathParameter.Error(), http.StatusBadRequest)

		jsn.WriteResponseBody(w, resp, http.StatusBadRequest)
		return
	}

	song, err := s.service.GetSongByID(id)
	if err != nil {
		if errors.Is(err, storage.ErrSongNotFound) {
			resp = srv.NewErrResponse("Song was not found", http.StatusNotFound)

			jsn.WriteResponseBody(w, resp, http.StatusNotFound)
			return
		}

		resp = srv.NewErrResponse(srv.ErrInternalServer.Error(), http.StatusInternalServerError)

		jsn.WriteResponseBody(w, resp, http.StatusInternalServerError)
		return
	}

	resp = srv.NewResponse("Successfully fetched song", http.StatusOK, song)

	jsn.WriteResponseBody(w, resp, http.StatusOK)
}

// GetSongByName retrieves a song by its name.
// @Summary Get song by name
// @Description Get a specific song by its name
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.SongName true "Song Name"
// @Success 200 {object} httpserver.Response
// @Failure 400 {object} httpserver.Response
// @Failure 404 {object} httpserver.Response
// @Failure 500 {object} httpserver.Response
// @Router /song/name [get]
func (s *Server) GetSongByName(w http.ResponseWriter, r *http.Request) {
	var songReq models.SongRequest
	var resp srv.Response

	if err := jsn.ReadRequestBody(r, &songReq); err != nil {
		resp = srv.NewErrResponse(srv.ErrBadRequest.Error(), http.StatusBadRequest)

		jsn.WriteResponseBody(w, resp, http.StatusBadRequest)
		return
	}

	if songReq.Name == "" {
		resp = srv.NewErrResponse(fmt.Sprintf("'song' %s", srv.ErrFieldIsRequired.Error()), http.StatusBadRequest)

		jsn.WriteResponseBody(w, resp, http.StatusBadRequest)
		return
	}

	song, err := s.service.GetSongByName(songReq.Name)
	if err != nil {
		if errors.Is(err, storage.ErrSongNotFound) {
			resp = srv.NewErrResponse("Song was not found", http.StatusNotFound)

			jsn.WriteResponseBody(w, resp, http.StatusNotFound)
			return
		}

		resp = srv.NewErrResponse(srv.ErrInternalServer.Error(), http.StatusInternalServerError)

		jsn.WriteResponseBody(w, resp, http.StatusInternalServerError)
		return
	}

	resp = srv.NewResponse("Successfully fetched song", http.StatusOK, song)

	jsn.WriteResponseBody(w, resp, http.StatusOK)
}

// GetSongTextByName retrieves the text of a song by its name.
// @Summary Get song text by name
// @Description Fetch the text of a specific song using its name
// @Tags songs
// @Accept json
// @Produce json
// @Param name query string false "Name"
// @Param verse query int false "Verse"
// @Success 200 {object} httpserver.Response
// @Failure 400 {object} httpserver.Response
// @Failure 404 {object} httpserver.Response
// @Failure 500 {object} httpserver.Response
// @Router /song/name/text [get]
func (s *Server) GetSongTextByName(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	name := params.Get("name")
	verseParam := params.Get("verse")

	verse := 0

	if verseParam != "" {
		if parsedVerse, err := strconv.Atoi(verseParam); err == nil {
			verse = parsedVerse
		}
	}

	var resp srv.Response

	text, err := s.service.GetSongTextByName(name, verse)
	if err != nil {
		if errors.Is(err, storage.ErrSongNotFound) {
			resp = srv.NewErrResponse("Song was not found", http.StatusNotFound)

			jsn.WriteResponseBody(w, resp, http.StatusNotFound)
			return
		}

		resp = srv.NewErrResponse(srv.ErrInternalServer.Error(), http.StatusInternalServerError)

		jsn.WriteResponseBody(w, resp, http.StatusInternalServerError)
		return
	}

	resp = srv.NewResponse("Successfully fetched song text", http.StatusOK, text)

	jsn.WriteResponseBody(w, resp, http.StatusInternalServerError)
}

// GetSongTextByID retrieves the text of a song by its ID.
// @Summary Get song text by ID
// @Description Fetch the text of a specific song using its ID
// @Tags songs
// @Produce json
// @Param id path int true "Song ID"
// @Param verse query int false "Verse"
// @Success 200 {object} httpserver.Response
// @Failure 400 {object} httpserver.Response
// @Failure 404 {object} httpserver.Response
// @Failure 500 {object} httpserver.Response
// @Router /song/{id}/text [get]
func (s *Server) GetSongTextByID(w http.ResponseWriter, r *http.Request) {
	verseParam := r.URL.Query().Get("verse")
	verse := 0

	if verseParam != "" {
		if parsedVerse, err := strconv.Atoi(verseParam); err == nil {
			verse = parsedVerse
		}
	}

	idStr := srv.GetPathParameter(r, 2)
	id, err := strconv.ParseInt(idStr, 10, 64)

	var resp srv.Response

	if err != nil {
		resp = srv.NewErrResponse(srv.ErrWrongPathParameter.Error(), http.StatusBadRequest)

		jsn.WriteResponseBody(w, resp, http.StatusBadRequest)
		return
	}

	text, err := s.service.GetSongTextByID(id, verse)
	if err != nil {
		if errors.Is(err, storage.ErrSongNotFound) {
			resp = srv.NewErrResponse("Song was not found", http.StatusNotFound)

			jsn.WriteResponseBody(w, resp, http.StatusNotFound)
			return
		}

		resp = srv.NewErrResponse(srv.ErrInternalServer.Error(), http.StatusInternalServerError)

		jsn.WriteResponseBody(w, resp, http.StatusInternalServerError)
		return
	}

	resp = srv.NewResponse("Successfully fetched song", http.StatusOK, text)

	jsn.WriteResponseBody(w, resp, http.StatusOK)
}

// UpdateSong updates details of an existing song.
// @Summary Update song
// @Description Update the details of a song by its ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body models.SongResponse true "Updated Song Data"
// @Success 200 {object} httpserver.Response
// @Failure 400 {object} httpserver.Response
// @Failure 404 {object} httpserver.Response
// @Failure 500 {object} httpserver.Response
// @Router /song/{id} [put]
// @Router /song/{id} [patch]
func (s *Server) UpdateSong(w http.ResponseWriter, r *http.Request) {
	idStr := srv.GetPathParameter(r, 2)
	id, err := strconv.ParseInt(idStr, 10, 64)

	var resp srv.Response

	if err != nil {
		resp = srv.NewErrResponse(srv.ErrWrongPathParameter.Error(), http.StatusBadRequest)

		jsn.WriteResponseBody(w, resp, http.StatusBadRequest)
		return
	}

	var newSong models.Song

	if err := jsn.ReadRequestBody(r, &newSong); err != nil {
		resp = srv.NewErrResponse("Error during decoding JSON", http.StatusBadRequest)

		jsn.WriteResponseBody(w, resp, http.StatusBadRequest)
		return
	}

	releaseDate, err := srv.ParseReleaseDate(newSong.ReleaseDate)
	if err != nil {
		resp = srv.NewErrResponse("Error during parsing release date", http.StatusBadRequest)

		jsn.WriteResponseBody(w, resp, http.StatusBadRequest)
		return
	}

	songResp := srv.SongToSongResponse(newSong, releaseDate)

	song, err := s.service.UpdateSong(id, songResp)
	if err != nil {
		if errors.Is(err, storage.ErrSongNotFound) {
			resp = srv.NewErrResponse("Song was not found", http.StatusNotFound)

			jsn.WriteResponseBody(w, resp, http.StatusNotFound)
			return
		}

		resp = srv.NewErrResponse("Error during fetching song", http.StatusBadRequest)

		jsn.WriteResponseBody(w, resp, http.StatusBadRequest)
		return
	}

	resp = srv.NewResponse("Successfully updated song", http.StatusOK, song)

	jsn.WriteResponseBody(w, resp, http.StatusOK)
}

// DeleteSong removes a song from the library.
// @Summary Delete song
// @Description Delete a specific song by its ID
// @Tags songs
// @Param id path int true "Song ID"
// @Success 204 {object} httpserver.Response
// @Failure 400 {object} httpserver.Response
// @Failure 404 {object} httpserver.Response
// @Failure 500 {object} httpserver.Response
// @Router /song/{id} [delete]
func (s *Server) DeleteSong(w http.ResponseWriter, r *http.Request) {
	idStr := srv.GetPathParameter(r, 2)
	id, err := strconv.ParseInt(idStr, 10, 64)

	var resp srv.Response

	if err != nil {
		resp = srv.NewErrResponse(srv.ErrWrongPathParameter.Error(), http.StatusBadRequest)

		jsn.WriteResponseBody(w, resp, http.StatusBadRequest)
		return
	}

	err = s.service.DeleteSong(id)
	if err != nil {
		if errors.Is(err, storage.ErrSongNotFound) {
			resp = srv.NewErrResponse("Song was not found", http.StatusNotFound)

			jsn.WriteResponseBody(w, resp, http.StatusNotFound)
			return
		}

		resp = srv.NewErrResponse("Error during fetching song", http.StatusBadRequest)

		jsn.WriteResponseBody(w, resp, http.StatusBadRequest)
		return
	}

	resp = srv.NewResponse("Successfully deleted song", http.StatusNoContent, nil)

	jsn.WriteResponseBody(w, resp, http.StatusNoContent)
}

// GetAllSongs retrieves songs matching specific filters.
// @Summary Get all songs
// @Description Fetch songs that match the provided filters
// @Tags songs
// @Produce json
// @Param group query string false "Song group"
// @Param name query string false "Song name"
// @Param text query string false "Song text"
// @Param link query string false "External link"
// @Param releaseDate query string false "Release date (YYYY-MM-DD)"
// @Param page query int false "Page number"
// @Success 200 {object} httpserver.Response
// @Failure 400 {object} httpserver.Response
// @Failure 500 {object} httpserver.Response
// @Router /song/all [get]
func (s *Server) GetAllSongs(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	var resp srv.Response

	var filter models.SongFilter
	filter.Group = params.Get("group")
	filter.Name = params.Get("name")
	filter.Text = params.Get("text")
	filter.Link = params.Get("link")
	filter.ReleaseDate = time.Time{}

	rlsDateParam := params.Get("releaseDate")

	if rlsDateParam != "" {
		rlsDate, err := srv.ParseReleaseDate(rlsDateParam)

		if err != nil {
			resp = srv.NewErrResponse(err.Error(), http.StatusBadRequest)

			jsn.WriteResponseBody(w, resp, http.StatusBadRequest)
			return
		}

		filter.ReleaseDate = rlsDate
	}

	pageParam := params.Get("page")
	page := 0

	if pageParam != "" {
		if parsedPage, err := strconv.Atoi(pageParam); err == nil {
			page = parsedPage
		}
	}

	songs, err := s.service.GetAllSongs(filter, s.pageSize*page, s.pageSize)
	if err != nil {
		if errors.Is(err, storage.ErrGroupNotFound) {
			resp = srv.NewErrResponse("Group was not found", http.StatusBadRequest)

			jsn.WriteResponseBody(w, resp, http.StatusBadRequest)
			return
		}

		resp = srv.NewErrResponse("Error during fetching songs", http.StatusBadRequest)

		jsn.WriteResponseBody(w, resp, http.StatusBadRequest)
		return
	}

	resp = srv.NewResponse("Successfully fetched songs", http.StatusOK, songs)

	jsn.WriteResponseBody(w, resp, http.StatusOK)
}
