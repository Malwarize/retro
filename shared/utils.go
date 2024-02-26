package shared

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gopxl/beep/mp3"

	"github.com/Malwarize/goplay/config"
)

type Task struct {
	Type  int // download, search
	Error string
}

type Status struct {
	CurrMusicIndex    int
	CurrMusicPosition time.Duration
	CurrMusicDuration time.Duration
	PlayerState       int
	MusicQueue        []string
	Volume            int
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

func EscapeSpecialDirChars(path string) string {
	// escape special chars
	path = url.PathEscape(path)
	// first 40 chars
	path = path[:40]
	return path
}

func ParseCachedFileName(filename string) (string, string) {
	// split filename by __
	split := strings.Split(filename, config.GetConfig().Separator)
	if len(split) != 2 {
		return "", ""
	}
	return split[0], split[1]
}

func CombineNameWithKey(name string, key string) string {
	return name + config.GetConfig().Separator + key
}

type AddToPlayListArgs struct {
	PlayListName string
	Query        string
}

type RemoveSongFromPlayListArgs struct {
	PlayListName string
	IndexOrName  IntOrString
}

type PlayListPlaySongArgs struct {
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

func ViewParseName(name string) string {
	name = filepath.Base(name)
	if strings.Contains(name, config.GetConfig().Separator) {
		name = strings.Split(name, config.GetConfig().Separator)[0]
	}
	return name
}

type IntOrString struct {
	IntVal int
	StrVal string
	IsInt  bool
}
