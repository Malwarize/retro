package views

import (
	"fmt"
	"net/rpc"
	"sync"
	"time"

	"github.com/Malwarize/goplay/controller"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type searchResultItem struct {
	title string
	desc  string
	ftype string
}

func (i searchResultItem) Title() string       { return i.title }
func (i searchResultItem) Description() string { return i.desc }
func (i searchResultItem) FilterValue() string { return i.title }

const (
	notStarted = iota
	running
	finished
)

type model struct {
	client *rpc.Client
	query  string

	selectList  list.Model
	spin        spinner.Model
	searchState int
	mu          *sync.Mutex
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spin.Tick,
		m.Search,
	)
}

func NewModel(client *rpc.Client, query string) *model {
	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = spinnerStyle
	return &model{
		client: client,
		query:  query,
		spin:   spin,
		mu:     &sync.Mutex{},
	}
}

func (m model) View() string {
	if m.searchState == finished {
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
			i := m.selectList.Index()
			controller.DetectAndPlay(m.selectList.Items()[i].(searchResultItem).desc, m.client)
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		fmt.Println("window size msg")
		h, v := docStyle.GetFrameSize()
		m.selectList.SetSize(msg.Width-h, msg.Height-v)
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
		m.selectList = list.New(msg.(searchDone).results, list.NewDefaultDelegate(), 50, 14)
		m.selectList.Title = "Select a song"
		m.searchState = finished
		m.selectList.SetFilteringEnabled(false)
		m.selectList.SetShowHelp(false)
		return selectUpdate(msg, m)
	}
	if m.searchState == finished {
		return selectUpdate(msg, m)
	}
	return spinnerUpdate(msg, m)
}

type searchDone struct {
	results []list.Item
}

func (m model) Search() tea.Msg {
	var results []list.Item
	musics := controller.DetectAndPlay(m.query, m.client)
	for _, music := range musics {
		results = append(results, searchResultItem{
			title: music.Title,
			desc:  music.Destination,
			ftype: music.Type,
		})
	}
	time.Sleep(1 * time.Second)
	return searchDone{
		results: results,
	}
}

func SearchThenSelect(query string, client *rpc.Client) error {
	model := NewModel(client, query)
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
