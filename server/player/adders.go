package player

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Malwarize/goplay/logger"
	"github.com/Malwarize/goplay/server/player/db"
	"github.com/Malwarize/goplay/shared"
)

func (p *Player) AddMusicFromFile(path string) error {
	// read file
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	ok, err := p.Converter.IsMp3(data)
	if err != nil {
		return err
	}

	var mp3data []byte
	if !ok {
		if mp3data, err = p.Converter.ConvertToMP3(
			data,
		); err != nil {
			return err
		}
	} else {
		mp3data = data
	}
	music, err := NewMusic(
		filepath.Base(
			path,
		),
		mp3data,
	)
	if err != nil {
		logger.LogInfo(
			"here is the error",
			err,
		)
		return err
	}
	if music == nil {
		return fmt.Errorf(
			"can't enqueue music for some reason",
		)
	}
	logger.LogInfo(
		"Enqueue", music.Name,
	)
	p.Queue.Enqueue(*music)
	return nil
}

// this function is used to play music from a file that is not mp3/ it will convert it to mp3 in temp and add it to the queue
func (p *Player) AddMusicsFromDir(dirPath string) error {
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
			err := p.AddMusicFromFile(musicPath)
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
func (p *Player) AddMusicFromOnline(unique string, engineName string) error {
	p.addTask(unique, shared.Downloading)
	music, err := p.Director.Download(engineName, unique)
	if err != nil {
		p.errorifyTask(unique, err)
		return logger.LogError(
			logger.GError(
				"Failed to download music",
				err,
			),
		)
	}

	data, err := p.Converter.ConvertToMP3(music.Data)
	if err != nil {
		p.errorifyTask(unique, err)
		return logger.LogError(
			logger.GError(
				"Failed to convert music to mp3",
				err,
			),
		)
	}

	if len(data) == 0 {
		p.errorifyTask(unique, fmt.Errorf("failed to download music: %s", unique))
		return logger.LogError(
			logger.GError(
				"Failed to download music",
				fmt.Errorf(
					"failed to download music: %s",
					unique,
				),
			),
		)
	}

	err = p.EnqueueDbMusic(
		music,
	)
	if err != nil {
		return logger.LogError(
			logger.GError(
				"Failed to enqueue music",
				fmt.Errorf(
					"failed to enqueue music: %s",
					unique,
				),
				err,
			),
		)
	}
	p.removeTask(unique)
	return nil
}

//
// func (p *Player) addMusicFromPlaylistByIndex(pl *PlayList, index int) error {
// 	playlistPath := filepath.Join(
// 		config.GetConfig().PlaylistPath,
// 		pl.Name,
// 	)
// 	dir, err := os.Open(
// 		playlistPath,
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	entries, err := dir.Readdir(0)
// 	if err != nil {
// 		return err
// 	}
// 	if index < len(entries) && index >= 0 {
// 		err := p.AddMusicFromFile(
// 			filepath.Join(playlistPath, entries[index].Name()),
// 		)
// 		return err
// 	}
// 	return nil
// }
//
//
// func (p *Player) addMusicFromPlaylistByName(
// 	pl *PlayList,
// 	name string,
// ) error {
// 	playlistPath := filepath.Join(
// 		config.GetConfig().PlaylistPath,
// 		pl.Name,
// 	)
// 	dir, err := os.Open(
// 		playlistPath,
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	entries, err := dir.Readdir(0)
// 	if err != nil {
// 		return err
// 	}
// 	for _, entry := range entries {
// 		if shared.ViewParseName(entry.Name()) == name {
// 			err := p.AddMusicFromFile(
// 				filepath.Join(playlistPath, entry.Name()),
// 			)
// 			return err
// 		}
// 	}
// 	return os.ErrNotExist
// }
//
// func (p *Player) addMusicsFromPlaylist(pl *PlayList) error {
// 	playlistPath := filepath.Join(config.GetConfig().PlaylistPath, pl.Name)
// 	dir, err := os.Open(playlistPath)
// 	if err != nil {
// 		return err
// 	}
// 	entries, err := dir.Readdir(0)
// 	if err != nil {
// 		return err
// 	}
// 	for _, entry := range entries {
// 		if !entry.IsDir() {
// 			err := p.AddMusicFromFile(filepath.Join(playlistPath, entry.Name()))
// 			if err != nil {
// 				logger.LogWarn(err.Error())
// 			}
// 		}
// 	}
// 	return nil
// }
//
// wrapper for AddMusicsFromPlaylist() + validate playlist
// func (p *Player) AddMusicsFromPlaylist(
// 	plname string,
// ) error {
// 	pl, err := p.PlayListManager.GetPlayListByName(plname)
// 	if err != nil {
// 		return logger.LogError(
// 			logger.GError(
// 				"Failed to get playlist by name",
// 				err,
// 			),
// 		)
// 	}
// 	return p.addMusicsFromPlaylist(pl)
// }

func (p *Player) EnqueueDbMusic(dmusic *db.Music) error {
	playableMusic, err := NewMusic(
		dmusic.Name,
		dmusic.Data,
	)
	if err != nil {
		return err
	}
	p.Queue.Enqueue(
		*playableMusic,
	)
	return nil
}
