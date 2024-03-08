package views

import (
	"math/rand"
	"net/rpc"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Malwarize/retro/client/controller"
	"github.com/Malwarize/retro/shared"
)

func playCallback(m model) error {
	i := m.selectList.Index()
	_, err := controller.DetectAndPlay(m.selectList.Items()[i].(searchResultItem).desc, m.client)
	return err
}

func PlayQuitMessage(m model) string {
	randEmoji := playingEmojies[rand.Intn(len(playingEmojies))]
	return quitTextStyle.Render(
		randEmoji + " Playing song " + m.selectList.Items()[m.selectList.Index()].(searchResultItem).title + ", this may take a while if download needed",
	)
}

func (m model) PlaySearch() tea.Msg {
	var results []list.Item
	musics, err := controller.DetectAndPlay(m.query, m.client)
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
