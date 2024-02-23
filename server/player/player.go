package player

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"

	"github.com/Malwarize/goplay/config"
	"github.com/Malwarize/goplay/server/player/online"
	"github.com/Malwarize/goplay/shared"
)

var PlayerInstance *Player

var once sync.Once

const (
	Playing = iota
	Paused
	Stopped
)

type lmeta struct {
	_lcurrentPos time.Duration
	_lcurrentDur time.Duration
}

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
	_lmeta          lmeta
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

func (p *Player) _setlMeta(m lmeta) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p._lmeta = m
}

func (p *Player) _getlMeta() lmeta {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p._lmeta
}

// singleton player
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
	music := p.Queue.GetCurrMusic()
	if music == nil {
		return
	}
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
	currentMusic := p.Queue.GetCurrMusic()
	if currentMusic == nil {
		return
	}
	if err := currentMusic.SetPositionD(0); err != nil {
		log.Println("Error when seeking")
	}
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
	currentMusic := p.Queue.GetCurrMusic()
	if currentMusic == nil {
		return
	}
	if err := currentMusic.Seek(0); err != nil {
		log.Println("Error when seeking")
	}
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
	p._setlMeta(lmeta{
		_lcurrentDur: p.GetCurrMusicDuration(),
		_lcurrentPos: p.GetCurrMusicPosition(),
	})
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
	state := p.getPlayerState()
	if state == Stopped {
		return
	}
	currentMusic := p.Queue.GetCurrMusic()
	if currentMusic == nil {
		return
	}
	if state == Paused {
		p.Resume()
		defer p.Pause()
	}
	if err := currentMusic.Seek(d); err != nil {
		log.Println("Error in seek", err)
	}
}

func (p *Player) Volume(vp int /*volume percentage*/) {
	if p.getPlayerState() == Stopped {
		return
	}
	p.Vol = vp
	currentMusic := p.Queue.GetCurrMusic()
	if currentMusic == nil {
		return
	}
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
	} else if index == p.Queue.GetCurrIndex() {
		p.Queue.Remove(index)
		p.Next()
	}
}

// player playlist command
func (p *Player) CreatePlayList(name string) {
	p.PlayListManager.Create(name)
}

// delete playlist
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
		if music.Path == filepath.Join(
			p.PlayListManager.PlayListPath,
			name,
			p.PlayListManager.PlayListSongs(name)[index],
		) {
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

func (p *Player) GetCurrMusicPosition() time.Duration {
	state := p.getPlayerState()
	if p.getPlayerState() == Stopped {
		return 0
	}
	if state == Paused {
		return p._getlMeta()._lcurrentPos
	}
	currentMusic := p.Queue.GetCurrMusic()
	if currentMusic == nil {
		return 0
	}
	fmt.Println("Curr music", currentMusic)
	return currentMusic.PositionD()
}

func (p *Player) GetCurrMusicDuration() time.Duration {
	if p.Queue.IsEmpty() {
		return 0
	}
	music := p.Queue.GetCurrMusic()
	if music == nil {
		return 0
	}
	if p.getPlayerState() == Paused {
		return p._getlMeta()._lcurrentDur
	}
	return music.DurationD()
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

func (p *Player) searchWorker(
	ctx context.Context,
	engine string,
	query string,
	musicChan chan shared.SearchResult,
	wg *sync.WaitGroup,
) {
	defer func() {
		log.Println("Search worker done for", engine, query)
		wg.Done()
	}()

	searchRes, err := p.Director.Search(engine, query, 5)
	if err != nil {
		log.Println("Failed to search for", query, ":", err)
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

	// for music := range musicChan {
	// 	musics = append(musics, music)
	// }

	for {
		select {
		case music, ok := <-musicChan:
			if !ok {
				return musics
			}
			musics = append(musics, music)
			p.removeTask(query)
		case <-ctx.Done():
			log.Println("Search timed out")
			return musics
		}
	}
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
		go p.PlayListManager.AddToPlayListFromOnline(name, query, whatIsThis, p)
	}
	return []shared.SearchResult{}
}

// if result is empty, it means it detects and plays the music if succeed other wise it returns the search results
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
			p.Queue.SetCurrIndex(index)
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

func (p *Player) SetTheme(theme string) {
	err := config.EditConfigField("theme", theme)
	if err != nil {
		log.Println("Failed to set theme", err)
	}
}

func (p *Player) GetTheme() string {
	return config.GetConfig().Theme
}

func (p *Player) GetPlayerStatus() shared.Status {
	return shared.Status{
		CurrMusicIndex:    p.Queue.GetCurrIndex(),
		CurrMusicPosition: p.GetCurrMusicPosition(),
		CurrMusicDuration: p.GetCurrMusicDuration(),
		PlayerState:       p.getPlayerState(),
		MusicQueue:        p.Queue.GetTitles(),
		Volume:            p.Vol,
		Tasks:             p.Tasks,
	}
}
