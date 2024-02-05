package views

import (
	"github.com/Malwarize/goplay/shared"
	"github.com/charmbracelet/lipgloss"
)

var emojiesType = map[string]string{
	"youtube": "🎬",
	"cache":   "💾",
	"file":    "🎵",
	"dir":     "📁",
}

var playingEmojies = []string{
	"🎵",
	"🎶",
	"🎷",
	"🎸",
	"🎹",
	"🎺",
}

var emojiesStatus = map[int]string{
	shared.Playing: "▶️",
	shared.Stopped: "🛑",
	shared.Paused:  "⏸️",
}

var tasksEmojies = map[int]string{
	shared.Download: "📥",
	shared.Search:   "🔍",
}

var failedEmojie = "❌"
var docStyle = lipgloss.NewStyle().Margin(1, 2)
var quitTextStyle = lipgloss.NewStyle().Margin(1, 0, 2, 4)
var spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

var progressStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Margin(1, 0, 0, 3)
var runningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")).Margin(1, 0, 2, 3)
var stoppedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Margin(1, 0, 2, 3)
var pausedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#0000FF")).Margin(1, 0, 2, 3)

var positionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Margin(1, 0, 0, 3)
var selectMusicStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Margin(0, 0, 0, 1)
var durationStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Margin(1, 0, 0, 3)
var taskStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")).Margin(1, 0, 0, 3)
var failedtaskStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Margin(1, 0, 0, 3)

var playListNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Margin(1, 0, 0, 3)
