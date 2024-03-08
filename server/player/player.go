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

	"github.com/Malwarize/retro/config"
	"github.com/Malwarize/retro/logger"
	"github.com/Malwarize/retro/shared"
)

var Instance *Player

var once sync.Once

type lmeta struct {
	_lcurrentPos time.Duration
	_lcurrentDur time.Duration
}

type Player struct {
	Queue       *MusicQueue
	playerState shared.PState
	done        chan struct{}
	initialised bool
	Director    *Director
	Tasks       map[string]shared.Task
	Vol         uint8
	_lmeta      lmeta
	mu          sync.Mutex
}

func NewPlayer() *Player {
	if _, err := os.Stat(
		config.GetConfig().RetroPath,
	); os.IsNotExist(err) {
		logger.LogInfo("retro dir not found, creating it")
		err = os.Mkdir(config.GetConfig().RetroPath, 0o755)
		if err != nil {
			log.Fatal(err)
		}
	}
	director, err := NewDefaultDirector()
	if err != nil {
		log.Fatal(err)
	}

	return &Player{
		Queue:       NewMusicQueue(),
		playerState: shared.Stopped,
		done:        make(chan struct{}),
		initialised: false,
		Director:    director,
		Vol:         100,
		Tasks:       make(map[string]shared.Task),
	}
}

// meta field is to save the current position and duration of
// the music when paused because when paused the speaker is blocking
func (p *Player) _setlMeta(
	m lmeta,
) {
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
		Instance = NewPlayer()
	})
	return Instance
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
		err := speaker.Init(
			music.Format.SampleRate,
			music.Format.SampleRate.N(time.Second/10),
		)
		if err != nil {
			return err
		}
		p.initialised = true
	} else {
		speaker.Clear()
	}
	p.setPlayerState(
		shared.Playing,
	)
	go func() {
		done := make(
			chan struct{},
		)
		music.SetVolume(
			p.Vol,
		)
		speaker.Play(
			beep.Seq(music.Volume, beep.Callback(
				func() {
					done <- struct{}{}
				},
			)),
		)
		<-done
		p.Next()
	}()
	return nil
}

func (p *Player) getPlayerState() shared.PState {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.playerState
}

func (p *Player) setPlayerState(
	state shared.PState,
) {
	p.mu.Lock()
	p.playerState = state
	p.mu.Unlock()
	cur := p.Queue.GetCurrMusic()
	if cur == nil {
		go adjustDiscordRPC(p.playerState, "")
	} else {
		go adjustDiscordRPC(p.playerState, cur.Name)
	}
}

func (p *Player) Next() error {
	if p.Queue.IsEmpty() {
		return logger.LogError(
			logger.GError(
				"Queue is empty",
			),
		)
	}
	if p.getPlayerState() == shared.Stopped {
		return logger.LogError(
			logger.GError(
				"Player is stopped",
			),
		)
	}
	if p.getPlayerState() == shared.Paused {
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
	if state == shared.Stopped {
		return logger.LogError(
			logger.GError(
				"Player is stopped",
			),
		)
	}
	if state == shared.Paused {
		err := p.Resume()
		if err != nil {
			return err
		}
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
	if state == shared.Stopped {
		return nil
	}
	if state == shared.Paused {
		p.Resume()
	}
	speaker.Clear()
	p.Queue.Clear()
	p.setPlayerState(
		shared.Stopped,
	)
	return nil
}

func (p *Player) Pause() error {
	state := p.getPlayerState()
	if state == shared.Paused || state == shared.Stopped {
		return nil
	}
	p._setlMeta(
		lmeta{
			_lcurrentDur: p.GetCurrMusicDuration(),
			_lcurrentPos: p.GetCurrMusicPosition(),
		},
	)
	p.setPlayerState(shared.Paused)
	speaker.Lock()
	return nil
}

func (p *Player) Resume() error {
	state := p.getPlayerState()
	if state == shared.Playing || state == shared.Stopped {
		return nil
	}
	p.setPlayerState(shared.Playing)
	speaker.Unlock()
	return nil
}

func (p *Player) Seek(d time.Duration) error {
	state := p.getPlayerState()
	if state == shared.Stopped {
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
	if state == shared.Paused {
		err := p.Resume()
		if err != nil {
			return err
		}
		defer p.Pause()
	}
	if err := currentMusic.Seek(d); err != nil {
		logger.ERRORLogger.Println(
			err,
		)
		return logger.LogError(
			logger.GError(
				"Failed to seek",
			),
		)
	}
	return nil
}

func (p *Player) Volume(
	vp uint8,
) error {
	if p.getPlayerState() == shared.Stopped {
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

		if m.Name == p.Queue.GetCurrMusic().Name {
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
	// check if Exists
	_, err := p.Director.Db.GetPlaylist(
		plname,
	)

	if err == nil {
		return logger.LogError(
			logger.GError(
				"Playlist already exists",
			),
		)
	}
	err = p.Director.Db.AddPlaylist(
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
	pl, err := p.Director.Db.GetPlaylist(
		plname,
	)
	if err != nil {
		return logger.LogError(
			logger.GError(
				"Playlist does not exist",
			),
		)
	}
	err = p.Director.Db.RemovePlaylist(
		pl.Name,
	)
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

func (p *Player) PlayListsNames() ([]string, error) {
	lists, err := p.Director.Db.GetPlaylists()
	if err != nil {
		return nil, logger.LogError(
			logger.GError(
				"Failed to get playlists",
				err,
			),
		)
	}

	var names []string
	for _, list := range lists {
		names = append(names, list.Name)
	}
	return names, nil
}

func (p *Player) RemoveMusicFromPlayList(plname string, music shared.IntOrString) error {
	pl, err := p.Director.Db.GetPlaylist(
		plname,
	)
	if err != nil {
		return logger.LogError(
			logger.GError(
				"Playlist does not exist",
			),
		)
	}
	var m *Music
	ms, err := p.Director.Db.GetMusicsFromPlaylist(
		pl.Name,
	)
	if music.IsInt {
		index := music.IntVal
		if index < 0 || index >= len(ms) {
			return logger.LogError(
				logger.GError(
					"Index out of range",
				),
			)
		}
		err = p.Director.Db.RemoveMusicFromPlaylist(
			pl.Name,
			ms[index].Name,
		)
	} else {
		name := music.StrVal
		err = p.Director.Db.RemoveMusicFromPlaylist(
			pl.Name,
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
		if hash(music.Data) == hash(m.Data) {
			p.Queue.Remove(
				&music,
			)
		}
	}
	if err != nil {
		return logger.LogError(
			logger.GError(
				"Failed to remove song from playlist",
				err,
			),
		)
	}
	return nil
}

func (p *Player) GetPlayListMusicNames(plname string) ([]string, error) {
	pl, err := p.Director.Db.GetPlaylist(
		plname,
	)
	if err != nil {
		return nil, logger.LogError(
			logger.GError(
				"Failed to get playlist",
				err,
			),
		)
	}
	songs, err := p.Director.Db.GetMusicsFromPlaylist(
		pl.Name,
	)
	if err != nil {
		return nil, logger.LogError(
			logger.GError(
				"Failed to get songs from playlist",
				err,
			),
		)
	}
	var names []string
	logger.LogInfo(
		"Playlist name :",
		pl.Name,
	)
	for _, song := range songs {
		logger.LogInfo(
			"Music name :",
			song.Name,
		)
		names = append(names, song.Name)
	}
	return names, nil
}

func (p *Player) PlayListPlayMusic(plname string, music shared.IntOrString) error {
	pl, err := p.Director.Db.GetPlaylist(
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
	ms, err := p.Director.Db.GetMusicsFromPlaylist(
		pl.Name,
	)
	if err != nil {
		return logger.LogError(
			logger.GError(
				"Failed to get musics from playlist",
				err,
			),
		)
	}
	var m *Music
	if music.IsInt {
		index := music.IntVal
		if index < 0 || index >= len(ms) {
			return logger.LogError(
				logger.GError(
					"Index out of range",
				),
			)
		}
		m, err = NewMusic(
			ms[index].Name,
			ms[index].Data,
		)

	} else {
		name := music.StrVal
		for _, song := range ms {
			if song.Name == name {
				m, err = NewMusic(
					song.Name,
					song.Data,
				)
			}
		}
	}
	if err != nil || m == nil {
		return logger.LogError(
			logger.GError(
				"Error adding song to play list",
				err,
			),
		)
	}
	p.Queue.Enqueue(
		*m,
	)
	err = p.Play()
	if err != nil {
		return err
	}
	if p.getPlayerState() == shared.Stopped {
		err := p.Play()
		return err
	}
	return nil
}

func (p *Player) PlayListPlayAll(
	plname string,
) error {
	pl, err := p.Director.Db.GetPlaylist(
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
	ms, err := p.Director.Db.GetMusicsFromPlaylist(
		pl.Name,
	)
	if err != nil {
		return logger.LogError(
			logger.GError(
				"Failed to get musics from playlist",
				err,
			),
		)
	}

	for _, song := range ms {
		m, err := NewMusic(
			song.Name,
			song.Data,
		)
		if err != nil {
			logger.LogWarn(
				"skiping music",
				song.Name,
				err,
			)
		}
		p.Queue.Enqueue(
			*m,
		)
	}

	if p.getPlayerState() == shared.Stopped {
		err := p.Play()
		if err != nil {
			return err
		}
	}
	return nil
}

// ####################
// # Player Info      #
// ####################
func (p *Player) GetCurrMusicPosition() time.Duration {
	state := p.getPlayerState()
	if p.getPlayerState() == shared.Stopped {
		return 0
	}
	if state == shared.Paused {
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
	if p.getPlayerState() == shared.Paused {
		return p._getlMeta()._lcurrentDur
	}
	return music.DurationD()
}

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
