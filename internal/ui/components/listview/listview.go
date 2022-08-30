package listview

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/notjedi/gotem/internal/config"
	"github.com/notjedi/gotem/internal/context"
	"github.com/notjedi/gotem/internal/ui/common"
)

type torrentUpdateMsg []list.Item
type Model struct {
	List list.Model
	ctx  *context.ProgramContext
}

func New(ctx *context.ProgramContext) Model {
	listDelegate := list.NewDefaultDelegate()

	listModel := list.New([]list.Item{}, listDelegate, 0, 0)
	listModel.SetShowHelp(false)
	listModel.SetShowStatusBar(false)
	listModel.DisableQuitKeybindings()
	listModel.Title = config.ProgramName
	return Model{
		List: listModel,
		ctx:  ctx,
	}
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		return common.GenerateTorrentInfoMsg(m.ctx)
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)

	switch msg := msg.(type) {
	case common.TorrentInfoMsg:
		// BUG: the items are disappearing if i update the items while filterState != Unfiltered
		// TODO: only send next request if filterState != Unfiltered
		if m.List.FilterState() == list.Unfiltered {
			// update items only if `filterState` == Unfiltered
			m.List.SetItems(msg)
		}
		return m, common.TorrentInfoCmd(m.ctx)
	}

	return m, cmd
}

func (m Model) View() string {
	return m.List.View()
}
