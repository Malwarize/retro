package config

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

var (
	once sync.Once
	cfg  *Config // singleton instance
)

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
	DbPath        string        `json:"db_path"`   // path to the database
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
		if jsonConfig.DbPath == "" {
			jsonConfig.DbPath = os.Getenv("HOME") + "/.goplay/goplay.db"
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

func EditConfigField(field string, value string) error {
	jsonConfig := GetConfig()
	switch field {
	case "goplay_path":
		jsonConfig.GoPlayPath = value
	case "playlist_path":
		jsonConfig.PlaylistPath = value
	case "cache_dir":
		jsonConfig.CacheDir = value
	case "path_ytldpl":
		jsonConfig.Pathytldpl = value
	case "path_ffmpeg":
		jsonConfig.Pathffmpeg = value
	case "path_ffprobe":
		jsonConfig.Pathffprobe = value
	case "search_timeout":
		jsonConfig.SearchTimeOut, _ = time.ParseDuration(value)
	case "separator":
		jsonConfig.Separator = value
	case "theme":
		jsonConfig.Theme = value
	case "db_path":
		jsonConfig.DbPath = value
	default:
		return errors.New("Unknown field : " + field)
	}
	// write config file
	jsonData, err := json.MarshalIndent(jsonConfig, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(configPath, jsonData, 0o644)
	if err != nil {
		return err
	}
	// update the singleton instance
	cfg = jsonConfig
	return nil
}

func DebugConfig() *Config {
	homeDir := os.Getenv("HOME")
	return &Config{
		GoPlayPath:    "./goplay_storage/", // in the current directory
		PlaylistPath:  "./goplay_storage/playlists/",
		CacheDir:      "./goplay_storage/cache/",
		DbPath:        homeDir + "/.goplay/goplay.db",
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
		DbPath:        homeDir + "/.goplay/goplay.db",
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
