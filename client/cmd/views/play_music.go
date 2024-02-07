package views

import (
	"math/rand"
	"net/rpc"

	"github.com/Malwarize/goplay/client/controller"
	"github.com/Malwarize/goplay/shared"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func playCallback(m model) {
	i := m.selectList.Index()
	controller.DetectAndPlay(m.selectList.Items()[i].(searchResultItem).desc, m.client)
}

func PlayQuitMessage(m model) string {
	randEmoji := playingEmojies[rand.Intn(len(playingEmojies))]
	return quitTextStyle.Render(randEmoji + " Playing song " + m.selectList.Items()[m.selectList.Index()].(searchResultItem).title + ", this may take a while if download needed")
}

func (m model) PlaySearch() tea.Msg {
	var results []list.Item
	musics := controller.DetectAndPlay(m.query, m.client)
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

func SearchThenSelect(query string, client *rpc.Client) error {
	model := NewModel(client, query)
	p := tea.NewProgram(model)
	model.callback = playCallback
	model.quitMessage = PlayQuitMessage
	model.initCmd = model.PlaySearch
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
