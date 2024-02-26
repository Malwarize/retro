package player

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/Malwarize/goplay/config"
	"github.com/Malwarize/goplay/logger"
	"github.com/Malwarize/goplay/shared"
)

type PlayList struct {
	Name  string
	Items []*Music
}

type PlayListManager struct {
	PlayListPath string
	PlayLists    map[string]*PlayList // map playlistname to playlist
	mu           *sync.Mutex
}

func NewPlayListManager() (*PlayListManager, error) {
	path := config.GetConfig().PlaylistPath
	err := os.MkdirAll(path, 0o755)
	if err != nil {
		return nil, err
	}
	return &PlayListManager{
		PlayLists:    make(map[string]*PlayList),
		PlayListPath: path,
		mu:           &sync.Mutex{},
	}, nil
}

func (plm *PlayListManager) Fetch() error {
	// fetch all playlists
	files, err := os.ReadDir(plm.PlayListPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		pl := &PlayList{Name: file.Name()}
		pl.Items = make([]*Music, 0)
		songs, err := os.ReadDir(filepath.Join(plm.PlayListPath, file.Name()))
		if err != nil {
			return err
		}
		for _, song := range songs {
			if song.IsDir() {
				continue
			}
			pl.Items = append(
				pl.Items,
				&Music{
					Path: filepath.Join(plm.PlayListPath,
						file.Name(),
						song.Name(),
					),
				},
			)
		}

		plm.PlayLists[file.Name()] = pl
	}
	return nil
}

// Create a new playlist
func (plm *PlayListManager) CreatePlayList(name string) error {
	err := os.Mkdir(filepath.Join(plm.PlayListPath, name), 0o755)
	if err != nil {
		return err
	}
	pl := &PlayList{Name: name}
	// check if playlist already exists
	_, ok := plm.PlayLists[name]
	if ok {
		return os.ErrExist
	}
	plm.mu.Lock()
	defer plm.mu.Unlock()
	plm.PlayLists[name] = pl
	return nil
}

func (plm *PlayListManager) Remove(pl *PlayList) error {
	err := os.RemoveAll(filepath.Join(plm.PlayListPath, pl.Name))
	if err != nil {
		return err
	}
	plm.mu.Lock()
	defer plm.mu.Unlock()
	delete(plm.PlayLists, pl.Name)
	return nil
}

func (plm *PlayListManager) AddMusic(pl *PlayList, music *Music) error {
	err := copyFile(music.Path, filepath.Join(plm.PlayListPath, pl.Name, music.Name()))
	if err != nil {
		return err
	}
	plm.mu.Lock()
	defer plm.mu.Unlock()
	pl.Items = append(pl.Items, music)
	plm.PlayLists[pl.Name] = pl
	return nil
}

func (plm *PlayListManager) RemoveMusic(pl *PlayList, music *Music) error {
	plm.mu.Lock()
	defer plm.mu.Unlock()

	err := os.Remove(music.Path)
	if err != nil {
		return err
	}
	for i, m := range pl.Items {
		if m.Path == music.Path {
			pl.Items = append(pl.Items[:i], pl.Items[i+1:]...)
		}
	}
	plm.PlayLists[pl.Name] = pl
	return nil
}

func (plm *PlayListManager) GetPlayListByName(name string) (*PlayList, error) {
	plm.mu.Lock()
	defer plm.mu.Unlock()
	if _, ok := plm.PlayLists[name]; !ok {
		return nil, os.ErrNotExist
	}
	return plm.PlayLists[name], nil
}

func (plm *PlayListManager) GetPlayListSongByName(pl *PlayList, name string) (*Music, error) {
	plm.mu.Lock()
	defer plm.mu.Unlock()
	for _, song := range pl.Items {
		if song.Name() == name {
			return song, nil
		}
	}
	return nil, os.ErrNotExist
}

func (plm *PlayListManager) plSize(pl *PlayList) int {
	plm.mu.Lock()
	defer plm.mu.Unlock()
	return len(pl.Items)
}

func (plm *PlayListManager) GetPlayListSongByIndex(pl *PlayList, index int) (*Music, error) {
	plm.mu.Lock()
	defer plm.mu.Unlock()
	if index < 0 || index >= plm.plSize(pl) {
		return nil, os.ErrNotExist
	}
	return pl.Items[index], nil
}

func (plm *PlayListManager) PlayListsNames() []string {
	plm.mu.Lock()
	defer plm.mu.Unlock()
	var names []string
	for name := range plm.PlayLists {
		names = append(names, name)
	}
	return names
}

func (plm *PlayListManager) GetPlayListSongs(pl *PlayList) []*Music {
	plm.mu.Lock()
	defer plm.mu.Unlock()
	return pl.Items
}

func (plm *PlayListManager) AddToPlayListFromDir(
	pl *PlayList,
	dir string,
	converter *Converter,
) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		// check if its mp3
		logger.LogInfo("Checking if", file.Name(), "is mp3")
		path := filepath.Join(dir, file.Name())
		isMp3, err := converter.IsMp3(path)
		if err != nil {
			return err
		}
		if !isMp3 {
			logger.LogWarn(
				"Skipping file", file.Name(),
				"in directory", dir,
				"because it is not an mp3 file",
			)
			continue
		}
		err = plm.AddMusic(
			pl,
			&Music{Path: path},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (plm *PlayListManager) AddToPlayListFromFile(pl *PlayList, file string) error {
	err := plm.AddMusic(
		pl,
		&Music{Path: file},
	)
	if err != nil {
		return err
	}
	return nil
}

func (plm *PlayListManager) AddToPlayListFromOnline(
	pl *PlayList,
	query string,
	engineName string,
	p *Player,
) {
	p.addTask(query, shared.Downloading)
	path, err := p.Director.Download(engineName, query)
	if err != nil {
		logger.LogWarn(
			"Error downloading music from", engineName,
			"with query", query,
		)
		p.errorifyTask(query, err)
		return
	}
	err = p.Converter.ConvertToMP3(path)
	if err != nil {
		logger.ERRORLogger.Println(
			"Error converting music from", engineName,
			"with query", query,
			"path", path,
		)
		p.errorifyTask(query, err)
		return
	}
	err = plm.AddMusic(pl, &Music{Path: path})
	if err != nil {
		logger.ERRORLogger.Println(
			"Error adding music to playlist", pl.Name,
			"with query", query,
			"path", path,
		)
		p.errorifyTask(query, err)
		return
	}
	p.removeTask(query)
}

func (plm *PlayListManager) Exists(plname string) bool {
	_, ok := plm.PlayLists[plname]
	return ok
}
