package httpserver

import (
	"effectivemobiletesttask/internal/domain/models"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	ErrBadRequest         = errors.New("Bad request")
	ErrInternalServer     = errors.New("Internal server error")
	ErrWrongPathParameter = errors.New("Wrong path parameter")
	ErrFieldIsRequired    = errors.New("field is required")
)

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    any    `json:"data,omitempty"`
}

func NewResponse(message string, status int, data any) Response {
	return Response{
		Message: message,
		Status:  status,
		Data:    data,
	}
}

func NewErrResponse(message string, status int) Response {
	return NewResponse(message, status, nil)
}

func GetPathParameter(r *http.Request, pNum int) string {
	params := strings.Split(r.URL.Path, "/")

	if pNum <= 0 {
		pNum = 1
	}

	if len(params) == 1 {
		return ""
	} else {
		return params[pNum]
	}
}

func ValidateSongRequest(songReq models.SongRequest) (string, error) {
	if songReq.Name == "" {
		return "song", ErrFieldIsRequired
	}

	if songReq.Group == "" {
		return "group", ErrFieldIsRequired
	}

	return "", nil
}

func ParseReleaseDate(rlsDateStr string) (time.Time, error) {
	const layout = "2006-01-02"

	releaseDate, err := time.Parse(layout, rlsDateStr)

	if err != nil {
		log.Printf("Invalid releaseDate format: %v", err)

		return time.Time{}, errors.New("invalid releaseDate format. use YYYY-MM-DD")
	}

	return releaseDate, nil
}

func SongToSongResponse(song models.Song, releaseDate time.Time) models.SongResponse {
	var songResp models.SongResponse

	songResp.ID = song.ID
	songResp.Group = song.Group
	songResp.Name = song.Name
	songResp.ReleaseDate = releaseDate
	songResp.Text = song.Text
	songResp.Link = song.Link

	return songResp
}
