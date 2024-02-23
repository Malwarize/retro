package views

import (
	"github.com/Malwarize/goplay/client/controller"
	"github.com/Malwarize/goplay/shared"
	"github.com/charmbracelet/bubbles/list"
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
	shared.Downloading: "📥",
	shared.Searching:   "🔍",
}

var volumeLevels = []string{
	"🔇",
	"🔈",
	"🔉",
	"🔊",
}

var failedEmojie = "❌"

var docStyle = lipgloss.NewStyle().Margin(1, 2)

var quitTextStyle = lipgloss.NewStyle().Margin(1, 0, 2, 4)

// var spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

// var progressStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Margin(1, 0, 0, 3)
// var runningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")).Margin(1, 0, 2, 3)
// var stoppedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Margin(1, 0, 2, 3)
// var pausedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#0000FF")).Margin(1, 0, 2, 3)

// var positionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Margin(1, 0, 0, 3)
// var selectMusicStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Margin(0, 0, 0, 1)
// var durationStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Margin(1, 0, 0, 3)
// var taskStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")).Margin(1, 0, 0, 3)
// var failedtaskStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")).Margin(1, 0, 0, 3)

// var playListNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Margin(1, 0, 0, 3)

type Themes struct {
	DocStyle         lipgloss.Style
	QuitTextStyle    lipgloss.Style
	SpinnerStyle     lipgloss.Style
	ProgressStyle    lipgloss.Style
	RunningStyle     lipgloss.Style
	StoppedStyle     lipgloss.Style
	PausedStyle      lipgloss.Style
	PositionStyle    lipgloss.Style
	SelectMusicStyle lipgloss.Style
	FailStyle        lipgloss.Style
	TaskStyle        lipgloss.Style
	MainColor        string
	ListDelegate     list.DefaultDelegate
}

func purpleItemStyle() (s list.DefaultItemStyles) {
	purple := lipgloss.AdaptiveColor{Light: "#D8BFD8", Dark: "#800080"}

	s.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
		Padding(0, 0, 0, 2)

	s.NormalDesc = s.NormalTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

	s.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(purple).
		Foreground(purple).
		Padding(0, 0, 0, 1)

	s.SelectedDesc = s.SelectedTitle.Copy().
		Foreground(purple)

	s.DimmedTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
		Padding(0, 0, 0, 2)

	s.DimmedDesc = s.DimmedTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#C2B8C2", Dark: "#4D4D4D"})

	s.FilterMatch = lipgloss.NewStyle().Underline(true)
	return s
}

func ListPurpleDelegate() list.DefaultDelegate {
	def := list.NewDefaultDelegate()
	def.Styles = purpleItemStyle()
	return def
}

func blueItemStyle() (s list.DefaultItemStyles) {
	blue := lipgloss.AdaptiveColor{Light: "#ADD8E6", Dark: "#0000FF"}

	s.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
		Padding(0, 0, 0, 2)

	s.NormalDesc = s.NormalTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

	s.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(blue).
		Foreground(blue).
		Padding(0, 0, 0, 1)

	s.SelectedDesc = s.SelectedTitle.Copy().
		Foreground(blue)

	s.DimmedTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
		Padding(0, 0, 0, 2)

	s.DimmedDesc = s.DimmedTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#C2B8C2", Dark: "#4D4D4D"})

	s.FilterMatch = lipgloss.NewStyle().Underline(true)

	return s
}

func ListBlueDelegate() list.DefaultDelegate {
	def := list.NewDefaultDelegate()
	def.Styles = blueItemStyle()
	return def
}

func NewPurpleTheme() Themes {
	purple := lipgloss.AdaptiveColor{Light: "#D8BFD8", Dark: "#800080"}
	return Themes{
		DocStyle:      lipgloss.NewStyle().Margin(1, 2),
		QuitTextStyle: lipgloss.NewStyle().Margin(1, 0, 2, 4),
		SpinnerStyle:  lipgloss.NewStyle().Foreground(purple),
		ProgressStyle: lipgloss.NewStyle().Foreground(purple).Margin(1, 0, 0, 3),
		RunningStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Margin(1, 0, 2, 3),
		StoppedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Margin(1, 0, 2, 3),
		PausedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#0000FF")).
			Margin(1, 0, 2, 3),
		PositionStyle:    lipgloss.NewStyle().Foreground(purple).Margin(1, 0, 0, 3),
		SelectMusicStyle: lipgloss.NewStyle().Foreground(purple).Margin(0, 0, 0, 1),
		FailStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			Margin(1, 0, 0, 3),
		TaskStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFFF")).
			Margin(1, 0, 0, 3),
		MainColor:    "#A020F0",
		ListDelegate: ListPurpleDelegate(),
	}
}

func NewPinkTheme() Themes {
	return Themes{
		DocStyle:      lipgloss.NewStyle().Margin(1, 2),
		QuitTextStyle: lipgloss.NewStyle().Margin(1, 0, 2, 4),
		SpinnerStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("205")),
		ProgressStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Margin(1, 0, 0, 3),
		RunningStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFF00")).
			Margin(1, 0, 2, 3),
		StoppedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff0000")).
			Margin(1, 0, 2, 3),
		PausedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#0000FF")).
			Margin(1, 0, 2, 3),
		PositionStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Margin(1, 0, 0, 3),
		SelectMusicStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Margin(0, 0, 0, 1),
		FailStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			Margin(1, 0, 0, 3),
		TaskStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFFF")).
			Margin(1, 0, 0, 3),
		MainColor:    "205",
		ListDelegate: list.NewDefaultDelegate(),
	}
}

func NewBlueTheme() Themes {
	blue := lipgloss.AdaptiveColor{Light: "#ADD8E6", Dark: "#0000FF"}
	return Themes{
		DocStyle:      lipgloss.NewStyle().Margin(1, 2),
		QuitTextStyle: lipgloss.NewStyle().Margin(1, 0, 2, 4),
		SpinnerStyle:  lipgloss.NewStyle().Foreground(blue),
		ProgressStyle: lipgloss.NewStyle().Foreground(blue).Margin(1, 0, 0, 3),
		RunningStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Margin(1, 0, 2, 3),
		StoppedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Margin(1, 0, 2, 3),
		PausedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#0000FF")).
			Margin(1, 0, 2, 3),
		PositionStyle:    lipgloss.NewStyle().Foreground(blue).Margin(1, 0, 0, 3),
		SelectMusicStyle: lipgloss.NewStyle().Foreground(blue).Margin(0, 0, 0, 1),
		FailStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			Margin(1, 0, 0, 3),
		TaskStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFFF")).
			Margin(1, 0, 0, 3),
		MainColor:    "#0000FF",
		ListDelegate: ListBlueDelegate(),
	}
}

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
