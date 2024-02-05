package player

import (
	"log"
	"os"
	"path/filepath"

	"github.com/Malwarize/goplay/config"
	"github.com/Malwarize/goplay/server/player/online"
)

type PlayList struct {
	Name  string
	Items []Music
}

type PlayListManager struct {
	PlayLists map[string]PlayList // map playlistname to playlist
}

func NewPlayListManager() *PlayListManager {
	return &PlayListManager{PlayLists: make(map[string]PlayList)}
}

func (plm *PlayListManager) Fetch() error {
	// fetch all playlists
	files, err := os.ReadDir(config.GetConfig().PlaylistPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		pl := PlayList{Name: file.Name()}
		pl.Items = make([]Music, 0)
		songs, err := os.ReadDir(filepath.Join(config.GetConfig().PlaylistPath, file.Name()))
		if err != nil {
			return err
		}
		for _, song := range songs {
			if song.IsDir() {
				continue
			}
			pl.Items = append(pl.Items, Music{Path: filepath.Join(config.GetConfig().PlaylistPath, file.Name(), song.Name())})
		}

		plm.PlayLists[file.Name()] = pl
	}

	return nil
}

// Create a new playlist
func (plm *PlayListManager) Create(name string) error {
	err := os.Mkdir(filepath.Join(config.GetConfig().PlaylistPath, name), 0755)
	if err != nil {
		return err
	}
	pl := PlayList{Name: name}
	plm.PlayLists[name] = pl
	return nil
}

// remove a playlist
func (plm *PlayListManager) Remove(name string) error {
	err := os.RemoveAll(filepath.Join(config.GetConfig().PlaylistPath, name))
	if err != nil {
		return err
	}
	delete(plm.PlayLists, name)
	return nil
}

// add music to a playlist
func (plm *PlayListManager) AddMusic(name string, music Music) error {
	err := copyFile(music.Path, filepath.Join(config.GetConfig().PlaylistPath, name, music.Name()))
	if err != nil {
		return err
	}
	pl, ok := plm.PlayLists[name]
	if !ok {
		return os.ErrNotExist
	}
	pl.Items = append(pl.Items, music)

	plm.PlayLists[name] = pl

	return nil
}

func (plm *PlayListManager) RemoveMusic(name string, index int) error {
	music := plm.PlayLists[name].Items[index]
	err := os.Remove(filepath.Join(config.GetConfig().PlaylistPath, name, music.Name()))
	if err != nil {
		return err
	}
	pl, ok := plm.PlayLists[name]
	if !ok {
		return os.ErrNotExist
	}
	pl.Items = append(pl.Items[:index], pl.Items[index+1:]...)

	plm.PlayLists[name] = pl
	return nil
}

func (plm *PlayListManager) PlayListsNames() []string {
	var names []string
	for name := range plm.PlayLists {
		names = append(names, name)
	}
	return names
}

func (plm *PlayListManager) PlayListSongs(name string) []string {
	pl, ok := plm.PlayLists[name]
	if !ok {
		return nil
	}
	var songs []string
	for _, song := range pl.Items {
		songs = append(songs, song.Name())
	}
	return songs
}

func (plm *PlayListManager) AddToPlayListFromDir(name string, dir string, converter *Converter) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		// check if its mp3
		log.Println("Checking if", file.Name(), "is mp3")
		path := filepath.Join(dir, file.Name())
		isMp3, err := converter.IsMp3(path)
		if err != nil {
			return err
		}
		if !isMp3 {
			continue
		}
		err = plm.AddMusic(name, Music{Path: path})
		if err != nil {
			return err
		}
	}
	return nil
}

func (plm *PlayListManager) AddToPlayListFromFile(name string, file string) error {
	err := plm.AddMusic(name, Music{Path: file})
	if err != nil {
		return err
	}
	return nil
}

func (plm *PlayListManager) AddToPlayListFromOnline(name string, query string, engineName string, director *online.OnlineDirector, converter *Converter) error {
	path, err := director.Download(engineName, query)
	if err != nil {
		return err
	}
	err = plm.AddMusic(name, Music{Path: path})
	if err != nil {
		return err
	}
	err = converter.ConvertToMP3(path)
	if err != nil {
		return err
	}
	return nil
}

func (plm *PlayListManager) Exists(name string) bool {
	_, ok := plm.PlayLists[name]
	return ok
}
