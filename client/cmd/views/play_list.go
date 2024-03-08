package views

import (
	"net/rpc"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Malwarize/retro/client/controller"
	"github.com/Malwarize/retro/shared"
)

func addToPlayListCallback(m model) error {
	i := m.selectList.Index()
	_, err := controller.DetectAndAddToPlayList(
		m.args[0].(string),
		m.selectList.Items()[i].(searchResultItem).desc,
		m.client,
	)
	return err
}

func AddToPlayListQuitMessage(m model) string {
	return quitTextStyle.Render(
		"🔋 Adding music " + m.selectList.Items()[m.selectList.Index()].(searchResultItem).title + " to playlist " + m.args[0].(string),
	)
}

func (m model) AddSearch() tea.Msg {
	var results []list.Item
	musics, err := controller.DetectAndAddToPlayList(m.args[0].(string), m.query, m.client)
	if err != nil {
		return searchDone{nil, err}
	}
	for _, music := range musics {
		results = append(results, searchResultItem{
			title:    music.Title,
			desc:     music.Destination,
			ftype:    music.Type,
			duration: shared.DurationToString(music.Duration),
		})
	}
	return searchDone{
		results: results,
	}
}

func SearchThenAddToPlayList(playlist, query string, client *rpc.Client) error {
	model := NewModel(client, query)
	model.callback = addToPlayListCallback
	model.args = []any{playlist}
	model.quitMessage = AddToPlayListQuitMessage
	model.initCmd = model.AddSearch
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
