package common

import (
	c "context"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hekmon/transmissionrpc/v2"
	"github.com/notjedi/gotem/internal/context"
)

type AllTorrentInfoMsg []list.Item
type TorrentInfoMsg transmissionrpc.Torrent

func GenerateAllTorrentInfoMsg(ctx *context.ProgramContext) tea.Msg {
	torrents, _ := ctx.Client().TorrentGet(c.TODO(), torrentFields, nil)
	var items []list.Item
	// TODO: nitpick: use make() to predefine size of array, so we don't copy back and forth
	for _, torrent := range torrents {
		items = append(items, TorrentItem{torrent, ctx})
	}
	return AllTorrentInfoMsg(items)
}

func AllTorrentInfoCmd(ctx *context.ProgramContext) tea.Cmd {
	return tea.Tick(time.Second*time.Duration(1), func(t time.Time) tea.Msg {
		return GenerateAllTorrentInfoMsg(ctx)
	})
}

func GenerateTorrentInfoMsg(ctx *context.ProgramContext, id int64) tea.Msg {
	torrentInfo, _ := ctx.Client().TorrentGet(c.TODO(), torrentFields, []int64{id})
	return TorrentInfoMsg(torrentInfo[0])
}

func TorrentInfoCmd(ctx *context.ProgramContext, id int64) tea.Cmd {
	return tea.Tick(time.Second*time.Duration(1), func(t time.Time) tea.Msg {
		return GenerateTorrentInfoMsg(ctx, id)
	})
}
