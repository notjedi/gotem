package ui

import (
	"fmt"

	// "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/statusbar"
	"github.com/notjedi/gotem/internal/context"
	"github.com/notjedi/gotem/internal/theme"
	"github.com/notjedi/gotem/internal/ui/components/listview"
)

type Model struct {
	currView  int
	context   *context.Context
	listview  listview.Model
	statusbar statusbar.Bubble
}

func New(ctx *context.Context) Model {
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

	listviewModel := listview.New(ctx)

	return Model{
		currView:  1,
		context:   ctx,
		listview:  listviewModel,
		statusbar: statusbarModel,
	}
}

func (m Model) Init() tea.Cmd {
	return m.listview.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.statusbar.SetSize(msg.Width)
		m.statusbar.SetContent(
			"NORMAL",
			"~/.config/nvim/lua/options.lua",
			"lua",
			fmt.Sprintf("%dx%d", msg.Width, msg.Height),
		)

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc {
			return m, tea.Quit
		}
	}
	m.listview, cmd = m.listview.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Top,
		m.listview.View(),
		m.statusbar.View(),
	)
}
