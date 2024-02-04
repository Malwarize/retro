package views

import (
	"net/rpc"

	"github.com/Malwarize/goplay/client/controller"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func addToPlayListCallback(m model) {
	i := m.selectList.Index()
	controller.DetectAndAddToPlayList(m.args[0].(string), m.selectList.Items()[i].(searchResultItem).desc, m.client)
}

func AddToPlayListQuitMessage(m model) string {
	return quitTextStyle.Render("🔋 Adding music " + m.selectList.Items()[m.selectList.Index()].(searchResultItem).title + " to playlist " + m.args[0].(string))
}

func (m model) AddSearch() tea.Msg {
	var results []list.Item
	musics := controller.DetectAndAddToPlayList(m.args[0].(string), m.query, m.client)
	for _, music := range musics {
		results = append(results, searchResultItem{
			title: music.Title,
			desc:  music.Destination,
			ftype: music.Type,
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
