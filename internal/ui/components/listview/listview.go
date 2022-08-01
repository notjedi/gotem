package listview

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/notjedi/gotem/internal/context"
	"github.com/notjedi/gotem/internal/theme"
)

var (
	torrentFields = []string{"id", "name", "downloadDir", "status", "desiredAvailable",
		"rateDownload", "rateUpload", "eta", "uploadRatio", "sizeWhenDone", "haveValid",
		"haveUnchecked", "addedDate", "uploadedEver", "errorString", "recheckProgress",
		"peersConnected", "uploadLimit", "downloadLimit", "uploadLimited", "downloadLimited",
		"bandwidthPriority", "peersSendingToUs", "peersGettingFromUs", "seedRatioLimit",
		"seedRatioMode", "magnetLink", "honorsSessionLimits", "metadataPercentComplete"}
)

type torrentUpdateMsg []list.Item

type Model struct {
	List list.Model
	ctx  *context.Context
}

func New(ctx *context.Context, theme theme.Theme) Model {
	listDelegate := NewCustomDelegate()
	listModel := list.New([]list.Item{}, listDelegate, 0, 0)
	// TODO: make constant
	listModel.Title = "gotem"
	listModel.SetShowHelp(false)
	listModel.SetShowStatusBar(false)
	listModel.DisableQuitKeybindings()
	listModel.Styles.Title = listModel.Styles.Title.Copy().
		Bold(true).
		Italic(true).
		Background(theme.TitleBackgroundColor).
		Foreground(theme.TitleForegroundColor)
	return Model{
		List: listModel,
		ctx:  ctx,
	}
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		return generateTorrentUpdateMsg(m)
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)

	switch msg := msg.(type) {
	case torrentUpdateMsg:
		m.List.SetItems(msg)
		return m, m.updateTorrentsCmd()
	}

	return m, cmd
}

func (m Model) View() string {
	return m.List.View()
}
