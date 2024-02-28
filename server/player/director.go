package player

import (
	"errors"
	"fmt"
	"io"

	"github.com/Malwarize/goplay/config"
	"github.com/Malwarize/goplay/logger"
	"github.com/Malwarize/goplay/server/player/db"
	en "github.com/Malwarize/goplay/server/player/engines"
	"github.com/Malwarize/goplay/shared"
)

type Director struct {
	Converter *Converter
	Db        *db.Db
	engines   map[string]en.Engine // key: engine name, value: engine
}

func NewDirector(db *db.Db) (*Director, error) {
	c, err := NewConverter()
	if err != nil {
		return nil, err
	}
	return &Director{
		engines:   make(map[string]en.Engine),
		Db:        db,
		Converter: c,
	}, nil
}

func NewDefaultDirector() (*Director, error) {
	db, err := db.LoadDb(config.GetConfig().DbPath)
	if err != nil {
		return nil, err
	}
	director, err := NewDirector(
		db,
	)
	if err != nil {
		return nil, err
	}

	youtubeEngine, err := en.NewYoutubeEngine()
	if err != nil {
		return director, fmt.Errorf("failed to create youtube engine: %w", err)
	}

	// register the engines here
	director.Register(youtubeEngine)
	return director, nil
}

func (od *Director) Register(engine en.Engine) {
	od.engines[engine.Name()] = engine
}

var times = 0

func (od *Director) Search(
	engineName, query string,
) ([]shared.SearchResult, error) {
	engine, ok := od.engines[engineName]
	if !ok {
		return nil, errors.New("engine not found")
	}

	return engine.Search(query, engine.MaxResults())
}

func (od *Director) Download(engineName, url string) (*db.Music, error) {
	engine, ok := od.engines[engineName]
	if !ok {
		return nil, errors.New("engine not found")
	}
	// check if file is Cached
	music, err := od.Db.GetMusic(
		engine.Name(),
		url,
	)
	if err == nil {
		return &music, nil
	}

	logger.LogInfo("Downloading file from ", url)
	reader, name, err := engine.Download(url)
	logger.LogInfo("Downloaded file from ", url)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	isMp3, err := od.Converter.IsMp3(data)
	var mp3data []byte
	if !isMp3 {
		mp3data, err = od.Converter.ConvertToMP3(data)
		if err != nil {
			return nil, err
		}
	} else {
		mp3data = data
	}
	// cache it to db
	music = db.Music{
		Name:   name,
		Source: engine.Name(),
		Key:    url,
		Data:   mp3data,
	}
	err = od.Db.AddMusic(&music)
	if err != nil {
		logger.LogWarn(
			"failed to add music to db",
			err,
		)
		return &music, err
	}
	return &music, nil
}

func (od *Director) GetEngines() map[string]en.Engine {
	return od.engines
}
