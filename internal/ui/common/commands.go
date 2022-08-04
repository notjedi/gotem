package common

import (
	c "context"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/notjedi/gotem/internal/context"
)

type TorrentInfoMsg []list.Item

// TODO: make this the same for both listview and detailview, accept ctx as arg instead of model, so
// we can use this for both views
func GenerateTorrentInfoMsg(ctx context.Context) tea.Msg {
	torrents, _ := ctx.Client().TorrentGet(c.TODO(), torrentFields, nil)
	var items []list.Item
	// TODO: nitpick: use make() to predefine size of array, so we don't copy back and forth
	for _, torrent := range torrents {
		items = append(items, TorrentItem{torrent, ctx.TitleSpacing(), ctx.DescSpacing()})
	}
	return TorrentInfoMsg(items)
}

// TODO: rename to torrentUpdateCmd
// TODO: move this to listview.go? ig the code would be more readable then?
func TorrentInfoCmd(ctx context.Context) tea.Cmd {
	return tea.Tick(time.Second*time.Duration(1), func(t time.Time) tea.Msg {
		return GenerateTorrentInfoMsg(ctx)
	})
}
