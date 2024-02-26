package player

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/Malwarize/goplay/config"
	"github.com/Malwarize/goplay/logger"
	"github.com/Malwarize/goplay/shared"
)

func (p *Player) CheckWhatIsThis(unknown string) string {
	// check if its a dir or file
	if fi, err := os.Stat(unknown); err == nil {
		if err == nil {
			if fi.IsDir() {
				// check if there is music files in the dir
				files, err := os.Open(unknown)
				if err != nil {
					return "unknown"
				}

				entries, err := files.Readdir(0)
				if err != nil {
					return "unknown"
				}

				for _, entry := range entries {
					if !entry.IsDir() {
						isMp3, _ := p.Converter.IsMp3(filepath.Join(unknown, entry.Name()))
						if isMp3 {
							return "dir"
						}
					}
				}
				return "unknown"
			} else {
				isMp3, _ := p.Converter.IsMp3(unknown)
				if isMp3 {
					return "file"
				}
			}
		}
	}
	// check if its play list name
	if p.PlayListManager.Exists(unknown) {
		return "playlist"
	}

	// check if its queue index
	if index, err := strconv.Atoi(unknown); err == nil && index >= 0 && index < p.Queue.Size() {
		return "queue"
	}
	engines := p.Director.GetEngines()
	for _, engine := range engines {
		ok, _ := engine.Exists(unknown)
		if ok {
			return engine.Name()
		}
	}
	return "unknown"
}

func (p *Player) searchWorker(
	ctx context.Context,
	engine string,
	query string,
	musicChan chan shared.SearchResult,
	wg *sync.WaitGroup,
) {
	defer func() {
		logger.LogInfo("Search worker done for", engine, query)
		wg.Done()
	}()

	searchRes, err := p.Director.Search(engine, query, 5)
	if err != nil {
		logger.LogWarn("Failed to search for", query, ":", err)
	}

	for _, music := range searchRes {
		musicChan <- music
	}
}

func (p *Player) GetAvailableMusicOptions(query string) []shared.SearchResult {
	// add task : this task displayed in the status: if the task is done, it will be removed
	p.addTask(query, shared.Searching)
	wg := &sync.WaitGroup{}
	musicChan := make(
		chan shared.SearchResult,
	)
	var musics []shared.SearchResult
	ctx, cancel := context.WithTimeout(context.Background(), config.GetConfig().SearchTimeOut)
	defer cancel()
	for name := range p.Director.GetEngines() {
		wg.Add(1)
		go p.searchWorker(ctx, name, query, musicChan, wg)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		cacheMusic := p.Director.Cached.Search(query)
		fmt.Println("Cached files", cacheMusic)
		for _, music := range cacheMusic {
			musicChan <- music
		}
	}()

	go func() {
		wg.Wait()
		close(musicChan)
	}()

	for {
		select {
		case music, ok := <-musicChan:
			if !ok {
				return musics
			}
			musics = append(musics, music)
			p.removeTask(
				query,
			)
		case <-ctx.Done():
			p.errorifyTask(
				query, fmt.Errorf("Search timed out"),
			)
			logger.LogWarn(
				"Search timed out",
				query,
			)
			return musics
		}
	}
}

func (p *Player) DetectAndAddToPlayList(
	plname string,
	query string,
) ([]shared.SearchResult, error) {
	whatIsThis := p.CheckWhatIsThis(query)
	pl, err := p.PlayListManager.GetPlayListByName(plname)
	if err != nil {
		return nil, logger.LogError(
			logger.GError(
				"Playlist not found",
			),
		)
	}
	switch whatIsThis {
	case "dir":
		logger.LogInfo(
			"detected dir for",
			query,
		)
		err = p.PlayListManager.AddToPlayListFromDir(
			pl,
			query,
			p.Converter,
		)
		if err != nil {
			return nil, logger.LogError(
				logger.GError(
					"Failed to add to playlist",
					err,
				),
			)
		}
	case "file":
		logger.LogInfo(
			"Detected file",
			query,
		)
		err := p.PlayListManager.AddToPlayListFromFile(
			pl,
			query,
		)
		if err != nil {
			return nil, logger.LogError(
				logger.GError(
					"Failed to add to playlist",
					err,
				),
			)
		}
	case "queue":
		logger.LogInfo(
			"Detected queue",
			query,
		)
		index, _ := strconv.Atoi(query)
		music := p.Queue.GetMusicByIndex(index)
		err := p.PlayListManager.AddToPlayListFromFile(
			pl,
			music.Path,
		)
		if err != nil {
			return nil, logger.LogError(
				logger.GError(
					"Failed to add to playlist",
					err,
				),
			)
		}
	case "unknown":
		logger.LogInfo(
			"Detected unknown",
			query,
		)
		return p.GetAvailableMusicOptions(query), nil
	default:
		logger.LogInfo("Detected Engine", whatIsThis)
		go p.PlayListManager.AddToPlayListFromOnline(pl, query, whatIsThis, p)
	}
	return []shared.SearchResult{}, nil
}

// if result is empty, it means it detects and plays the music if succeed other wise it returns the search results
func (p *Player) DetectAndPlay(unknown string) []shared.SearchResult {
	logger.LogInfo("Checking what is this", unknown)
	whatIsThis := p.CheckWhatIsThis(unknown)
	switch whatIsThis {
	case "dir":
		logger.LogInfo("Detected dir")
		p.AddMusicsFromDir(unknown)
		p.Play()
	case "file":
		logger.LogInfo("Detected file")
		go func() {
			p.AddMusicFromFile(unknown)
			p.Play()
		}()
	case "queue":
		logger.LogInfo("Detected queue")
		index, _ := strconv.Atoi(unknown)
		go func() {
			p.Queue.SetCurrIndex(index)
			p.Play()
		}()
	case "playlist":
		logger.LogInfo("Detected playlist")
		go func() {
			p.AddMusicsFromPlaylist(unknown)
			p.Play()
		}()
	case "unknown":
		logger.LogInfo("Detected unknown, searching for", unknown)
		return p.GetAvailableMusicOptions(unknown)
	default:
		logger.LogInfo("Detected Engine", whatIsThis)
		go func() {
			p.AddMusicFromOnline(unknown, whatIsThis)
			p.Play()
		}()
	}
	return []shared.SearchResult{}
}
