package ui

import (
	"fmt"
	"math"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/hekmon/transmissionrpc/v2"
	"github.com/mistakenelf/teacup/statusbar"
	"github.com/notjedi/gotem/internal/config"
	"github.com/notjedi/gotem/internal/context"
	"github.com/notjedi/gotem/internal/theme"
	"github.com/notjedi/gotem/internal/ui/common"
	"github.com/notjedi/gotem/internal/ui/components/detailview"
	"github.com/notjedi/gotem/internal/ui/components/listview"
	"github.com/notjedi/tabs"
)

type (
	View  int32
	Model struct {
		currView   View
		ctx        *context.ProgramContext
		listView   listview.Model
		detailView detailview.Model
		statusbar  statusbar.Model
	}
)

const (
	TorrentListView View = iota + 1
	TorrentDetailView
	HelpView
)

var (
	appStyle     = lipgloss.NewStyle().Margin(1, 2, 1, 2)
	sessionStats common.SessionStatsMsg
)

func New(ctx *context.ProgramContext) Model {
	theme := theme.GetTheme("default")
	statusbarModel := statusbar.New(
		statusbar.ColorConfig{
			Foreground: theme.StatusbarSelectedFileForegroundColor,
			Background: theme.StatusbarSelectedFileBackgroundColor,
		},
		statusbar.ColorConfig{
			Foreground: theme.StatusbarBarForegroundColor,
			Background: theme.StatusbarBarBackgroundColor,
		},
		statusbar.ColorConfig{
			Foreground: theme.StatusbarTotalFilesForegroundColor,
			Background: theme.StatusbarTotalFilesBackgroundColor,
		},
		statusbar.ColorConfig{
			Foreground: theme.StatusbarLogoForegroundColor,
			Background: theme.StatusbarLogoBackgroundColor,
		},
	)

	listViewModel := listview.New(ctx)

	return Model{
		currView:  TorrentListView,
		ctx:       ctx,
		listView:  listViewModel,
		statusbar: statusbarModel,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.listView.Init(), common.SessionStatsCmdInstant(m.ctx))
}

// BUG: TorrentListView doesn't update once we are out of TorrentDetailView
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.ctx.Width = msg.Width
		m.ctx.Height = msg.Height
		h, v := appStyle.GetFrameSize()
		width, height := msg.Width-h, msg.Height-v

		// TODO: convert this to a method and make `List` private
		m.listView.List.SetSize(width, height-statusbar.Height)
		m.statusbar.SetSize(width)
		/*
		   NOTE: take padding of NormalTitle style if we modify it
		   https://github.com/notjedi/gotem/blob/9471ba90d28728b1dccc96cdea0a8db20c53b6de/internal/ui/components/listview/listview.go#L25
		   HACK: directly using `percent` * msg.Width is a little buggy, cause sometimes in
		   descSpacing `0.25 * 3` < `0.75`, which leads to off by 1 or 2 spacing issues
		*/
		textWidth := uint(math.Ceil(float64(msg.Width-h) * 0.05)) // 5%
		m.ctx.SetTitleSpacing([...]uint{
			12 * textWidth, // 60%
			4 * textWidth,  // 20%
			4 * textWidth,  // 20%
		})
		m.ctx.SetDescSpacing([...]uint{
			4 * textWidth, // 20%
			4 * textWidth, // 20%
			4 * textWidth, // 20%
			4 * textWidth, // 20%
			4 * textWidth, // 20%
		})

		if m.currView == TorrentDetailView {
			msg.Width = width
			msg.Height = height - tabs.TabHeight
			m.detailView, cmd = m.detailView.Update(msg)
		}
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		// TODO: replace if statements with switch statements
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		} else if msg.Type == tea.KeyEsc || msg.String() == "q" {
			if m.currView == TorrentDetailView {
				cmds = append(cmds, common.AllTorrentInfoMsgInstant(m.ctx))
				m.currView = TorrentListView
			} else {
				return m, tea.Quit
			}
		} else if msg.Type == tea.KeyRight || msg.String() == "l" {
			if m.currView == TorrentListView {
				// TODO: handle index when items are filtered
				// https://stackoverflow.com/questions/43883502/how-to-invoke-a-method-with-pointer-receiver-after-type-assertion
				torrent := m.listView.List.Items()[m.listView.List.Index()].(common.TorrentItem).Item()
				h, v := appStyle.GetFrameSize()
				width := m.ctx.Width - h
				height := m.ctx.Height - v - tabs.TabHeight

				m.detailView = detailview.New(*torrent.HashString, *torrent.ID, width, height, m.ctx)
				cmds = append(cmds, m.detailView.Init())

				m.currView = TorrentDetailView
				return m, tea.Batch(cmds...)
			}
		}

	// TODO: implement statusbarUpdateMsg
	// creating separate msg for this cause, doing all compute in a go routine
	// case statusbarUpdateMsg:
	case common.AllTorrentInfoMsg:
		torrentItems := convertListItemToTorrentItem(msg)
		m.statusbar.SetContent(getStatusBarContent(torrentItems))

	case common.SessionStatsMsg:
		sessionStats = msg
		cmds = append(cmds, common.SessionStatsCmd(m.ctx))
	}

	// TODO: should i move this before the switch statement, so i don't need to check for unwanted
	// messages
	if m.currView == TorrentListView {
		m.listView, cmd = m.listView.Update(msg)
	} else if m.currView == TorrentDetailView {
		m.detailView, cmd = m.detailView.Update(msg)
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.currView == TorrentListView {
		return lipgloss.JoinVertical(lipgloss.Top,
			appStyle.Render(fmt.Sprintf("%s\n%s", m.listView.View(), m.statusbar.View())),
		)
	} else if m.currView == TorrentDetailView {
		return lipgloss.JoinVertical(lipgloss.Top,
			appStyle.Render(fmt.Sprintf("%s\n%s", m.detailView.View(), m.statusbar.View())),
		)
	}
	return ""
}

func getSessionStatsString(stats common.SessionStatsMsg) string {
	return fmt.Sprintf(" %s  %s", humanize.Bytes(uint64(stats.CumulativeStats.DownloadedBytes)), humanize.Bytes(uint64(stats.CumulativeStats.UploadedBytes)))
}

func convertListItemToTorrentItem(items []list.Item) []transmissionrpc.Torrent {
	torrentItems := make([]transmissionrpc.Torrent, len(items))
	for i, item := range items {
		torrentItems[i] = item.(common.TorrentItem).Item()
	}
	return torrentItems
}

func getStatusBarContent(torrents []transmissionrpc.Torrent) (string, string, string, string) {
	/*
		1. total torrent info -            - do in a go routine, every 2 seconds if it's a large list and every 1 second if it's small
		2. filter by and sort by values - 料  惡  僚 寮           -- don't really know how to go about this as of now
		3. total gb uploaded? time elapsed? file count? -              神 羽 ﮫ ﲊ            -- do this on list item change
	*/
	netDownloadSpeed := getTotalDownloadSpeed(&torrents)
	netUploadSpeed := getTotalUploadSpeed(&torrents)
	statusString := getStatusString(&torrents)
	return config.ProgramName, statusString, fmt.Sprintf(" %s  %s", humanize.Bytes(uint64(netDownloadSpeed)), humanize.Bytes(uint64(netUploadSpeed))), getSessionStatsString(sessionStats)
}

func getStatusString(torrents *[]transmissionrpc.Torrent) string {
	downloading, seeding, paused := 0, 0, 0
	for _, torrent := range *torrents {
		switch *torrent.Status {
		case transmissionrpc.TorrentStatusDownload:
			downloading += 1
		case transmissionrpc.TorrentStatusSeed:
			seeding += 1
		case transmissionrpc.TorrentStatusStopped,
			transmissionrpc.TorrentStatusCheckWait,
			transmissionrpc.TorrentStatusCheck,
			transmissionrpc.TorrentStatusDownloadWait,
			transmissionrpc.TorrentStatusSeedWait,
			transmissionrpc.TorrentStatusIsolated:
			paused += 1
		default:
			paused += 1
		}
	}
	return fmt.Sprintf("Downloading: %d Seeding: %d Paused: %d", downloading, seeding, paused)
}

func getTotalDownloadSpeed(torrents *[]transmissionrpc.Torrent) int64 {
	var netDownloadSpeed int64 = 0
	for _, torrent := range *torrents {
		netDownloadSpeed += *torrent.RateDownload
	}
	return netDownloadSpeed
}

func getTotalUploadSpeed(torrents *[]transmissionrpc.Torrent) int64 {
	var netUploadSpeed int64 = 0
	for _, torrent := range *torrents {
		netUploadSpeed += *torrent.RateUpload
	}
	return netUploadSpeed
}
