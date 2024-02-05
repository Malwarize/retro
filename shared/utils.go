package shared

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"
)

type Task struct {
	Type  int // download, search
	Error string
}

type Status struct {
	CurrentMusicIndex    int
	CurrentMusicPosition time.Duration
	CurrentMusicLength   time.Duration
	PlayerState          int
	MusicQueue           []string
	Tasks                map[string]Task // key: target, value: task
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

type AddToPlayListArgs struct {
	PlayListName string
	Query        string
}

type RemoveSongFromPlayListArgs struct {
	PlayListName string
	Index        int
}
type PlayListPlaySongArgs struct {
	PlayListName string
	Index        int
}
