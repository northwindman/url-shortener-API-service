package sqlite

import (
	"database/sql"
	"errors"
	"github.com/mattn/go-sqlite3"
	"github.com/northwindman/url-shortener-API-service/internal/storage"
	"github.com/northwindman/url-shortener-API-service/pkg/lib/e"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url (
	    id INTEGER PRIMARY KEY,
	    alias TEXT NOT NULL UNIQUE,
	    url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS url_alias ON url (alias);
	`)

	if err != nil {
		return nil, e.Wrap(op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url (url, alias) VALUES (?, ?)")
	if err != nil {
		return 0, e.Wrap(op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, e.Wrap(op, storage.ErrURLExists)
		}

		return 0, e.Wrap(op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, e.Wrap(op+"failed to get last insert id: ", err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", e.Wrap(op, err)
	}

	var url string
	err = stmt.QueryRow(alias).Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", e.Wrap(op, storage.ErrURLNotFound)
		}

		return "", e.Wrap(op, err)
	}

	return url, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.sqlite.DeleteURL"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = ?")
	if err != nil {
		return e.Wrap(op, err)
	}

	_, err = stmt.Exec(alias)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.Wrap(op, storage.ErrURLNotFound)
		}

		return e.Wrap(op, err)
	}

	return nil
}
