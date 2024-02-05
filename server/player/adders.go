package player

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/Malwarize/goplay/config"
	"github.com/Malwarize/goplay/shared"
)

func (p *Player) AddMusicFromFile(path string) {
	music, err := NewMusic(path)
	if err != nil {
		log.Println(err)
	}
	p.AddMusicToQueue(music)
}

// this function is used to play music from a file that is not mp3/ it will convert it to mp3 in temp and add it to the queue
func (p *Player) addConvertedMp3InTemp(path string) bool {
	f, err := os.CreateTemp("", "goplay")
	defer os.Remove(f.Name())
	if err != nil {
		log.Println(err)
		return false
	}
	sourceFile, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return false
	}
	_, err = io.Copy(f, sourceFile)
	if err != nil {
		log.Println(err)
	}
	err = p.Converter.ConvertToMP3(f.Name())
	if err != nil {
		log.Println(err)
		return false
	}
	p.AddMusicFromFile(f.Name())
	return true
}

func (p *Player) AddMusicsFromDir(dirPath string) {
	dir, err := os.Open(dirPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	entries, err := dir.Readdir(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			isMp3, err := p.Converter.IsMp3(dirPath + "/" + entry.Name())
			if err != nil {
				log.Println(err)
				continue
			}
			if isMp3 {
				log.Println("Playing music from dir", dirPath+"/"+entry.Name())
				p.AddMusicFromFile(filepath.Join(dirPath, entry.Name()))
			} else {
				p.addConvertedMp3InTemp(filepath.Join(dirPath, entry.Name()))
			}
		}
	}
}

// the unique is the unique id of the music in the engine it can be url or id
func (p *Player) AddMusicFromOnline(unique string, engineName string) {
	p.addTask(unique, shared.Downloading)
	path, err := p.Director.Download(engineName, unique)
	if err != nil {
		log.Println(err)
		p.errorifyTask(unique, err)
		return
	}
	err = p.Converter.ConvertToMP3(path)
	if err != nil {
		log.Println(err)
		p.errorifyTask(unique, err)
		return
	}

	if err != nil {
		log.Println(err)
		p.errorifyTask(unique, err)
		return
	}
	if path == "" {
		p.errorifyTask(unique, fmt.Errorf("failed to download music: %s", unique))
		return
	}
	p.AddMusicFromFile(path)
	p.removeTask(unique)
}

func (p *Player) AddMusicFromPlaylistByName(playlistName string, musicName string) {
	playlistPath := filepath.Join(config.GetConfig().GoPlayPath, config.GetConfig().PlaylistPath, playlistName, musicName)
	p.AddMusicFromFile(playlistPath)
}

func (p *Player) AddMusicFromPlaylistByIndex(playlistName string, index int) {
	playlistPath := filepath.Join(config.GetConfig().GoPlayPath, config.GetConfig().PlaylistPath, playlistName)
	dir, err := os.Open(playlistPath)
	if err != nil {
		log.Println(err)
		return
	}
	entries, err := dir.Readdir(0)
	if err != nil {
		log.Println(err)
		return
	}
	if index < len(entries) && index >= 0 {
		p.AddMusicFromFile(filepath.Join(playlistPath, entries[index].Name()))
	}
}

func (p *Player) AddMusicsFromPlaylist(playlistName string) {
	playlistPath := filepath.Join(config.GetConfig().GoPlayPath, config.GetConfig().PlaylistPath, playlistName)
	dir, err := os.Open(playlistPath)
	if err != nil {
		log.Println(err)
		return
	}
	entries, err := dir.Readdir(0)
	if err != nil {
		log.Println(err)
		return
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			p.AddMusicFromFile(filepath.Join(playlistPath, entry.Name()))
		}
	}
}
