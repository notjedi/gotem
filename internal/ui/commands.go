package ui

/* import (
	c "context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hekmon/transmissionrpc/v2"
	"github.com/notjedi/gotem/internal/context"
) */

var (
	torrentFields = []string{"status", "rateDownload",
		"rateUpload", "uploadedEver"}
)

/*
   TODO: should i make another RPC call for updating status bar?
   although the number of fields are low, since we are getting these
   fields in other view, i guess we can update them in a global context?
   for download and upload speeds, we don't have to depend on the update
   interval if we can keep a running sum of the net speeds during the
   interval. will come back to this once i'm done with the detail view,
   to make a better design choice.
*/

/* type statusbarUpdateMsg []statusbarInfo

func generateStatusbarUpdateMsg(ctx context.Context) tea.Msg {
	torrents, _ := ctx.Client().TorrentGet(c.TODO(), torrentFields, nil)
	return statusbarUpdateMsg(torrents)
}

func statusbarUpdateCmd(ctx context.Context) tea.Cmd {
	return tea.Tick(time.Second*time.Duration(1), func(t time.Time) tea.Msg {
		return generateStatusbarUpdateMsg(ctx)
	})
} */
