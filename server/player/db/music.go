package db

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Music struct {
	Name   string
	Source string
	Key    string
	Data   []byte
	_hash  string
}

func (d *Db) InitMusic() error {
	_, err := d.db.Exec(
		`CREATE TABLE IF NOT EXISTS music (
      name TEXT UNIQUE,
      source TEXT,
      key TEXT,
      data BLOB,
      hash TEXT UNIQUE NOT NULL,
      PRIMARY KEY (source, key)
    )`,
	)
	return err
}

func (d *Db) GetMusic(source string, key string) (Music, error) {
	var music Music
	err := d.db.QueryRow(
		`SELECT name, source, key, data, hash FROM music WHERE source = ? AND key = ?`,
		source,
		key,
	).Scan(
		&music.Name,
		&music.Source,
		&music.Key,
		&music.Data,
		&music._hash,
	)
	return music, err
}

func (d *Db) UpdateMusic(
	name string,
	source string,
	key string,
	data []byte,
) error {
	_, err := d.db.Exec(
		`UPDATE music SET data = ?, hash = ? WHERE name = ? AND source = ? AND key = ?`,
		data,
		hash(data),
		name,
		source,
		key,
	)
	return err
}

// insertUniqueMusicName is a helper function to insert music with a unique name
func (d *Db) insertUniqueMusicName(music *Music) error {
	var newName string
	err := d.db.QueryRow(
		`SELECT name FROM music WHERE name LIKE ? ORDER BY name DESC LIMIT 1`,
		fmt.Sprintf("%s%%", music.Name),
	).Scan(&newName)
	if err == sql.ErrNoRows {
		// No similar names found, append _1
		newName = fmt.Sprintf("%s_1", music.Name)
	} else if err != nil {
		return err
	} else {
		// Extract the current numeric suffix and increment
		var suffix int
		_, err := fmt.Sscanf(newName, "%s_%d", music.Name, &suffix)
		if err != nil {
			return err
		}
		newName = fmt.Sprintf("%s_%d", music.Name, suffix+1)
	}

	// Insert music with the new unique name
	_, err = d.db.Exec(
		`INSERT INTO music (name, source, key, data, hash) VALUES (?, ?, ?, ?,?)`,
		newName,
		music.Source,
		music.Key,
		music.Data,
		hash(music.Data),
	)
	return err
}

func (d *Db) AddMusic(music *Music) error {
	// Check if the name is already used

	var count int
	err := d.db.QueryRow(
		`SELECT COUNT(*) FROM music WHERE name = ?`,
		music.Name,
	).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		// If the name is already used, increment it
		return d.insertUniqueMusicName(music)
	}

	// If the name is not used, insert the music with hash
	_, err = d.db.Exec(
		`INSERT INTO music (name, source, key, data, hash) VALUES (?, ?, ?, ?, ?)`,
		music.Name,
		music.Source,
		music.Key,
		music.Data,
		hash(music.Data),
	)
	return err
}

func (d *Db) GetMusicByName(name string) (Music, error) {
	var music Music
	err := d.db.QueryRow(
		`SELECT name, source, key, data, hash FROM music WHERE name = ?`,
		name,
	).Scan(
		&music.Name,
		&music.Source,
		&music.Key,
		&music.Data,
		&music._hash,
	)
	return music, err
}

func (d *Db) GetMusicByKeySource(source string, key string) (Music, error) {
	var music Music
	err := d.db.QueryRow(
		`SELECT name, source, key, data, hash FROM music WHERE source = ? AND key = ?`,
		source,
		key,
	).Scan(
		&music.Name,
		&music.Source,
		&music.Key,
		&music.Data,
		&music._hash,
	)
	return music, err
}

func hash(data []byte) string {
	hasher := md5.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}

func (m Music) GetHash() string {
	return m._hash
}

func (d *Db) GetMusicByHash(hash string) (Music, error) {
	var music Music
	err := d.db.QueryRow(
		`SELECT name, source, key, data, hash FROM music WHERE hash = ?`,
		hash,
	).Scan(
		&music.Name,
		&music.Source,
		&music.Key,
		&music.Data,
		&music._hash,
	)
	return music, err
}

func (d *Db) FilterMusic(query string) ([]Music, error) {
	rows, err := d.db.Query(
		`SELECT name, source, key, data, hash FROM music WHERE name LIKE ?`,
		fmt.Sprintf("%%%s%%", query),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var musics []Music
	for rows.Next() {
		var music Music
		err := rows.Scan(
			&music.Name,
			&music.Source,
			&music.Key,
			&music.Data,
			&music._hash,
		)
		if err != nil {
			return nil, err
		}
		musics = append(musics, music)
	}
	return musics, nil
}

func (d *Db) CleanCache() error {
	// Delete music thats not in playlist
	_, err := d.db.Exec(
		`DELETE FROM music WHERE hash NOT IN (SELECT hash FROM playlist)`,
	)
	return err
}

func (d *Db) GetCachedMusics() ([]Music, error) {
	rows, err := d.db.Query(
		`SELECT * FROM music where music.hash NOT IN  (SELECT hash FROM music_playlist)`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var musics []Music
	for rows.Next() {
		var music Music
		err := rows.Scan(
			&music.Name,
			&music.Source,
			&music.Key,
			&music.Data,
			&music._hash,
		)
		if err != nil {
			return nil, err
		}
		musics = append(musics, music)
	}
	return musics, nil
}
