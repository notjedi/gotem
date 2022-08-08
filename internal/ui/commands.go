package ui

import (
	c "context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hekmon/transmissionrpc/v2"
	"github.com/notjedi/gotem/internal/context"
)

var (
	torrentFields = []string{"status", "rateDownload", "rateUpload",
		"uploadedEver"}
)

type statusbarUpdateMsg []statusbarInfo

func generateStatusbarUpdateMsg(ctx context.Context) tea.Msg {
	torrents, _ := ctx.Client().TorrentGet(c.TODO(), torrentFields, nil)
	return statusbarUpdateMsg(torrents)
}

func statusbarUpdateCmd(ctx context.Context) tea.Cmd {
	return tea.Tick(time.Second*time.Duration(1), func(t time.Time) tea.Msg {
		return generateStatusbarUpdateMsg(ctx)
	})
}
