package shared

import (
	"fmt"
	"time"
)

const (
	Playing = iota
	Paused
	Stopped
)

type Status struct {
	CurrentMusicIndex    int
	CurrentMusicPosition time.Duration
	CurrentMusicLength   time.Duration
	PlayerState          int
	MusicList            []string
}

func (s Status) String() string {
	var str string
	str += "CurrentMusicIndex: " + fmt.Sprintf("%d", s.CurrentMusicIndex) + "\n"
	str += "CurrentMusicPosition: " + s.CurrentMusicPosition.String() + "\n"
	str += "CurrentMusicLength: " + s.CurrentMusicLength.String() + "\n"
	switch s.PlayerState {
	case Playing:
		str += "PlayerState: Playing\n"
	case Paused:
		str += "PlayerState: Paused\n"
	case Stopped:
		str += "PlayerState: Stopped\n"
	}

	str += "MusicList: " + "\n"
	str += "[\n"
	for _, music := range s.MusicList {
		str += "\t" + music + "\n"
	}
	str += "]"
	return str
}
