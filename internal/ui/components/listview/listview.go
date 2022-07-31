package listview

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	// "github.com/hekmon/transmissionrpc/v2"
	"github.com/notjedi/gotem/internal/context"
)

var (
	torrentFields = []string{"id", "name", "downloadDir", "status", "desiredAvailable",
		"rateDownload", "rateUpload", "eta", "uploadRatio", "sizeWhenDone", "haveValid",
		"haveUnchecked", "addedDate", "uploadedEver", "errorString", "recheckProgress",
		"peersConnected", "uploadLimit", "downloadLimit", "uploadLimited", "downloadLimited",
		"bandwidthPriority", "peersSendingToUs", "peersGettingFromUs", "seedRatioLimit",
		"seedRatioMode", "magnetLink", "honorsSessionLimits", "metadataPercentComplete"}
	docStyle = lipgloss.NewStyle().Margin(1, 2)
)

type torrentUpdateMsg []list.Item

type Model struct {
	list list.Model
	ctx  *context.Context
}

func New(ctx *context.Context) Model {
	listDelegate := NewCustomDelegate()
	listModel := list.New([]list.Item{}, listDelegate, 0, 0)
	listModel.Styles.Title = listModel.Styles.Title.Copy().
		Bold(true).
		Italic(true)
	listModel.DisableQuitKeybindings()
	return Model{
		list: listModel,
		ctx:  ctx,
	}
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		return generateTorrentUpdateMsg(m)
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case torrentUpdateMsg:
		m.list.SetItems(msg)
		// TODO: should i really be using m.updateTorrentsCmd?
		return m, m.updateTorrentsCmd()
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return docStyle.Render(m.list.View())
}
