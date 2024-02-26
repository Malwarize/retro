package player

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Malwarize/goplay/config"
	"github.com/Malwarize/goplay/logger"
	"github.com/Malwarize/goplay/shared"
)

func (p *Player) addConvertedMp3InTemp(path string) bool {
	f, err := os.CreateTemp("", "goplay")
	defer os.Remove(f.Name())
	if err != nil {
		logger.LogError(
			logger.GError(
				"Failed to create temp file",
				err,
			),
		)

		return false
	}
	sourceFile, err := os.Open(path)
	if err != nil {
		logger.ERRORLogger.Println(err)
		return false
	}
	_, err = io.Copy(f, sourceFile)
	if err != nil {
		logger.ERRORLogger.Println(err)
		return false
	}
	err = p.Converter.ConvertToMP3(f.Name())
	if err != nil {
		logger.ERRORLogger.Println(err)
		return false
	}
	err = p.AddMusicFromFile(f.Name())
	if err != nil {
		return false
	}
	return true
}

func (p *Player) AddMusicFromFile(path string) error {
	music, err := NewMusic(path)
	if err != nil {
		return err
	}
	p.Queue.Enqueue(music)
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
			isMp3, err := p.Converter.IsMp3(dirPath + "/" + entry.Name())
			if err != nil {
				logger.WARNLogger.Println(
					"Failed to check if file is mp3",
					err,
				)
				continue
			}
			if isMp3 {
				logger.INFOLogger.Println("Playing music from dir", dirPath+"/"+entry.Name())
				err := p.AddMusicFromFile(filepath.Join(dirPath, entry.Name()))
				if err != nil {
					logger.WARNLogger.Println(
						"Failed to add music from dir",
						err,
					)
				}
			} else {
				if !p.addConvertedMp3InTemp(filepath.Join(dirPath, entry.Name())) {
					return logger.LogError(
						logger.GError(
							"file is not mp3 and failed to convert to mp3",
						),
					)
				}
			}
		}
	}
	return nil
}

// the unique is the unique id of the music in the engine it can be url or id
func (p *Player) AddMusicFromOnline(unique string, engineName string) error {
	p.addTask(unique, shared.Downloading)
	path, err := p.Director.Download(engineName, unique)
	if err != nil {
		p.errorifyTask(unique, err)
		return logger.LogError(
			logger.GError(
				"Failed to download music",
				err,
			),
		)
	}

	err = p.Converter.ConvertToMP3(path)
	if err != nil {
		p.errorifyTask(unique, err)
		return logger.LogError(
			logger.GError(
				"Failed to convert music to mp3",
				err,
			),
		)
	}

	if path == "" {
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
	err = p.AddMusicFromFile(path)
	if err != nil {
		return logger.LogError(
			logger.GError(
				"Failed to add music to queue",
				err,
			),
		)
	}
	p.removeTask(unique)
	return nil
}

func (p *Player) addMusicFromPlaylistByIndex(pl *PlayList, index int) error {
	playlistPath := filepath.Join(
		config.GetConfig().PlaylistPath,
		pl.Name,
	)
	dir, err := os.Open(
		playlistPath,
	)
	if err != nil {
		return err
	}
	entries, err := dir.Readdir(0)
	if err != nil {
		return err
	}
	if index < len(entries) && index >= 0 {
		err := p.AddMusicFromFile(
			filepath.Join(playlistPath, entries[index].Name()),
		)
		return err
	}
	return nil
}

func (p *Player) addMusicFromPlaylistByName(
	pl *PlayList,
	name string,
) error {
	playlistPath := filepath.Join(
		config.GetConfig().PlaylistPath,
		pl.Name,
	)
	dir, err := os.Open(
		playlistPath,
	)
	if err != nil {
		return err
	}
	entries, err := dir.Readdir(0)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if shared.ViewParseName(entry.Name()) == name {
			err := p.AddMusicFromFile(
				filepath.Join(playlistPath, entry.Name()),
			)
			return err
		}
	}
	return os.ErrNotExist
}

func (p *Player) addMusicsFromPlaylist(pl *PlayList) error {
	playlistPath := filepath.Join(config.GetConfig().PlaylistPath, pl.Name)
	dir, err := os.Open(playlistPath)
	if err != nil {
		return err
	}
	entries, err := dir.Readdir(0)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			err := p.AddMusicFromFile(filepath.Join(playlistPath, entry.Name()))
			if err != nil {
				logger.LogWarn(err.Error())
			}
		}
	}
	return nil
}

// wrapper for AddMusicsFromPlaylist() + validate playlist
func (p *Player) AddMusicsFromPlaylist(
	plname string,
) error {
	pl, err := p.PlayListManager.GetPlayListByName(plname)
	if err != nil {
		return logger.LogError(
			logger.GError(
				"Failed to get playlist by name",
				err,
			),
		)
	}
	return p.addMusicsFromPlaylist(pl)
}
