package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	once sync.Once
	cfg  *Config // singleton instance
)

var (
	DEBUG      = false // set to true for debug mode
	configPath = os.Getenv("HOME") + "/.config/retro.json"
)

type Config struct {
	RetroPath     string        `json:"retro_path"`     // path to retro
	PathYTDL      string        `json:"path_ytldpl"`    // path to yt-dlp
	PathFFmpeg    string        `json:"path_ffmpeg"`    // path to ffmpeg
	PathFFprobe   string        `json:"path_ffprobe"`   // path to ffprobe
	SearchTimeout time.Duration `json:"search_timeout"` // search timeout
	Theme         string        `json:"theme"`          // UI theme
	DBPath        string        `json:"db_path"`        // path to the database
	DiscordRPC    bool          `json:"discord_rpc"`    // Discord Rich Presence
	LogFile       string        `json:"log_file"`       // path to the log file
}

func initConfig() *Config {
	retro_path := os.Getenv("HOME") + "/.retro/"
	config := &Config{
		RetroPath:     retro_path,
		PathYTDL:      "yt-dlp",
		PathFFmpeg:    "ffmpeg",
		PathFFprobe:   "ffprobe",
		SearchTimeout: 60 * time.Second,
		Theme:         "pink",
		DiscordRPC:    true,
		LogFile:       filepath.Join(retro_path, "retro.log"),
		DBPath:        filepath.Join(retro_path, "retro.db"),
	}

	// Attempt to load from file
	if jsonFile, err := os.ReadFile(configPath); err == nil {
		if err = json.Unmarshal(jsonFile, config); err != nil {
			return config // Return default config if unmarshaling fails
		}
	}

	return config
}

func GetConfig() *Config {
	once.Do(func() {
		cfg = initConfig()
	})
	return cfg
}

func EditConfigField(field, value string) error {
	config := GetConfig()
	switch field {
	case "retro_path":
		config.RetroPath = value
	case "path_ytldpl":
		config.PathYTDL = value
	case "path_ffmpeg":
		config.PathFFmpeg = value
	case "path_ffprobe":
		config.PathFFprobe = value
	case "search_timeout":
		if duration, err := time.ParseDuration(value); err == nil {
			config.SearchTimeout = duration
		} else {
			return err
		}
	case "theme":
		config.Theme = value
	case "db_path":
		config.DBPath = value
	case "discord_rpc":
		if value == "true" {
			config.DiscordRPC = true
		} else if value == "false" {
			config.DiscordRPC = false
		}
	case "log_file":
		config.LogFile = value
	default:
		return errors.New("unknown field: " + field)
	}

	// Save updated config to file
	return saveConfig(config)
}

func saveConfig(config *Config) error {
	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, jsonData, 0o644)
}
