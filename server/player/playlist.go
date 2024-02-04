package player

import (
	"os"
	"path/filepath"

	"github.com/Malwarize/goplay/shared"
)

type PlayList struct {
	Name  string
	Items []string
}

func NewPlayList(name string) *PlayList {
	file, err := os.Open(
		filepath.Join(shared.PlaylistPath, name),
	)
	if err != nil {
		return nil
	}
	file.Close()

	return &PlayList{
		Name:  name,
		Items: make([]string, 0),
	}
}

func (pl *PlayList) Fetch() error {
	dir, err := os.Open(filepath.Join(shared.PlaylistPath, pl.Name))
	if err != nil {
		return err
	}
	defer dir.Close()
	names, err := dir.Readdirnames(-1)
	if err != nil {
		return err
	}
	pl.Items = names
	return nil
}

func (pl *PlayList) Add(path string, name string) error {
	err := os.Rename(path, filepath.Join(shared.PlaylistPath, pl.Name, name))
	if err != nil {
		return err
	}
	pl.Items = append(pl.Items, name)
	return nil
}

func (pl *PlayList) Remove(name string) error {
	err := os.Remove(filepath.Join(shared.PlaylistPath, pl.Name, name))
	if err != nil {
		return err
	}
	for i, n := range pl.Items {
		if n == name {
			pl.Items = append(pl.Items[:i], pl.Items[i+1:]...)
			break
		}
	}
	return nil
}

type PlayListManager struct {
	PlayLists map[string]*PlayList
}

func NewPlayListManager() *PlayListManager {
	return &PlayListManager{
		PlayLists: make(map[string]*PlayList),
	}
}

func (plm *PlayListManager) Fetch() error {
	dir, err := os.Open(shared.PlaylistPath)
	if err != nil {
		return err
	}
	defer dir.Close()
	names, err := dir.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		pl := NewPlayList(name)
		err := pl.Fetch()
		if err != nil {
			return err
		}
		plm.PlayLists[name] = pl
	}
	return nil
}

func (plm *PlayListManager) Add(name string) error {
	err := os.MkdirAll(filepath.Join(shared.PlaylistPath, name), os.ModePerm)
	if err != nil {
		return err
	}
	plm.PlayLists[name] = NewPlayList(name)
	return nil
}

func (plm *PlayListManager) Remove(name string) error {
	err := os.RemoveAll(filepath.Join(shared.PlaylistPath, name))
	if err != nil {
		return err
	}
	delete(plm.PlayLists, name)
	return nil
}

func (plm *PlayListManager) Get(name string) *PlayList {
	return plm.PlayLists[name]
}
