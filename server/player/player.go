package player

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"

	"github.com/Malwarize/goplay/config"
	"github.com/Malwarize/goplay/logger"
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
		logger.LogInfo("goplay dir not found, creating it")
		err = os.Mkdir(config.GetConfig().GoPlayPath, 0o755)
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

// meta field is to save the current position and duration of
// the music when paused because when paused the speaker is blocking
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

func GetPlayer() *Player {
	once.Do(func() {
		PlayerInstance = NewPlayer()
	})
	return PlayerInstance
}

// ############################
// # Core Player methods      #
// ############################
func (p *Player) Play() error {
	if p.Queue.IsEmpty() {
		return logger.LogError(
			logger.GError(
				"Queue is empty",
			),
		)
	}
	music := p.Queue.GetCurrMusic()
	if music == nil {
		return logger.LogError(
			logger.GError(
				"Failed to get current music",
			),
		)
	}

	if !p.initialised {
		speaker.Init(
			music.Format.SampleRate,
			music.Format.SampleRate.N(time.Second/10),
		)
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
	return nil
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

func (p *Player) Next() error {
	if p.Queue.IsEmpty() {
		return logger.LogError(
			logger.GError(
				"Queue is empty",
			),
		)
	}
	if p.getPlayerState() == Stopped {
		return logger.LogError(
			logger.GError(
				"Player is stopped",
			),
		)
	}
	if p.getPlayerState() == Paused {
		p.Resume()
	}
	currentMusic := p.Queue.GetCurrMusic()
	if currentMusic == nil {
		return logger.LogError(
			logger.GError(
				"Failed to get current music",
			),
		)
	}
	if err := currentMusic.SetPositionD(0); err != nil {
		return logger.LogError(
			logger.GError(
				"Failed to set position",
				err,
			),
		)
	}
	p.Queue.QueueNext()
	err := p.Play()
	if err != nil {
		return err
	}
	return nil
}

func (p *Player) Prev() error {
	state := p.getPlayerState()
	if p.Queue.IsEmpty() {
		return logger.LogError(
			logger.GError(
				"Queue is empty",
			),
		)
	}
	if state == Stopped {
		return logger.LogError(
			logger.GError(
				"Player is stopped",
			),
		)
	}
	if state == Paused {
		p.Resume()
	}
	currentMusic := p.Queue.GetCurrMusic()
	if currentMusic == nil {
		return logger.LogError(
			logger.GError(
				"Failed to get current music",
			),
		)
	}
	if err := currentMusic.SetPositionD(0); err != nil {
		logger.LogError(
			logger.GError(
				"Failed to seek",
				err,
			),
		)
	}
	p.Queue.QueuePrev()
	err := p.Play()
	if err != nil {
		return err
	}
	return nil
}

func (p *Player) Stop() error {
	clear(p.Tasks)

	state := p.getPlayerState()
	if state == Stopped {
		return nil
	}
	if state == Paused {
		p.Resume()
	}
	speaker.Clear()
	p.Queue.Clear()
	p.setPlayerState(Stopped)
	return nil
}

func (p *Player) Pause() error {
	state := p.getPlayerState()
	if state == Paused || state == Stopped {
		return nil
	}
	p._setlMeta(lmeta{
		_lcurrentDur: p.GetCurrMusicDuration(),
		_lcurrentPos: p.GetCurrMusicPosition(),
	})
	p.setPlayerState(Paused)
	speaker.Lock()
	return nil
}

func (p *Player) Resume() error {
	state := p.getPlayerState()
	if state == Playing || state == Stopped {
		return nil
	}
	p.setPlayerState(Playing)
	speaker.Unlock()
	return nil
}

func (p *Player) Seek(d time.Duration) error {
	state := p.getPlayerState()
	if state == Stopped {
		return logger.LogError(
			logger.GError(
				"player is not running",
			),
		)
	}
	if p.Queue.IsEmpty() {
		return logger.LogError(
			logger.GError(
				"Queue is empty",
			),
		)
	}
	currentMusic := p.Queue.GetCurrMusic()
	if currentMusic == nil {
		return logger.LogError(
			logger.GError(
				"Failed to get current music",
			),
		)
	}
	if state == Paused {
		p.Resume()
		defer p.Pause()
	}
	if err := currentMusic.Seek(d); err != nil {
		logger.ERRORLogger.Println(err)
		return logger.LogError(
			logger.GError(
				"Failed to seek",
			),
		)
	}
	return nil
}

func (p *Player) Volume(vp int /*volume percentage*/) error {
	if p.getPlayerState() == Stopped {
		return nil
	}
	p.Vol = vp
	currentMusic := p.Queue.GetCurrMusic()
	if currentMusic == nil {
		return nil
	}
	speaker.Lock()
	currentMusic.SetVolume(vp)
	speaker.Unlock()
	return nil
}

func (p *Player) Remove(music shared.IntOrString) error {
	if p.Queue.IsEmpty() {
		return logger.LogError(
			logger.GError(
				"Queue is empty",
			),
		)
	}
	if p.Queue.Size() == 1 {
		p.Stop()
		return nil
	} else {
		var m *Music
		if music.IsInt {
			musicIndex := music.IntVal
			logger.LogInfo(
				"Removing music by index",
				strconv.Itoa(
					musicIndex,
				),
			)
			m = p.Queue.GetMusicByIndex(
				musicIndex,
			)
		} else {
			musicName := music.StrVal
			fmt.Println(musicName)
			m = p.Queue.GetMusicByName(
				musicName,
			)
			logger.LogInfo(
				"Removing music by name",
				musicName,
			)
		}
		if m == nil {
			return logger.LogError(
				logger.GError(
					"Music not found",
				),
			)
		}

		if m.Name() == p.Queue.GetCurrMusic().Name() {
			p.Queue.Remove(
				m,
			)
			return p.Next()
		}

		p.Queue.Remove(
			m,
		)
	}
	return nil
}

// ####################
// # Playlist methods #
// ####################

func (p *Player) CreatePlayList(plname string) error {
	// check if exists
	if p.PlayListManager.Exists(plname) {
		return logger.LogError(
			logger.GError(
				"Playlist already exists",
			),
		)
	}
	err := p.PlayListManager.CreatePlayList(
		plname,
	)
	if err != nil {
		logger.ERRORLogger.Println(err)
		return logger.LogError(
			logger.GError(
				"Failed to create playlist",
			),
		)
	}
	return nil
}

func (p *Player) RemovePlayList(plname string) error {
	if !p.PlayListManager.Exists(plname) {
		return logger.LogError(
			logger.GError(
				"Playlist does not exist",
			),
		)
	}
	pl, err := p.PlayListManager.GetPlayListByName(
		plname,
	)
	if err != nil {
		return logger.LogError(
			logger.GError(
				"Failed to get playlist",
				err,
			),
		)
	}
	err = p.PlayListManager.Remove(pl)
	if err != nil {
		logger.ERRORLogger.Println(err)
		return logger.LogError(
			logger.GError(
				"Failed to remove playlist",
				err,
			),
		)
	}
	return nil
}

func (p *Player) PlayListsNames() []string {
	return p.PlayListManager.PlayListsNames()
}

func (p *Player) RemoveSongFromPlayList(plname string, music shared.IntOrString) error {
	pl, err := p.PlayListManager.GetPlayListByName(plname)
	if err != nil {
		return logger.LogError(
			logger.GError(
				"Failed to get playlist",
				err,
			),
		)
	}
	var m *Music
	if music.IsInt {
		index := music.IntVal
		m, err = p.PlayListManager.GetPlayListSongByIndex(
			pl,
			index,
		)
	} else {
		name := music.StrVal
		m, err = p.PlayListManager.GetPlayListSongByName(
			pl,
			name,
		)
	}

	if err != nil {
		return logger.LogError(
			logger.GError(
				"Failed to get song from playlist",
				err,
			),
		)
	}

	// check if the exists in the queue and remove it
	for _, music := range p.Queue.queue {
		if music.Path == m.Path {
			p.Queue.Remove(
				&music,
			)
		}
	}
	err = p.PlayListManager.RemoveMusic(
		pl,
		m,
	)
	if err != nil {
		logger.ERRORLogger.Println(err)
		return logger.LogError(
			logger.GError(
				"Failed to remove song from playlist",
				err,
			),
		)
	}
	return nil
}

func (p *Player) GetPlayListSongNames(plname string) ([]string, error) {
	pl, err := p.PlayListManager.GetPlayListByName(plname)
	if err != nil {
		return nil, logger.LogError(
			logger.GError(
				"Failed to get playlist",
				err,
			),
		)
	}
	songs := p.PlayListManager.GetPlayListSongs(pl)
	var names []string
	logger.LogInfo(
		"Playlist name :",
		pl.Name,
	)
	for _, song := range songs {
		logger.LogInfo(
			"Song name :",
			song.Name(),
		)
		names = append(names, song.Name())
	}
	return names, nil
}

func (p *Player) PlayListPlaySong(plname string, indexOrName shared.IntOrString) error {
	pl, err := p.PlayListManager.GetPlayListByName(plname)
	if err != nil {
		return logger.LogError(
			logger.GError(
				"Failed to get playlist",
				err,
			),
		)
	}

	if indexOrName.IsInt {
		err = p.addMusicFromPlaylistByIndex(
			pl,
			indexOrName.IntVal,
		)
	} else {
		err = p.addMusicFromPlaylistByName(
			pl,
			indexOrName.StrVal,
		)
	}
	if err != nil {
		return logger.LogError(
			logger.GError(
				"Error adding song to play list",
				err,
			),
		)
	}
	err = p.Play()
	if err != nil {
		return err
	}
	if p.getPlayerState() == Stopped {
		err := p.Play()
		return err
	}
	return nil
}

func (p *Player) PlayListPlayAll(
	plname string,
) error {
	pl, err := p.PlayListManager.GetPlayListByName(plname)
	if err != nil {
		return logger.LogError(
			logger.GError(
				"Failed to get playlist",
				err,
			),
		)
	}
	err = p.addMusicsFromPlaylist(pl)
	if err != nil {
		return logger.LogError(
			logger.GError(
				"Failed to add musics from playlist",
				err,
			),
		)
	}
	if p.getPlayerState() == Stopped {
		err := p.Play()
		if err != nil {
			return err
		}
	}
	return nil
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

// ####################
// # Player Info      #
// ####################
func (p *Player) GetTheme() string {
	return config.GetConfig().Theme
}

func (p *Player) SetTheme(theme string) error {
	err := config.EditConfigField(
		"theme",
		theme,
	)
	if err != nil {
		return logger.LogError(
			logger.GError(
				"Failed to set theme",
				err,
			),
		)
	}
	return nil
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
