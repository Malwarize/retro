package online

import (
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/Malwarize/goplay/shared"
)

type OnlineEngine interface {
	Search(query string, maxResults int) ([]shared.SearchResult, error)
	Download(url string) (io.ReadCloser, string, error)
	Exists(url string) (bool, error)
	Name() string
}

type OnlineDirector struct {
	engines map[string]OnlineEngine // key: engine name, value: engine
	Cached  *CachedFiles
}

func NewOnlineDirector() *OnlineDirector {
	return &OnlineDirector{
		engines: make(map[string]OnlineEngine),
		Cached:  NewCachedFiles(),
	}
}

func NewDefaultDirector() (*OnlineDirector, error) {
	director := NewOnlineDirector()
	director.Cached.Fetch()
	youtubeEngine, err := newYoutubeEngine()
	if err != nil {
		return director, fmt.Errorf("failed to create youtube engine: %w", err)
	}

	// register the engines here
	director.Register("youtube", youtubeEngine)
	return director, nil
}

func (od *OnlineDirector) Register(name string, engine OnlineEngine) {
	od.engines[name] = engine
}

var times = 0

func (od *OnlineDirector) Search(engineName, query string, maxResults int) ([]shared.SearchResult, error) {
	engine, ok := od.engines[engineName]
	if !ok {
		return nil, errors.New("engine not found")
	}
	return engine.Search(query, maxResults)
}

// func (od *OnlineDirector) Download(engineName, url string) (io.ReadCloser, string, error) {
// 	engine, ok := od.engines[engineName]
// 	if !ok {
// 		return nil, "", errors.New("engine not found")
// 	}
// 	// check if file is Cached
// 	name, err := od.Cached.GetFileByKey(url, engineName)
// 	if err == nil {
// 		f, err := os.Open(name)
// 		if err != nil {
// 			return nil, "", err
// 		}
// 		return f, name, nil
// 	}

// 	log.Println("Downloading file from ", url)
// 	reader, name, err := engine.Download(url)
// 	log.Println("Downloaded file from ", url)
// 	if err != nil {
// 		return nil, "", err
// 	}
// 	// cache it
// 	data, err := io.ReadAll(reader)
// 	if err != nil {
// 		return nil, "", err
// 	}
// 	path := od.Cached.AddFile(data, name, engineName, url)
// 	if path == "" {
// 		return nil, "", errors.New("failed to cache file")
// 	}
// 	f, err := os.Open(path)
// 	if err != nil {
// 		return nil, "", err
// 	}
// 	return f, path, nil
// }

func (od *OnlineDirector) Download(engineName, url string) (path string, err error) {
	engine, ok := od.engines[engineName]
	if !ok {
		return "", errors.New("engine not found")
	}
	// check if file is Cached
	name, err := od.Cached.GetFileByKey(url, engineName)
	if err == nil {
		return name, nil
	}
	log.Println("Downloading file from ", url)
	reader, name, err := engine.Download(url)
	log.Println("Downloaded file from ", url)
	if err != nil {
		return "", err
	}
	// cache it
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	path = od.Cached.AddFile(data, name, engineName, url)
	if path == "" {
		return "", errors.New("failed to cache file")
	}
	return path, nil
}

func (od *OnlineDirector) GetEngines() map[string]OnlineEngine {
	return od.engines
}
