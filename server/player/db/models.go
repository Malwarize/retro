package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	db   *sql.DB
	path string
}

func NewDb(
	path string,
) (*Db, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return &Db{
		db:   db,
		path: path,
	}, nil
}

func (d *Db) Close() error {
	return d.db.Close()
}

func LoadDb(
	path string,
) (*Db, error) {
	db, err := NewDb(path)
	if err != nil {
		return nil, err
	}
	err = db.InitMusic()
	if err != nil {
		return nil, err
	}

	return db, nil
}
