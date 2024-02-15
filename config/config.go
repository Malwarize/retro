package config

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

var once sync.Once
var cfg *Config // singleton instance

var DEBUG = false // set to true for debug mode

type Config struct {
	GoPlayPath    string        `json:"goplay_path"`   // path to goplay
	PlaylistPath  string        `json:"playlist_path"` // path to playlists
	CacheDir      string        `json:"cache_dir"`     // path to cache
	Pathytldpl    string        `json:"path_ytldpl"`   // path to yt-dlp
	Pathffmpeg    string        `json:"path_ffmpeg"`   // path to ffmpeg
	Pathffprobe   string        `json:"path_ffprobe"`  // path to ffprobe
	SearchTimeOut time.Duration `json:"search_timeout"`
	Separator     string        `json:"separator"` // separator for file names cache file
	Theme         string        `json:"theme"`     // blue, purple, pink
}

var configPath = os.Getenv("HOME") + "/.config/goplay.json"

func loadConfig() *Config {
	// read config file
	if _, err := os.Stat(configPath); err == nil {
		jsonFile, err := os.ReadFile(configPath)
		if err != nil {
			return nil
		}
		var jsonConfig Config
		err = json.Unmarshal(jsonFile, &jsonConfig)
		if err != nil {
			return nil
		}
		// default values
		if jsonConfig.GoPlayPath == "" {
			jsonConfig.GoPlayPath = os.Getenv("HOME") + "/.goplay/"
		}
		if jsonConfig.PlaylistPath == "" {
			jsonConfig.PlaylistPath = os.Getenv("HOME") + "/.goplay/playlists/"
		}
		if jsonConfig.CacheDir == "" {
			jsonConfig.CacheDir = os.Getenv("HOME") + "/.goplay/cache/"
		}
		if jsonConfig.Pathytldpl == "" {
			jsonConfig.Pathytldpl = "yt-dlp"
		}
		if jsonConfig.Pathffmpeg == "" {
			jsonConfig.Pathffmpeg = "ffmpeg"
		}
		if jsonConfig.Pathffprobe == "" {
			jsonConfig.Pathffprobe = "ffprobe"
		}
		if jsonConfig.SearchTimeOut == 0 {
			jsonConfig.SearchTimeOut = 60 * time.Second
		}
		if jsonConfig.Separator == "" {
			jsonConfig.Separator = "_#__#_"
		}
		if jsonConfig.Theme == "" {
			jsonConfig.Theme = "pink"
		}

		return &jsonConfig
	}
	return nil
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
		Theme:         "pink",
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
		Theme:         "pink",
	}
}

func GetConfig() *Config {
	once.Do(func() {
		if DEBUG {
			cfg = DebugConfig()
		} else {
			if cfg = loadConfig(); cfg != nil {
				return
			}
			cfg = ReleaseConfig()
		}
	})
	return cfg
}
