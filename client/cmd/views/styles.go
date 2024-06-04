package views

import (
	"github.com/Malwarize/retro/client/controller"
	"github.com/Malwarize/retro/shared"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// Emoji and status mappings
var (
	emojiesType = map[string]string{
		"youtube": "🎬",
		"cache":   "💾",
		"file":    "🎵",
		"dir":     "📁",
	}

	playingEmojies = []string{
		"🎵",
		"🎶",
		"🎷",
		"🎸",
		"🎹",
		"🎺",
	}

	emojiesStatus = map[shared.PState]string{
		shared.Playing: "▶️",
		shared.Stopped: "🛑",
		shared.Paused:  "⏸️",
	}

	tasksEmojies = map[int]string{
		shared.Downloading: "📥",
		shared.Searching:   "🔍",
	}

	volumeLevels = []string{
		"🔇",
		"🔈",
		"🔉",
		"🔊",
	}

	failedEmojie  = "❌"
	defaultMargin = lipgloss.NewStyle().Margin(1, 2)
)

// Themes struct encapsulates all theme-related styles
type Themes struct {
	MainColorStyle   string
	DocStyle         lipgloss.Style
	QuitTextStyle    lipgloss.Style
	SpinnerStyle     lipgloss.Style
	ProgressStyle    lipgloss.Style
	ColoredTextStyle lipgloss.Style
	RunningStyle     lipgloss.Style
	StoppedStyle     lipgloss.Style
	PausedStyle      lipgloss.Style
	PositionStyle    lipgloss.Style
	SelectMusicStyle lipgloss.Style
	FailStyle        lipgloss.Style
	TaskStyle        lipgloss.Style
	ListDelegate     list.DefaultDelegate
}

// Helper function to create item styles
func createItemStyle(mainColor lipgloss.AdaptiveColor) (s list.DefaultItemStyles) {
	s.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
		Padding(0, 0, 0, 2)

	s.NormalDesc = s.NormalTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

	s.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(mainColor).
		Foreground(mainColor).
		Padding(0, 0, 0, 1)

	s.SelectedDesc = s.SelectedTitle.Copy().
		Foreground(mainColor)

	s.DimmedTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
		Padding(0, 0, 0, 2)

	s.DimmedDesc = s.DimmedTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#C2B8C2", Dark: "#4D4D4D"})

	s.FilterMatch = lipgloss.NewStyle().Underline(true)

	return s
}

// Create list delegate with specific color
func createListDelegate(mainColor lipgloss.AdaptiveColor) list.DefaultDelegate {
	def := list.NewDefaultDelegate()
	def.Styles = createItemStyle(mainColor)
	return def
}

// Common theme attributes
func commonStyles(mainColor lipgloss.AdaptiveColor) Themes {
	positionStyle := lipgloss.NewStyle().Margin(0, 0, 0, 3)
	colorStyle := lipgloss.NewStyle().Foreground(mainColor)
	return Themes{
		DocStyle:      defaultMargin,
		QuitTextStyle: lipgloss.NewStyle().Margin(1, 0, 2, 4),
		SpinnerStyle:  lipgloss.NewStyle().Foreground(mainColor),
		ProgressStyle: positionStyle.Copy().Inherit(colorStyle).MarginTop(1),
		RunningStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Margin(1, 0, 2, 3),
		StoppedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Margin(1, 0, 2, 3),
		PausedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#0000FF")).
			Margin(1, 0, 2, 3),
		PositionStyle:    positionStyle,
		SelectMusicStyle: lipgloss.NewStyle().Foreground(mainColor).Margin(0, 0, 0, 1).Bold(true),
		FailStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			Margin(1, 0, 0, 3),
		TaskStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFFF")).
			Margin(1, 0, 0, 3),
		ColoredTextStyle: colorStyle,
	}
}

// Purple theme settings
func NewPurpleTheme() Themes {
	purple := lipgloss.AdaptiveColor{Light: "#D8BFD8", Dark: "#800080"}
	theme := commonStyles(purple)
	theme.MainColorStyle = "#800080"
	theme.ListDelegate = createListDelegate(purple)
	return theme
}

// Pink theme settings
func NewPinkTheme() Themes {
	pink := lipgloss.AdaptiveColor{Light: "#FFC0CB", Dark: "#FF1493"}
	theme := commonStyles(pink)
	theme.MainColorStyle = "#FF1493"
	theme.ListDelegate = list.NewDefaultDelegate()
	return theme
}

// Blue theme settings
func NewBlueTheme() Themes {
	blue := lipgloss.AdaptiveColor{Light: "#ADD8E6", Dark: "#0000C7"}
	theme := commonStyles(blue)
	theme.MainColorStyle = "#0000C7"
	theme.ListDelegate = createListDelegate(blue)
	return theme
}

// GetTheme returns the theme based on client settings
func GetTheme() Themes {
	client, err := controller.GetClient()
	if client == nil || err != nil {
		return NewPinkTheme()
	}
	theme := controller.GetTheme(client)
	switch theme {
	case "purple":
		return NewPurpleTheme()
	case "blue":
		return NewBlueTheme()
	default:
		return NewPinkTheme()
	}
}
