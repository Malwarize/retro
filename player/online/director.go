package online

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Malwarize/goplay/shared"
)

type OnlineEngine interface {
	Search(query string, maxResults int) ([]shared.SearchResult, error)
	Download(url string) (io.ReadCloser, string, error)
	Exists(url string) (bool, error)
}

type OnlineDirector struct {
	engines map[string]OnlineEngine // key: engine name, value: engine
	cached  *CachedFiles
}

func NewOnlineDirector() *OnlineDirector {
	return &OnlineDirector{
		engines: make(map[string]OnlineEngine),
		cached:  NewCachedFiles("./cache"),
	}
}
func NewDefaultDirector() (*OnlineDirector, error) {
	director := NewOnlineDirector()
	director.cached.Fetch()
	youtubeEngine, err := newYoutubeEngine()
	if err != nil {
		return director, fmt.Errorf("failed to create youtube engine: %w", err)
	}
	director.Register("youtube", youtubeEngine)
	return director, nil
}
func (od *OnlineDirector) Register(name string, engine OnlineEngine) {
	od.engines[name] = engine
}

func (od *OnlineDirector) Search(engineName, query string, maxResults int) ([]shared.SearchResult, error) {
	engine, ok := od.engines[engineName]
	if !ok {
		return nil, errors.New("engine not found")
	}
	return engine.Search(query, maxResults)
}

func (od *OnlineDirector) Download(engineName, url string) (io.ReadCloser, string, error) {
	engine, ok := od.engines[engineName]
	if !ok {
		return nil, "", errors.New("engine not found")
	}
	// check if file is cached
	name, err := od.cached.GetFileByUrl(url, engineName)
	if err == nil {
		f, err := os.Open(name)
		if err != nil {
			return nil, "", err
		}
		return f, name, nil
	}

	log.Println("Downloading file from ", url)
	reader, name, err := engine.Download(url)
	fmt.Println("NAMEEEEEE ", name)
	log.Println("Downloaded file from ", url)
	if err != nil {
		return nil, "", err
	}
	// cache it
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, "", err
	}
	od.cached.AddFile(data, name, engineName, url)
	path := od.cached.BaseDir + "/" + engineName + "/" + shared.EscapeSpecialDirChars(name) + "_" + getYoutubeIdFromUrl(url)

	f, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	return f, path, nil
}

func (od *OnlineDirector) GetEngines() map[string]OnlineEngine {
	return od.engines
}
