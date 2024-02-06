package config

import (
	"os"
	"sync"
	"time"
)

var once sync.Once
var cfg *Config // singleton instance

var DEBUG = false // set to true for debug mode

type Config struct {
	GoPlayPath    string // used to store cache, playlists, etc
	PlaylistPath  string // path to playlists
	CacheDir      string // path to cache
	Pathytldpl    string // path to yt-dlp
	Pathffmpeg    string // path to ffmpeg
	Pathffprobe   string // path to ffprobe
	SearchTimeOut time.Duration
	Separator     string // separator for cache files filename_#__#_id
}

func DebugConfig() *Config {
	return &Config{
		GoPlayPath:    "./goplay_storage/", // in the current directory
		PlaylistPath:  "./goplay_storage/playlists/",
		CacheDir:      "./goplay_storage/cache/",
		Pathytldpl:    "yt-dlp",
		Pathffmpeg:    "ffmpeg",
		Pathffprobe:   "ffprobe",
		SearchTimeOut: 60 * time.Second,
		Separator:     "_#__#_",
	}
}

func ReleaseConfig() *Config {
	homeDir := os.Getenv("HOME")
	return &Config{
		GoPlayPath:    homeDir + "/.goplay/",
		PlaylistPath:  homeDir + "/.goplay/playlists/",
		CacheDir:      homeDir + "/.goplay/cache/",
		Pathytldpl:    "yt-dlp",
		Pathffmpeg:    "ffmpeg",
		Pathffprobe:   "ffprobe",
		SearchTimeOut: 60 * time.Second,
		Separator:     "_#__#_",
	}
}

func GetConfig() *Config {
	once.Do(func() {
		if DEBUG {
			cfg = DebugConfig()
		} else {
			cfg = ReleaseConfig()
		}
	})
	return cfg
}
