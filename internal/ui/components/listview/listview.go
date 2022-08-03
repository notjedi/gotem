package listview

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/notjedi/gotem/internal/context"
	"github.com/notjedi/gotem/internal/theme"
)

const (
	// TODO: move this to context? as this is a global thing
	programName string = "gotem"
)

/*
	TODO: make this part of context? i'll be using more or less the same fields in detailview

i can prolly use package specific fields to squeeze a tiny tiny amount of performance. since the
response to the request made comes from c, i can kinda assume that adding more fields is basically
free. so the performance gain comes down to json serialization and deserialization? and as go is
also kinda fast, ig the performance gain here is immeasurable? anyways, i'll use package specific
fields for now and make it part of context as we go?
*/
var (
	torrentFields = []string{"id", "hashString", "name", "status", "rateDownload", "rateUpload",
		"eta", "uploadRatio", "sizeWhenDone", "haveValid", "uploadedEver", "recheckProgress",
		"peersConnected", "uploadLimited", "downloadLimited", "bandwidthPriority",
		"peersSendingToUs", "peersGettingFromUs", "seedRatioLimit", "trackerStats", "magnetLink",
		"honorsSessionLimits", "metadataPercentComplete", "percentDone"}
)

type torrentUpdateMsg []list.Item
type Model struct {
	List         list.Model
	ctx          *context.Context
	TitlePadding int
}

func New(ctx *context.Context, theme theme.Theme) Model {
	listDelegate := NewCustomDelegate()
	titlePadding := listDelegate.Styles.NormalTitle.GetPaddingLeft() +
		listDelegate.Styles.NormalTitle.GetPaddingRight()

	listModel := list.New([]list.Item{}, listDelegate, 0, 0)
	listModel.SetShowHelp(false)
	listModel.SetShowStatusBar(false)
	listModel.DisableQuitKeybindings()
	listModel.Title = programName
	listModel.Styles.Title = listModel.Styles.Title.Copy().
		Bold(true).
		Italic(true).
		Background(theme.TitleBackgroundColor).
		Foreground(theme.TitleForegroundColor)
	return Model{
		List:         listModel,
		ctx:          ctx,
		TitlePadding: titlePadding,
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
		// BUG: the items are disappearing if i update the items while filterState != Unfiltered
		// FIXME: ig it makes sense to stop updating the items while filtering
		if m.List.FilterState() == list.Unfiltered {
			m.List.SetItems(msg)
		}
		return m, m.updateTorrentsCmd()
	}

	return m, cmd
}

func (m Model) View() string {
	return m.List.View()
}
