package postgres

import (
	"database/sql"
	"effectivemobiletesttask/internal/domain/models"
	"effectivemobiletesttask/internal/storage"
	"errors"
	"fmt"
	"strings"
)

func (s *Storage) CreateSong(song models.SongStorage) (int64, error) {
	const op = "storage.postgres.CreateSong"

	stmt, err := s.db.Prepare(
		"INSERT INTO songs(group_id, name, release_date, text, link) VALUES ($1, $2, $3, $4, $5) RETURNING id",
	)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var id int64

	err = stmt.QueryRow(song.GroupID, song.Name, song.ReleaseDate, song.Text, song.Link).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetSongByID(id int64) (models.SongStorage, error) {
	const op = "storage.postgres.GetSongByID"

	stmt, err := s.db.Prepare("SELECT * FROM songs WHERE id = $1")
	if err != nil {
		return models.SongStorage{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)

	var song models.SongStorage

	err = row.Scan(&song.ID, &song.Name, &song.GroupID, &song.ReleaseDate, &song.Text, &song.Link)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.SongStorage{}, storage.ErrSongNotFound
		}

		return models.SongStorage{}, fmt.Errorf("%s: %w", op, err)
	}

	return song, nil
}

func (s *Storage) GetSongByName(songName string) (models.SongStorage, error) {
	const op = "storage.postgres.GetSongByName"

	stmt, err := s.db.Prepare("SELECT * FROM songs WHERE name = $1")
	if err != nil {
		return models.SongStorage{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(songName)

	var song models.SongStorage

	err = row.Scan(&song.ID, &song.Name, &song.GroupID, &song.ReleaseDate, &song.Text, &song.Link)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.SongStorage{}, storage.ErrSongNotFound
		}

		return models.SongStorage{}, fmt.Errorf("%s: %w", op, err)
	}

	return song, nil
}

func (s *Storage) UpdateSong(id int64, song models.SongStorage) (models.SongStorage, error) {
	const op = "storage.postgres.UpdateSong"

	stmt, err := s.db.Prepare(
		"UPDATE songs SET name = $1, group_id = $2, release_date = $3, text = $4, link = $5 WHERE id = $6",
	)
	if err != nil {
		return models.SongStorage{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(song.Name, song.GroupID, song.ReleaseDate, song.Text, song.Link, id)
	if err != nil {
		return models.SongStorage{}, fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.SongStorage{}, fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return models.SongStorage{}, storage.ErrSongNotFound
	}

	return song, nil
}

func (s *Storage) DeleteSong(id int64) error {
	const op = "storage.postgres.DeleteSong"

	stmt, err := s.db.Prepare("DELETE FROM songs WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: unable to fetch affected rows: %w", op, err)
	}

	if rowsAffected == 0 {
		return storage.ErrSongNotFound
	}

	fmt.Println("No error")

	return nil
}

func (s *Storage) GetAllSongs(filter models.SongFilter, groupID int64, offset int, limit int) ([]models.SongStorage, error) {
	const op = "storage.postgres.GetAllSongs"

	baseQuery := "SELECT * FROM songs WHERE 1=1"
	var conditions []string
	var args []interface{}

	if groupID != 0 {
		conditions = append(conditions, "group_id = $"+fmt.Sprint(len(args)+1))
		args = append(args, groupID)
	}
	if filter.Name != "" {
		conditions = append(conditions, "name = $"+fmt.Sprint(len(args)+1))
		args = append(args, filter.Name)
	}
	if !filter.ReleaseDate.IsZero() {
		conditions = append(conditions, "release_date = $"+fmt.Sprint(len(args)+1))
		args = append(args, filter.ReleaseDate.GoString())
	}
	if filter.Text != "" {
		conditions = append(conditions, "text = $"+fmt.Sprint(len(args)+1))
		args = append(args, filter.Text)
	}
	if filter.Link != "" {
		conditions = append(conditions, "link = $"+fmt.Sprint(len(args)+1))
		args = append(args, filter.Link)
	}

	args = append(args, offset, limit)

	finalQuery := baseQuery
	if len(conditions) > 0 {
		finalQuery += " AND " + strings.Join(conditions, " AND ")
	}
	finalQuery += " OFFSET $" + fmt.Sprint(len(args)-1) + " LIMIT $" + fmt.Sprint(len(args))

	rows, err := s.db.Query(finalQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var songs []models.SongStorage
	for rows.Next() {
		var song models.SongStorage
		if err := rows.Scan(&song.ID, &song.Name, &song.GroupID, &song.ReleaseDate, &song.Text, &song.Link); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		songs = append(songs, song)
	}

	return songs, nil
}
