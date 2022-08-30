package ui

import (
	"fmt"
	"math"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hekmon/transmissionrpc/v2"
	"github.com/knipferrc/teacup/statusbar"
	"github.com/notjedi/gotem/internal/config"
	"github.com/notjedi/gotem/internal/context"
	"github.com/notjedi/gotem/internal/theme"
	"github.com/notjedi/gotem/internal/ui/components/detailview"
	"github.com/notjedi/gotem/internal/ui/components/listview"
)

type View int32
type Model struct {
	currView   View
	ctx        *context.ProgramContext
	listView   listview.Model
	detailView detailview.Model
	statusbar  statusbar.Bubble
}

const (
	TorrentListView View = iota + 1
	TorrentDetailView
	HelpView
)

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2, 1, 2)
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
	return m.listView.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.currView == TorrentListView {
		m.listView, cmd = m.listView.Update(msg)
	} else if m.currView == TorrentDetailView {
		m.detailView, cmd = m.detailView.Update(msg)
	}
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.listView.List.SetSize(msg.Width-h, msg.Height-statusbar.Height-v)
		m.statusbar.SetSize(msg.Width - h)

		/*
		   NOTE: take padding of NormalTitle style if we modify it
		   https://github.com/notjedi/gotem/blob/9471ba90d28728b1dccc96cdea0a8db20c53b6de/internal/ui/components/listview/listview.go#L25
		   HACK: directly using `percent` * msg.Width is a little buggy, cause sometimes in
		   descSpacing `0.25 * 3` < `0.75`, which leads to off by 1 or 2 spacing issues
		*/
		textWidth := uint(math.Ceil(float64(msg.Width-h) * 0.05)) // 5%
		m.ctx.SetTitleSpacing([...]uint{12 * textWidth,           // 60%
			4 * textWidth, // 20%
			4 * textWidth, // 20%
		})
		m.ctx.SetDescSpacing([...]uint{4 * textWidth, // 20%
			4 * textWidth, // 20%
			4 * textWidth, // 20%
			4 * textWidth, // 20%
			4 * textWidth, // 20%
		})

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc || msg.String() == "q" {
			return m, tea.Quit
		} else if msg.Type == tea.KeyRight || msg.String() == "l" {
			// TODO: make current view a field of global context
			// update view in listview, on the item selected
			// continue if no item is selected
			m.currView = TorrentDetailView
		}

		// TODO: implement statusbarUpdateMsg
		// creating separate msg for this cause, doing all compute in a go routine
		// case statusbarUpdateMsg:
		// 	m.statusbar.SetContent(getStatusBarContent(msg))
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.currView == TorrentListView {
		return lipgloss.JoinVertical(lipgloss.Top,
			docStyle.Render(fmt.Sprintf("%s\n%s", m.listView.View(), m.statusbar.View())),
		)
	} else if m.currView == TorrentDetailView {
		return lipgloss.JoinVertical(lipgloss.Top,
			docStyle.Render(fmt.Sprintf("%s\n%s", m.detailView.View(), m.statusbar.View())),
		)
	}
	return ""
}

func getStatusBarContent(torrents []transmissionrpc.Torrent) (string, string, string, string) {
	/*
		1. total torrent info -            - do in a go routine, every 2 seconds if it's a large list and every 1 second if it's small
		2. filter by and sort by values - 料  惡  僚 寮           -- don't really know how to go about this as of now
		3. total gb uploaded? time elapsed? file count? -              神 羽 ﮫ ﲊ            -- do this on list item change
		4. net download and upload speed -           --- add all the speeds and update on torrentInfoUpdateMsg
	*/
	return config.ProgramName, "", "", ""
}
