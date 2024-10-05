package config

import (
	"encoding/json"
	"errors"
	"fmt"
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
	configPath = os.Getenv("HOME") + "/.retro/config.json"
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
	ServerPort    string        `json:"server_port"`    // port to run the server on
}

// Merges file config with default config
func mergeConfigs(config, defaultConfig *Config) *Config {
	if config.RetroPath == "" {
		config.RetroPath = defaultConfig.RetroPath
	}
	if config.PathYTDL == "" {
		config.PathYTDL = defaultConfig.PathYTDL
	}
	if config.PathFFmpeg == "" {
		config.PathFFmpeg = defaultConfig.PathFFmpeg
	}
	if config.PathFFprobe == "" {
		config.PathFFprobe = defaultConfig.PathFFprobe
	}
	if config.SearchTimeout == 0 {
		config.SearchTimeout = defaultConfig.SearchTimeout
	}
	if config.Theme == "" {
		config.Theme = defaultConfig.Theme
	}
	if config.LogFile == "" {
		config.LogFile = defaultConfig.LogFile
	}
	if config.DBPath == "" {
		config.DBPath = defaultConfig.DBPath
	}
	if config.ServerPort == "" {
		config.ServerPort = defaultConfig.ServerPort
	}
	// No need to check boolean field (DiscordRPC) since false is a meaningful value
	return config
}

func initConfig() *Config {
	retro_path := os.Getenv("HOME") + "/.retro/"
	configPath := filepath.Join(retro_path, "config.json")
	var config *Config

	// Load default config
	defaultConfig := &Config{
		RetroPath:     retro_path,
		PathYTDL:      "yt-dlp",
		PathFFmpeg:    "ffmpeg",
		PathFFprobe:   "ffprobe",
		SearchTimeout: 60 * time.Second,
		Theme:         "pink",
		DiscordRPC:    true,
		LogFile:       filepath.Join(retro_path, "retro.log"),
		DBPath:        filepath.Join(retro_path, "retro.db"),
		ServerPort:    "3131",
	}

	// Attempt to load from file
	if jsonFile, err := os.ReadFile(configPath); err == nil {
		config = &Config{}
		if err = json.Unmarshal(jsonFile, config); err == nil {
			// Merge file config with default config
			config = mergeConfigs(config, defaultConfig)
			return config
		} else {
			fmt.Println("Error loading config file:", err)
		}
	}

	// If we can't load from file, return the default config
	return defaultConfig
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
	case "server_port":
		config.ServerPort = value
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
