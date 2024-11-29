package song

import (
	"bytes"
	"effectivemobiletesttask/internal/config"
	"effectivemobiletesttask/internal/domain/models"
	lg "effectivemobiletesttask/internal/utils/logger"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type Client struct {
	log *slog.Logger
	api config.APIClient
}

func NewClient(log *slog.Logger, api config.APIClient) *Client {
	return &Client{
		log: log,
		api: api,
	}
}

func (c *Client) GetSongDetail(songReq models.SongRequest) (models.SongDetail, error) {
	const op = "client.song.GetSongDetail"

	reqBody, err := json.Marshal(songReq)
	if err != nil {
		c.log.Error("error marshaling song request", lg.Err(err))
		return models.SongDetail{}, fmt.Errorf("%s: %w", op, err)
	}

	c.log.Debug("start preparing request")
	req, err := http.NewRequest("GET", c.api.Protocol+"://"+c.api.Address+c.api.Url, bytes.NewReader(reqBody))
	if err != nil {
		c.log.Error("error creating the HTTP request", lg.Err(err))
		return models.SongDetail{}, fmt.Errorf("%s: %w", op, err)
	}
	req.Header.Set("Content-Type", "application/json")
	c.log.Debug("successfully prepared request")

	c.log.Debug("start sending request to API")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.log.Error("error during the request to API", lg.Err(err))
		return models.SongDetail{}, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()
	c.log.Debug("successfully sent request to API")

	c.log.Debug("start fetching song detail")
	var songDetail models.SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		c.log.Error("error during the fetching song detail")
		return models.SongDetail{}, fmt.Errorf("%s: %w", op, err)
	}
	c.log.Debug("fetched song detail: ", slog.Any("detail", songDetail))

	return songDetail, nil
}
