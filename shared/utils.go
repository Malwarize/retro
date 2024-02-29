package shared

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gopxl/beep/mp3"
)

type Task struct {
	Type  int // download, search
	Error string
}

type Status struct {
	CurrMusicIndex    int
	CurrMusicPosition time.Duration
	CurrMusicDuration time.Duration
	PlayerState       PState
	MusicQueue        []string
	Volume            uint8
	Tasks             map[string]Task // key: target, value: task
}

func (s Status) String() string {
	var str string
	str += "CurrMusicIndex: " + fmt.Sprintf("%d", s.CurrMusicIndex) + "\n"
	str += "CurrMusicPosition: " + s.CurrMusicPosition.String() + "\n"
	str += "CurrMusicLength: " + s.CurrMusicDuration.String() + "\n"
	switch s.PlayerState {
	case Playing:
		str += "PlayerState: Playing\n"
	case Paused:
		str += "PlayerState: Paused\n"
	case Stopped:
		str += "PlayerState: Stopped\n"
	}

	str += "Volume: " + fmt.Sprintf("%d", s.Volume) + "\n"

	str += "MusicQueue " + "\n"
	str += "[\n"
	for _, music := range s.MusicQueue {
		str += "\t" + music + "\n"
	}
	str += "]"

	for target, task := range s.Tasks {
		str += fmt.Sprintf("Target: %s, Type: %d, Error: %v\n", target, task.Type, task.Error)
	}

	return str
}

type SearchResult struct {
	Title       string
	Destination string
	Type        string
	Duration    time.Duration
}

type AddToPlayListArgs struct {
	PlayListName string
	Query        string
}

type RemoveMusicFromPlayListArgs struct {
	PlayListName string
	IndexOrName  IntOrString
}

type PlayListPlayMusicArgs struct {
	PlayListName string
	IndexOrName  IntOrString
}

// helper function to get mp3 duration
func GetMp3Duration(path string) (time.Duration, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	music, format, err := mp3.Decode(f)
	if err != nil {
		return 0, err
	}
	return format.SampleRate.D(music.Len()), nil
}

// function converts 00:00:00 to time.Duration
func StringToDuration(s string) (time.Duration, error) {
	sp := strings.Split(s, ":")
	if len(sp) < 2 {
		return 0, fmt.Errorf("invalid duration: %s", s)
	}
	l := len(sp)
	sec := "0"
	min := "0"
	hour := "0"
	if l > 0 {
		sec = sp[l-1]
	}

	if l > 1 {
		min = sp[l-2]
	}
	if l > 2 {
		hour = sp[l-3]
	}
	return time.ParseDuration(hour + "h" + min + "m" + sec + "s")
}

func DurationToString(d time.Duration) string {
	// to format 00:00:00
	return fmt.Sprintf("%02d:%02d:%02d", int(d.Hours()), int(d.Minutes())%60, int(d.Seconds())%60)
}

type IntOrString struct {
	IntVal int
	StrVal string
	IsInt  bool
}
