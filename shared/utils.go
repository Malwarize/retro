package shared

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"
)

const (
	Playing = iota
	Paused
	Stopped
)

type Task struct {
	Target string // download "url", search "query"
	Type   string // download search
}

type Status struct {
	CurrentMusicIndex    int
	CurrentMusicPosition time.Duration
	CurrentMusicLength   time.Duration
	PlayerState          int
	MusicList            []string
	Tasks                []Task
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

	for _, task := range s.Tasks {
		str += "Task: " + task.Type + " " + task.Target + "\n"
	}

	return str
}

type SearchResult struct {
	Title       string
	Destination string
	Type        string
}

func EscapeSpecialDirChars(path string) string {
	// escape special chars
	path = url.PathEscape(path)
	//first 40 chars
	path = path[:40]
	return path
}

func ParseCachedFileName(filename string) (string, string) {
	// split filename by __
	split := strings.Split(filename, Separator)
	if len(split) != 2 {
		log.Println("Invalid cached file name: ", filename)
		return "", ""
	}
	return split[0], split[1]
}

func CombineNameWithKey(name string, key string) string {
	return name + Separator + key
}
