package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hekmon/transmissionrpc/v2"
	"github.com/knipferrc/teacup/statusbar"
	"github.com/notjedi/gotem/internal/theme"
)

type Model struct {
	client    *transmissionrpc.Client
	statusbar statusbar.Bubble
}

func New(client *transmissionrpc.Client) Model {
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
	return Model{
		statusbar: statusbarModel,
		client:    client,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	return m, nil
}

func (m Model) View() string {
	return m.statusbar.View()
}
