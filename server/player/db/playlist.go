package db

import (
	"fmt"
	"strings"
)

type Playlist struct {
	Name string
}

func (d *Db) GetPlaylist(name string) (Playlist, error) {
	var playlist Playlist
	err := d.db.QueryRow(
		`SELECT name FROM playlist WHERE name = ?`,
		name,
	).Scan(
		&playlist.Name,
	)
	return playlist, err
}

func (d *Db) GetPlaylists() ([]Playlist, error) {
	rows, err := d.db.Query(
		`SELECT name FROM playlist`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var playlists []Playlist
	for rows.Next() {
		var playlist Playlist
		err := rows.Scan(
			&playlist.Name,
		)
		if err != nil {
			return nil, err
		}
		playlists = append(playlists, playlist)
	}
	return playlists, nil
}

func (d *Db) AddPlaylist(plname string) error {
	_, err := d.db.Exec(
		`INSERT OR IGNORE INTO playlist (name) VALUES (?)`,
		plname,
	)
	if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed") {
		return fmt.Errorf(
			"Playlist %s already exists",
			plname,
		)
	}
	return err
}

func (d *Db) RemovePlaylist(name string) error {
	_, err := d.db.Exec(
		`DELETE FROM playlist WHERE name = ?`,
		name,
	)
	_, err = d.db.Exec(
		`DELETE FROM music_playlist WHERE playlist_name = ?`,
		name,
	)
	return err
}

func (d *Db) AddMusicToPlaylist(musicName, playlistName string) error {
	_, err := d.db.Exec(
		`INSERT INTO music_playlist (music_name, playlist_name) VALUES (?, ?)`,
		musicName,
		playlistName,
	)
	if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed") {
		return fmt.Errorf(
			"Music %s already in playlist %s",
			musicName,
			playlistName,
		)
	}
	return err
}

func (d *Db) RemoveMusicFromPlaylist(
	playlistName string,
	musicName string,
) error {
	_, err := d.db.Exec(
		`DELETE FROM music_playlist WHERE music_name= ? AND playlist_name = ?`,
		musicName,
		playlistName,
	)
	return err
}

func (d *Db) GetMusicsFromPlaylist(playlistName string) ([]Music, error) {
	rows, err := d.db.Query(
		`SELECT m.name, m.source, m.key, m.data
         FROM music m
        JOIN music_playlist mp ON m.name = mp.music_name
         WHERE mp.playlist_name = ?`,
		playlistName,
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
		)
		if err != nil {
			return nil, err
		}
		musics = append(musics, music)
	}

	return musics, nil
}

func (d *Db) InitPlaylist() error {
	_, err := d.db.Exec(
		`CREATE TABLE IF NOT EXISTS playlist (
      name TEXT PRIMARY KEY
    )`,
	)
	return err
}

// relation ship
func (d *Db) InitMusicPlaylist() error {
	_, err := d.db.Exec(
		`CREATE TABLE IF NOT EXISTS music_playlist (
      music_name TEXT,
      playlist_name TEXT,
      PRIMARY KEY (music_name, playlist_name),
      FOREIGN KEY (music_name) REFERENCES music (name),
      FOREIGN KEY (playlist_name) REFERENCES playlist (name)
    )`,
	)
	return err
}
