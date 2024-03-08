package player

import (
	"os"
	"path/filepath"

	"github.com/Malwarize/retro/logger"
	"github.com/Malwarize/retro/server/player/db"
	"github.com/Malwarize/retro/shared"
)

type callback func(m db.Music) error

func (p *Player) AddMusicFromFile(
	path string,
	how callback,
) error {
	// read file
	logger.LogInfo(
		"Reading file",
	)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	// check if exists in db
	var m db.Music
	m, err = p.Director.Db.GetMusicByHash(
		hash(
			data,
		),
	)
	if err != nil {
		music, err := NewMusic(
			filepath.Base(
				path,
			),
			data,
		)
		if err != nil {
			logger.LogInfo(
				"Failed to create music",
				err,
			)
			return err
		}
		if err = p.Director.Db.AddMusic(
			&db.Music{
				Name:   music.Name,
				Data:   music.Data,
				Source: "local",
				Key:    path,
			},
		); err != nil {
			return err
		}
		m, err = p.Director.Db.GetMusicByName(music.Name)
		if err != nil {
			return logger.LogError(
				logger.GError(
					"Failed to get music from db",
					err,
				),
			)
		}
	} else {
		logger.LogWarn(
			"File Found in db",
			path,
		)
	}
	if how != nil {
		err = how(m)
		if err != nil {
			return err
		}
	}
	return nil
}

// this function is used to play music from a file that is not mp3/ it will convert it to mp3 in temp and add it to the queue
func (p *Player) AddMusicsFromDir(dirPath string, how callback) error {
	dir, err := os.Open(dirPath)
	if err != nil {
		return err
	}
	entries, err := dir.Readdir(0)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			musicPath := filepath.Join(dirPath, entry.Name())
			err := p.AddMusicFromFile(musicPath, how)
			if err != nil {
				logger.LogWarn(
					"skipping music",
					musicPath,
					"because of error",
					err,
				)
				continue
			}
		}
	}
	return nil
}

// the unique is the unique id of the music in the engine it can be url or id
func (p *Player) AddMusicFromOnline(
	unique string,
	engineName string,
	how callback,
) error {
	p.addTask(unique, shared.Downloading)
	music, err := p.Director.Download(engineName, unique)
	if err != nil || len(music.Data) == 0 {
		p.errorTask(unique, err)
		return logger.LogError(
			logger.GError(
				"Failed to download music",
				err,
			),
		)
	}

	if how != nil {
		err := how(*music)
		if err != nil {
			return err
		}
	}
	p.removeTask(unique)
	return nil
}
