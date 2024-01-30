package online

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Malwarize/goplay/player"
)

type OnlineEngine interface {
	Search(query string, maxResults int) ([]string, error)
	Download(url string) (io.ReadCloser, string, error)
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

func (od *OnlineDirector) Register(name string, engine OnlineEngine) {
	od.engines[name] = engine
}

func (od *OnlineDirector) Search(engineName, query string, maxResults int) ([]string, error) {
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
	name, err = od.cached.GetFileByName(url)
	if err == nil {
		f, err := os.Open(name)
		if err != nil {
			return nil, "", err
		}
		return f, name, nil
	}

	// if not, download it

	reader, name, err := engine.Download(url)
	if err != nil {
		return nil, "", err
	}
	// cache it
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, "", err
	}
	od.cached.AddFile(data, name, engineName, url)
	path := od.cached.BaseDir + "/" + engineName + "/" + name + "_" + url

	f, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	return f, path, nil
}

func Test() {
	director := NewOnlineDirector()
	director.cached.Fetch()
	youtubeEngine, err := newYoutubeEngine()
	if err != nil {
		panic(err)
	}
	director.Register("youtube", youtubeEngine)

	musics, err := director.Search("youtube", "golang", 10)
	if err != nil {
		panic(err)
	}
	for _, music := range musics {
		fmt.Println(music)
	}
	fmt.Println("=====================================")
	_, path, err := director.Download("youtube", musics[0])
	if err != nil {
		panic(err)
	}
	converter, err := player.NewConverter("ffmpeg", "ffprobe")
	if err != nil {
		panic(err)
	}
	is, err := converter.IsMp3(path)
	if err != nil {
		panic(err)
	}
	if is {
		println("is mp3")
	} else {
		println("not mp3")
		println("converting...")
		err = converter.ConvertToMP3(path)
		if err != nil {
			panic(err)
		}
		println("converted")

		// check again
		is, err = converter.IsMp3(path)
		if err != nil {
			panic(err)
		}
		if is {
			println("is mp3")
		} else {
			println("not mp3")
		}

	}

}
