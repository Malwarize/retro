package views

import (
	"fmt"
	"net/rpc"
	"sync"

	"github.com/Malwarize/goplay/shared"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type searchResultItem struct {
	title    string
	desc     string
	ftype    string
	duration string
}

func (i searchResultItem) Title() string {
	if i.ftype == "cache" {
		return parseName(i.title)
	}
	return i.title
}
func (i searchResultItem) Description() string {
	return emojiesType[i.ftype] + " " + i.ftype + " " + i.duration
}
func (i searchResultItem) FilterValue() string { return "" }

type model struct {
	client *rpc.Client
	query  string

	// select callback to be called when a song is selected
	// func of reciever model
	callback func(model)

	initCmd     tea.Cmd
	quitMessage func(model) string
	args        []any
	selectList  list.Model
	spin        spinner.Model
	searchState int
	quit        bool
	mu          *sync.Mutex
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spin.Tick,
		m.initCmd,
	)
}

func NewModel(client *rpc.Client, query string) *model {
	spin := spinner.New()
	spin.Spinner = spinner.Points
	spin.Style = spinnerStyle
	return &model{
		client: client,
		query:  query,
		spin:   spin,
		mu:     &sync.Mutex{},
	}
}

func NewList(items []list.Item) list.Model {
	listModel := list.New(items, list.NewDefaultDelegate(), 50, 14)
	listModel.Title = "Select a song 👇"
	listModel.SetFilteringEnabled(false)
	listModel.SetShowHelp(false)
	listModel.Styles.Title = lipgloss.NewStyle()
	return listModel
}

func (m model) View() string {
	if m.quit {
		return m.quitMessage(m)
	}
	if m.searchState == shared.Finished {
		return docStyle.Render(m.selectList.View())
	}
	return fmt.Sprintf("%s Searching for %q...", m.spin.View(), m.query)
}

func spinnerUpdate(msg tea.Msg, m model) (model, tea.Cmd) {
	var cmd tea.Cmd
	m.spin, cmd = m.spin.Update(msg)
	return m, cmd
}

func selectUpdate(msg tea.Msg, m model) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			m.callback(m)
			m.quit = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.selectList, cmd = m.selectList.Update(msg)
	return m, cmd
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch t := msg.(type) {
	case tea.KeyMsg:
		switch t.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	case searchDone:
		if len(t.results) == 0 {
			return m, tea.Quit
		}
		m.selectList = NewList(t.results)
		m.searchState = shared.Finished
		return selectUpdate(msg, m)
	}
	if m.searchState == shared.Finished {
		return selectUpdate(msg, m)
	}
	return spinnerUpdate(msg, m)
}

type searchDone struct {
	results []list.Item
}
