package common

import (
	c "context"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hekmon/transmissionrpc/v2"
	"github.com/notjedi/gotem/internal/context"
)

type (
	AllTorrentInfoMsg []list.Item
	TorrentInfoMsg    transmissionrpc.Torrent
	SessionStatsMsg   transmissionrpc.SessionStats
)

func GenerateAllTorrentInfoMsg(ctx *context.ProgramContext) tea.Msg {
	torrents, _ := ctx.Client().TorrentGet(c.TODO(), allTorrentInfoFields, nil)
	items := make([]list.Item, 0, len(torrents))
	for _, torrent := range torrents {
		items = append(items, TorrentItem{torrent, ctx})
	}
	return AllTorrentInfoMsg(items)
}

func AllTorrentInfoCmd(ctx *context.ProgramContext) tea.Cmd {
	return tea.Tick(time.Second*time.Duration(2), func(t time.Time) tea.Msg {
		return GenerateAllTorrentInfoMsg(ctx)
	})
}

func AllTorrentInfoMsgInstant(ctx *context.ProgramContext) tea.Cmd {
	return func() tea.Msg {
		return GenerateAllTorrentInfoMsg(ctx)
	}
}

func GenerateTorrentInfoMsg(ctx *context.ProgramContext, id int64) tea.Msg {
	torrentInfo, _ := ctx.Client().TorrentGet(c.TODO(), torrentInfoFields, []int64{id})
	return TorrentInfoMsg(torrentInfo[0])
}

func TorrentInfoCmd(ctx *context.ProgramContext, id int64) tea.Cmd {
	return tea.Tick(time.Second*time.Duration(2), func(t time.Time) tea.Msg {
		return GenerateTorrentInfoMsg(ctx, id)
	})
}

func TorrentInfoCmdInstant(ctx *context.ProgramContext, id int64) tea.Cmd {
	return func() tea.Msg {
		return GenerateTorrentInfoMsg(ctx, id)
	}
}

func GenerateSessionStatsMsg(ctx *context.ProgramContext) tea.Msg {
	sessionStats, _ := ctx.Client().SessionStats(c.TODO())
	return SessionStatsMsg(sessionStats)
}

func SessionStatsCmd(ctx *context.ProgramContext) tea.Cmd {
	return tea.Tick(time.Second*time.Duration(2), func(t time.Time) tea.Msg {
		return GenerateSessionStatsMsg(ctx)
	})
}

func SessionStatsCmdInstant(ctx *context.ProgramContext) tea.Cmd {
	return func() tea.Msg {
		return GenerateSessionStatsMsg(ctx)
	}
}
