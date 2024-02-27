package db

import (
	_ "github.com/mattn/go-sqlite3"
)

type Music struct {
	Name   string
	Source string
	Key    string
	Data   []byte
}

func (d *Db) InitMusic() error {
	_, err := d.db.Exec(
		`CREATE TABLE IF NOT EXISTS music (
      name TEXT,
      source TEXT,
      key TEXT,
      data BLOB,
      PRIMARY KEY (source, key)
    )`,
	)
	return err
}

func (d *Db) AddMusic(music *Music) error {
	_, err := d.db.Exec(
		`INSERT INTO music (name, source, key, data) VALUES (?, ?, ?, ?)`,
		music.Name,
		music.Source,
		music.Key,
		music.Data,
	)
	return err
}

func (d *Db) GetMusic(source string, key string) (Music, error) {
	var music Music
	err := d.db.QueryRow(
		`SELECT name, source, key, data FROM music WHERE source = ? AND key = ?`,
		source,
		key,
	).Scan(
		&music.Name,
		&music.Source,
		&music.Key,
		&music.Data,
	)
	return music, err
}
