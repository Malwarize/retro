package player

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/Malwarize/goplay/config"
	"github.com/Malwarize/goplay/server/player/online"
	"github.com/Malwarize/goplay/shared"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
)

var PlayerInstance *Player
var once sync.Once

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
	Vol             int
	mu              sync.Mutex
}

func NewPlayer() *Player {
	if _, err := os.Stat(config.GetConfig().GoPlayPath); os.IsNotExist(err) {
		log.Println("goplay dir not found, creating it")
		err = os.Mkdir(config.GetConfig().GoPlayPath, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	converter, err := NewConverter()
	if err != nil {
		log.Fatal(err)
	}
	director, err := online.NewDefaultDirector()
	if err != nil {
		log.Fatal(err)
	}
	playlistManager, err := NewPlayListManager()
	if err != nil {
		log.Fatal(err)
	}
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
		Vol:             100,
		Tasks:           make(map[string]shared.Task),
	}
}

//singleton player
func GetPlayer() *Player {
	once.Do(func() {
		PlayerInstance = NewPlayer()
	})
	return PlayerInstance
}

func (p *Player) Play() {
	if p.Queue.IsEmpty() {
		return
	}
	music := p.Queue.GetCurrentMusic()
	if !p.initialised {
		speaker.Init(music.Format.SampleRate, music.Format.SampleRate.N(time.Second/10))
		p.initialised = true
	} else {
		speaker.Clear()
		speaker.Init(music.Format.SampleRate, music.Format.SampleRate.N(time.Second/10))
	}
	p.setPlayerState(Playing)
	go func() {
		done := make(chan struct{})
		music.SetVolume(p.Vol)
		speaker.Play(beep.Seq(music.Volume, beep.Callback(func() {
			done <- struct{}{}
		})))
		<-done
		p.Next()
	}()
}

func (p *Player) getPlayerState() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.playerState
}

func (p *Player) setPlayerState(state int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.playerState = state
}

func (p *Player) Next() {
	if p.Queue.IsEmpty() || p.getPlayerState() == Stopped {
		return
	}
	if p.getPlayerState() == Paused {
		p.Resume()
	}
	p.Queue.GetCurrentMusic().Streamer().Seek(0)
	p.Queue.QueueNext()
	p.Play()
}

func (p *Player) Prev() {
	state := p.getPlayerState()
	if p.Queue.IsEmpty() || state == Stopped {
		return
	}
	if state == Paused {
		p.Resume()
	}
	currentMusic := p.Queue.GetCurrentMusic()
	currentMusic.Streamer().Seek(0)
	p.Queue.QueuePrev()
	p.Play()
}

func (p *Player) Stop() {
	clear(p.Tasks)

	state := p.getPlayerState()
	if state == Stopped {
		return
	}
	if state == Paused {
		p.Resume()
	}
	speaker.Clear()
	p.Queue.Clear()
	p.setPlayerState(Stopped)
}

func (p *Player) Pause() {
	state := p.getPlayerState()
	if state == Paused || state == Stopped {
		return
	}
	p.setPlayerState(Paused)
	speaker.Lock()
}

func (p *Player) Resume() {
	state := p.getPlayerState()
	if state == Playing || state == Stopped {
		return
	}
	p.setPlayerState(Playing)
	speaker.Unlock()
}
func (p *Player) Seek(d time.Duration) {
	if p.getPlayerState() == Stopped {
		return
	}
	speaker.Lock()
	defer speaker.Unlock()
	currentMusic := p.Queue.GetCurrentMusic()
	currentSamplePos := currentMusic.Streamer().Position()
	curretnTimePos := currentMusic.Format.SampleRate.D(currentSamplePos)
	newTimePos := (curretnTimePos + d) % p.GetCurrentMusicLength()
	newSamplePos := currentMusic.Format.SampleRate.N(newTimePos)
	// check if seek is out of bounds
	if newTimePos < 0 {
		newSamplePos = 0
	}
	if newTimePos > p.GetCurrentMusicLength() {
		newSamplePos = currentMusic.Streamer().Len()
	}
	if err := currentMusic.Streamer().Seek(newSamplePos); err != nil {
		fmt.Println(err)
	}
}
func (p *Player) Volume(vp int /*volume percentage*/) {
	if p.getPlayerState() == Stopped {
		return
	}
	p.Vol = vp
	currentMusic := p.Queue.GetCurrentMusic()
	speaker.Lock()
	currentMusic.SetVolume(vp)
	speaker.Unlock()
}

func (p *Player) Remove(index int) {
	if p.Queue.IsEmpty() {
		return
	}
	if p.Queue.Size() == 1 {
		p.Stop()
	} else if index == p.Queue.GetCurrentIndex() {
		p.Next()
	}
	p.Queue.Remove(index)
}

// player playlist command
func (p *Player) CreatePlayList(name string) {
	p.PlayListManager.Create(name)
}

//delete playlist
func (p *Player) RemovePlayList(name string) {
	p.PlayListManager.Remove(name)
}

// list play list names
func (p *Player) PlayListsNames() []string {
	return p.PlayListManager.PlayListsNames()
}

func (p *Player) RemoveSongFromPlayList(name string, index int) {
	// check if the exists in the queue and remove it
	for i, music := range p.Queue.queue {
		if music.Path == filepath.Join(p.PlayListManager.PlayListPath, name, p.PlayListManager.PlayListSongs(name)[index]) {
			p.Queue.Remove(i)
		}
	}
	p.PlayListManager.RemoveMusic(name, index)
}

// list play list songs
func (p *Player) PlayListSongs(name string) []string {
	return p.PlayListManager.PlayListSongs(name)
}

func (p *Player) PlayListPlaySong(name string, index int) {
	p.AddMusicFromPlaylistByIndex(name, index)
	p.Play()
	if p.getPlayerState() == Stopped {
		p.Play()
	}
}

func (p *Player) PlayListPlayAll(name string) {
	p.AddMusicsFromPlaylist(name)
	if p.getPlayerState() == Stopped {
		p.Play()
	}
}
func (p *Player) GetCurrentMusicPosition() time.Duration {
	if p.getPlayerState() == Stopped {
		return 0
	}
	currentMusic := p.Queue.GetCurrentMusic()
	currentSamplePos := currentMusic.Streamer().Position()
	curretnTimePos := currentMusic.Format.SampleRate.D(currentSamplePos)
	return curretnTimePos
}

func (p *Player) GetCurrentMusicLength() time.Duration {
	if p.Queue.IsEmpty() {
		return 0
	}
	music := p.Queue.GetCurrentMusic()
	return music.Format.SampleRate.D(music.Streamer().Len())
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

func (p *Player) GetAvailableMusicOptions(query string) []shared.SearchResult {
	// add task : this task displayed in the status: if the task is done, it will be removed
	p.addTask(query, shared.Searching)
	var musics []shared.SearchResult
	for engineName := range p.Director.GetEngines() {
		searchDone := make(chan bool)
		go func() {
			defer close(searchDone)
			searchRes, err := p.Director.Search(engineName, query, 5)
			searchDone <- true
			if err != nil {
				p.errorifyTask(query, err)
				log.Println("Failed to search for", query, ":", err)
			}
			musics = append(musics, searchRes...)
		}()

		select {
		case <-time.After(config.GetConfig().SearchTimeOut):
			log.Println("Timeout searching for", query)
			p.errorifyTask(query, fmt.Errorf("timeout searching for %s", query))
		case <-searchDone:
		}
	}
	files := p.Director.Cached.Search(query)
	fmt.Println("Cached files", files)
	musics = append(musics, files...)
	p.removeTask(query)
	return musics
}

func (p *Player) DetectAndAddToPlayList(name string, query string) []shared.SearchResult {
	whatIsThis := p.CheckWhatIsThis(query)
	switch whatIsThis {
	case "dir":
		log.Println("Detected dir")
		err := p.PlayListManager.AddToPlayListFromDir(name, query, p.Converter)
		if err != nil {
			log.Println("Failed to add dir to playlist", err)
		}
	case "file":
		log.Println("Detected file")
		go p.PlayListManager.AddToPlayListFromFile(name, query)
	case "queue":
		index, _ := strconv.Atoi(query)
		music := p.Queue.GetMusicByIndex(index)
		go p.PlayListManager.AddToPlayListFromFile(name, music.Path)
	case "unknown":
		log.Println("Detected unknown, searching for", query)
		return p.GetAvailableMusicOptions(query)
	default:
		log.Println("Detected Engine", whatIsThis)
		go p.AddToPlayListFromOnline(name, query, whatIsThis, p.Director, p.Converter)
	}
	return []shared.SearchResult{}
}

//if result is empty, it means it detects and plays the music if succeed other wise it returns the search results
func (p *Player) DetectAndPlay(unknown string) []shared.SearchResult {
	log.Println("Checking what is this", unknown)
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
	case "playlist":
		log.Println("Detected playlist")
		go func() {
			p.AddMusicsFromPlaylist(unknown)
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
		PlayerState:          p.getPlayerState(),
		MusicQueue:           p.Queue.GetTitles(),
		Volume:               p.Vol,
		Tasks:                p.Tasks,
	}
}
