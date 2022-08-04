package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/statusbar"
	"github.com/notjedi/gotem/internal/context"
	"github.com/notjedi/gotem/internal/theme"
	"github.com/notjedi/gotem/internal/ui/components/detailview"
	"github.com/notjedi/gotem/internal/ui/components/listview"
)

type View int32
type Model struct {
	currView   View
	context    context.Context
	listView   listview.Model
	detailView detailview.Model
	statusbar  statusbar.Bubble
}

type Direction int

const (
	TorrentListView View = iota + 1
	TorrentDetailView
	HelpView
)

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2, 1, 2)
)

func New(ctx context.Context) Model {
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

	listViewModel := listview.New(ctx, theme)

	return Model{
		currView:  TorrentListView,
		context:   ctx,
		listView:  listViewModel,
		statusbar: statusbarModel,
	}
}

func (m Model) Init() tea.Cmd {
	return m.listView.Init()
}

// TODO: update statusbar content
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.listView, cmd = m.listView.Update(msg)
	cmds = append(cmds, cmd)
	// m.detailView, cmd = m.detailView.Update(msg)
	// cmds = append(cmds, cmd)

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		// m.context.ListWidth = float32(msg.Width - h - m.listView.TitlePadding)
		m.listView.List.SetSize(msg.Width-h, msg.Height-statusbar.Height-v)
		m.statusbar.SetSize(msg.Width - h)
		m.statusbar.SetContent(
			"NORMAL",
			"~/.config/nvim/lua/options.lua",
			"lua",
			fmt.Sprintf("%dx%d", msg.Width, statusbar.Height),
		)

		// TODO: find a better way to align fields
		textWidth := float32(msg.Width - h - m.listView.TitlePadding)
		if textWidth <= 140 {
			m.context.SetTitleSpacing([...]uint{uint(0.50 * textWidth), uint(0.25 * textWidth),
				uint(0.25 * textWidth)})
			m.context.SetDescSpacing([...]uint{uint(0.25 * textWidth), uint(0.25 * textWidth),
				uint(0.25 * textWidth), uint(0.25 * textWidth)})
		} else {
			m.context.SetTitleSpacing([...]uint{uint(0.75 * textWidth), uint(0.15 * textWidth),
				uint(0.10 * textWidth)})
			m.context.SetDescSpacing([...]uint{uint(0.25 * textWidth), uint(0.25 * textWidth),
				uint(0.25 * textWidth), uint(0.25 * textWidth)})
		}

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc || msg.String() == "q" {
			return m, tea.Quit
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.currView == TorrentListView {
		return lipgloss.JoinVertical(lipgloss.Top,
			docStyle.Render(fmt.Sprintf("%s\n%s", m.listView.View(), m.statusbar.View())),
		)
	} else if m.currView == TorrentDetailView {
		// TODO
		return ""
	}
	return ""
}
