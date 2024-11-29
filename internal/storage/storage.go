package storage

import "errors"

var (
	ErrSongNotFound  = errors.New("song was not found")
	ErrGroupNotFound = errors.New("group was not found")
)
