package listview

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func generateTorrentUpdateMsg(m Model) tea.Msg {
	torrents, _ := m.ctx.Client.TorrentGet(context.TODO(), torrentFields, nil)
	var items []list.Item
	for _, torrent := range torrents {
		items = append(items, TorrentItem{torrent, m.ctx.ListWidth})
	}
	return torrentUpdateMsg(items)
}

func (m Model) updateTorrentsCmd() tea.Cmd {
	return tea.Tick(time.Second*time.Duration(2), func(t time.Time) tea.Msg {
		return generateTorrentUpdateMsg(m)
	})
}
