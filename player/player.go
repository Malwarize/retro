package player

import (
	"fmt"
	"sync"
	"time"

	"github.com/Malwarize/goplay/shared"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
)

var PlayerInstance *Player

type Music struct {
	Name     string
	Streamer beep.StreamSeekCloser
	Format   beep.Format
}

func (m Music) String() string {
	return m.Name
}

const (
	Playing = iota
	Paused
	Stopped
)

type Player struct {
	MusicList         []Music
	CurrentMusicIndex int
	playerState       int
	done              chan struct{}
	initialised       bool
	mu                sync.Mutex
}

func NewMusic(name string, streamer beep.StreamSeekCloser, format beep.Format) Music {
	return Music{
		Name:     name,
		Streamer: streamer,
		Format:   format,
	}
}

func NewPlayer() *Player {
	return &Player{
		MusicList:         make([]Music, 0),
		CurrentMusicIndex: 0,
		playerState:       Stopped,
		done:              make(chan struct{}),
		initialised:       false,
	}
}

func GetPlayer() *Player {
	if PlayerInstance == nil {
		PlayerInstance = NewPlayer()
	}
	return PlayerInstance
}

func (p *Player) AddMusic(music Music) {
	p.MusicList = append(p.MusicList, music)
}

func (p *Player) Play() {
	music := p.MusicList[p.CurrentMusicIndex]
	if p.initialised == false {
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
	if len(p.MusicList) < 1 || p.playerState == Stopped {
		return
	}
	if p.playerState == Paused {
		p.Resume()
	}
	p.MusicList[p.CurrentMusicIndex].Streamer.Seek(0)
	p.CurrentMusicIndex++
	if p.CurrentMusicIndex >= len(p.MusicList) {
		p.CurrentMusicIndex = 0
	}
	p.Play()
}

func (p *Player) Prev() {
	if len(p.MusicList) < 1 || p.playerState == Stopped {
		return
	}
	if p.playerState == Paused {
		p.Resume()
	}
	p.MusicList[p.CurrentMusicIndex].Streamer.Seek(0)
	p.CurrentMusicIndex--
	if p.CurrentMusicIndex < 0 {
		p.CurrentMusicIndex = len(p.MusicList) - 1
	}
	p.Play()
}

func (p *Player) Stop() {
	if p.playerState == Stopped {
		return
	}
	if p.playerState == Paused {
		p.Resume()
	}
	speaker.Clear()
	for _, music := range p.MusicList {
		music.Streamer.Close()
	}
	p.MusicList = make([]Music, 0)
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
	currentSamplePos := p.MusicList[p.CurrentMusicIndex].Streamer.Position()
	curretnTimePos := p.MusicList[p.CurrentMusicIndex].Format.SampleRate.D(currentSamplePos)
	newTimePos := (curretnTimePos + d) % p.GetCurrentMusicLength()
	newSamplePos := p.MusicList[p.CurrentMusicIndex].Format.SampleRate.N(newTimePos)
	if err := p.MusicList[p.CurrentMusicIndex].Streamer.Seek(newSamplePos); err != nil {
		fmt.Println(err)
	}
	speaker.Unlock()
}

func (p *Player) GetCurrentMusic() Music {
	return p.MusicList[p.CurrentMusicIndex]
}

func (p *Player) GetCurrentMusicPosition() time.Duration {
	if p.playerState == Stopped {
		return 0
	}
	currentSamplePos := p.MusicList[p.CurrentMusicIndex].Streamer.Position()
	curretnTimePos := p.MusicList[p.CurrentMusicIndex].Format.SampleRate.D(currentSamplePos)
	return curretnTimePos
}

func (p *Player) GetCurrentMusicLength() time.Duration {
	if len(p.MusicList) == 0 {
		return 0
	}
	music := p.GetCurrentMusic()
	return music.Format.SampleRate.D(music.Streamer.Len())
}

func (p *Player) getMusicList() []string {
	musicList := make([]string, 0)
	for _, music := range p.MusicList {
		musicList = append(musicList, music.Name)
	}
	return musicList
}

func (p *Player) GetPlayerStatus() shared.Status {
	return shared.Status{
		CurrentMusicIndex:    p.CurrentMusicIndex,
		CurrentMusicPosition: p.GetCurrentMusicPosition(),
		CurrentMusicLength:   p.GetCurrentMusicLength(),
		PlayerState:          p.playerState,
		MusicList:            p.getMusicList(),
	}
}

func (p *Player) Wait() {
	<-p.done
}
