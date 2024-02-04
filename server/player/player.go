package player

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/Malwarize/goplay/server/player/online"
	"github.com/Malwarize/goplay/shared"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
)

var PlayerInstance *Player

const (
	Playing = iota
	Paused
	Stopped
)

type Player struct {
	Queue           *MusicQueue
	playerState     int
	done            chan struct{}
	initialised     bool
	Converter       *Converter
	Director        *online.OnlineDirector
	Tasks           map[string]shared.Task
	PlayListManager *PlayListManager
	mu              sync.Mutex
}

func NewPlayer() *Player {
	converter, err := NewConverter()
	if err != nil {
		log.Fatal(err)
	}
	director, err := online.NewDefaultDirector()
	if err != nil {
		log.Fatal(err)
	}
	playlistManager := NewPlayListManager()
	err = playlistManager.Fetch()
	if err != nil {
		log.Fatal(err)
	}

	return &Player{
		Queue:           NewMusicQueue(),
		playerState:     Stopped,
		done:            make(chan struct{}),
		initialised:     false,
		Converter:       converter,
		Director:        director,
		PlayListManager: playlistManager,
		Tasks:           make(map[string]shared.Task),
	}
}

//singleton player
func GetPlayer() *Player {
	if PlayerInstance == nil {
		PlayerInstance = NewPlayer()
	}
	return PlayerInstance
}

func (p *Player) Play() {
	music := p.GetCurrentMusic()
	if !p.initialised {
		speaker.Init(music.Format.SampleRate, music.Format.SampleRate.N(time.Second/10))
		p.initialised = true
	} else {
		speaker.Clear()
		speaker.Init(music.Format.SampleRate, music.Format.SampleRate.N(time.Second/10))
	}
	p.playerState = Playing
	go func() {
		done := make(chan struct{})
		speaker.Play(beep.Seq(music.Streamer, beep.Callback(func() {
			done <- struct{}{}
		})))
		<-done
		p.Next()
	}()
}

func (p *Player) Next() {
	if p.Queue.IsEmpty() || p.playerState == Stopped {
		return
	}
	if p.playerState == Paused {
		p.Resume()
	}
	p.GetCurrentMusic().Streamer.Seek(0)
	p.QueueNext()
	log.Println("Next music index", p.Queue.GetCurrentIndex())
	p.Play()
}

func (p *Player) Prev() {
	if p.Queue.IsEmpty() || p.playerState == Stopped {
		return
	}
	if p.playerState == Paused {
		p.Resume()
	}
	p.GetCurrentMusic().Streamer.Seek(0)
	p.QueuePrev()
	p.Play()
}

func (p *Player) Stop() {
	clear(p.Tasks)
	if p.playerState == Stopped {
		return
	}
	if p.playerState == Paused {
		p.Resume()
	}
	speaker.Clear()
	p.Queue.Clear()
	p.playerState = Stopped
}

func (p *Player) Pause() {
	if p.playerState == Paused || p.playerState == Stopped {
		return
	}
	p.playerState = Paused
	speaker.Lock()
}

func (p *Player) Resume() {
	if p.playerState == Playing || p.playerState == Stopped {
		return
	}
	p.playerState = Playing
	speaker.Unlock()
}

func (p *Player) Seek(d time.Duration) {
	if p.playerState == Stopped {
		return
	}
	speaker.Lock()
	defer speaker.Unlock()
	currentMusic := p.GetCurrentMusic()
	currentSamplePos := currentMusic.Streamer.Position()
	curretnTimePos := currentMusic.Format.SampleRate.D(currentSamplePos)
	newTimePos := (curretnTimePos + d) % p.GetCurrentMusicLength()
	newSamplePos := currentMusic.Format.SampleRate.N(newTimePos)
	if err := currentMusic.Streamer.Seek(newSamplePos); err != nil {
		fmt.Println(err)
	}
}

func (p *Player) GetCurrentMusicPosition() time.Duration {
	if p.playerState == Stopped {
		return 0
	}
	currentMusic := p.GetCurrentMusic()
	currentSamplePos := currentMusic.Streamer.Position()
	curretnTimePos := currentMusic.Format.SampleRate.D(currentSamplePos)
	return curretnTimePos
}

func (p *Player) GetCurrentMusicLength() time.Duration {
	if p.Queue.IsEmpty() {
		return 0
	}
	music := p.GetCurrentMusic()
	return music.Format.SampleRate.D(music.Streamer.Len())
}

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

func (p *Player) GetAvailableMusicOptions(query string) []shared.SearchResult {
	var musics []shared.SearchResult
	for engineName := range p.Director.GetEngines() {
		searchDone := make(chan bool)
		go func() {
			searchRes, err := p.Director.Search(engineName, query, 5)
			searchDone <- true
			if err != nil {
				p.errorifyTask(query, err)
				log.Println("Failed to search for", query, ":", err)
			}
			musics = append(musics, searchRes...)
		}()

		select {
		case <-time.After(60 * time.Second):
			p.errorifyTask(query, errors.New("Timeout searching for "+query))
			log.Println("Timeout searching for", query)
		case <-searchDone:
			close(searchDone)
		}
	}
	p.Director.Cached.Fetch()
	files := p.Director.Cached.Search(query)
	for _, f := range files {
		name := filepath.Base(f)
		musics = append(
			musics,
			shared.SearchResult{
				Title:       name,
				Type:        "cache",
				Destination: f,
			},
		)
	}
	// etc ...
	return musics
}

//if result is empty, it means it detects and plays the music if succeed other wise it returns the search results
func (p *Player) DetectAndPlay(unknown string) []shared.SearchResult {
	fmt.Println("Checking what is this")
	whatIsThis := p.CheckWhatIsThis(unknown)
	switch whatIsThis {
	case "dir":
		log.Println("Detected dir")
		go func() {
			p.AddMusicsFromDir(unknown)
			p.Play()
		}()
	case "file":
		log.Println("Detected file")
		go func() {
			p.AddMusicFromFile(unknown)
			p.Play()
		}()
	case "queue":
		log.Println("Detected queue")
		index, _ := strconv.Atoi(unknown)
		go func() {
			p.Queue.SetCurrentIndex(index)
			p.Play()
		}()
	case "unknown":
		log.Println("Detected unknown, searching for", unknown)
		return p.GetAvailableMusicOptions(unknown)
	default:
		log.Println("Detected Engine", whatIsThis)
		go func() {
			p.AddMusicFromOnline(unknown, whatIsThis)
			p.Play()
		}()
	}
	return []shared.SearchResult{}
}

func (p *Player) GetPlayerStatus() shared.Status {
	return shared.Status{
		CurrentMusicIndex:    p.Queue.GetCurrentIndex(),
		CurrentMusicPosition: p.GetCurrentMusicPosition(),
		CurrentMusicLength:   p.GetCurrentMusicLength(),
		PlayerState:          p.playerState,
		MusicQueue:           p.Queue.GetTitles(),
		Tasks:                p.Tasks,
	}
}
