package sqlite

import (
	"database/sql"
	_ "errors"
	_ "github.com/mattn/go-sqlite3"
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
