package models

import "time"

type SongName struct {
	Name string `json:"song"`
}

type SongRequest struct {
	Group string `json:"group"`
	Name  string `json:"song"`
}

type SongDetail struct {
	ReleaseDate time.Time `json:"releaseDate"`
	Text        string    `json:"text,omitempty"`
	Link        string    `josn:"link,omitempty"`
}

type SongResponse struct {
	ID int64 `json:"id"`
	SongRequest
	SongDetail
}

type SongFilter struct {
	SongRequest
	SongDetail
}

type SongStorage struct {
	ID      int64
	GroupID int64
	Name    string
	SongDetail
}

type Song struct {
	ID int64 `json:"id"`
	SongRequest
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text,omitempty"`
	Link        string `josn:"link,omitempty"`
}

type SongUpdate struct {
	Name        string `json:"song"`
	Group       string `json:"group"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text,omitempty"`
	Link        string `josn:"link,omitempty"`
}
