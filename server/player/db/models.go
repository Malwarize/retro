package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
)

type Db struct {
	db   *sql.DB
	path string
}

func NewDb(path string) (*Db, error) {
	// create one if not exists mkdirAll
	fmt.Println(path)
	dir := filepath.Dir(
		filepath.Dir(path),
	)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}
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

	err = db.InitPlaylist()
	if err != nil {
		return nil, err
	}

	err = db.InitMusicPlaylist()
	if err != nil {
		return nil, err
	}

	return db, nil
}
