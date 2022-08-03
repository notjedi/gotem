package listview

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO: make this the same for both listview and detailview, accept ctx as arg instead of model, so
// we can use this for both views
func generateTorrentUpdateMsg(m Model) tea.Msg {
	torrents, _ := m.ctx.Client.TorrentGet(context.TODO(), torrentFields, nil)
	var items []list.Item
	// TODO: nitpick: use make() to predefine size of array, so we don't copy back and forth
	for _, torrent := range torrents {
		items = append(items, TorrentItem{torrent, m.ctx.TitleSpacing, m.ctx.DescSpacing})
	}
	return torrentUpdateMsg(items)
}

// TODO: rename to torrentUpdateCmd
// TODO: move this to listview.go? ig the code would be more readable then?
func (m Model) updateTorrentsCmd() tea.Cmd {
	return tea.Tick(time.Second*time.Duration(1), func(t time.Time) tea.Msg {
		return generateTorrentUpdateMsg(m)
	})
}
