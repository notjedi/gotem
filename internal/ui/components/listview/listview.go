package listview

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/notjedi/gotem/internal/config"
	"github.com/notjedi/gotem/internal/context"
	"github.com/notjedi/gotem/internal/ui/common"
)

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
		return common.GenerateAllTorrentInfoMsg(m.ctx)
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.List, cmd = m.List.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case common.AllTorrentInfoMsg:
		// BUG: the items are disappearing if i update the items while filterState != Unfiltered
		// TODO: only send next request if filterState != Unfiltered
		// TODO: try removing this if check now
		if m.List.FilterState() == list.Unfiltered {
			// update items only if `filterState` == Unfiltered
			cmd = m.List.SetItems(msg)
			cmds = append(cmds, cmd)
		}
		cmds = append(cmds, common.AllTorrentInfoCmd(m.ctx))
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.List.View()
}
