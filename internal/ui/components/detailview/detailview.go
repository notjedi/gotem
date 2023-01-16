package detailview

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/notjedi/gotem/internal/context"
	"github.com/notjedi/gotem/internal/ui/common"
	"github.com/notjedi/gotem/internal/ui/components/overviewtab"
	"github.com/notjedi/tabs"
)

type (
	Tab            int
	torrentInfoMsg []list.Item
	Model          struct {
		hash string
		id   int64
		ctx  *context.ProgramContext
		Tabs tabs.Model
	}
)

const (
	OverviewTab Tab = iota + 1
	FilesTab
	ChunksTab
	// PeersTab
	// TrackersTab
)

// TODO: do we need both hash and id?
// TODO: add width arg
func New(hash string, id int64, width int, height int, ctx *context.ProgramContext) Model {
	overviewTab := overviewtab.New(hash, id, width, height)
	var models []tea.Model = []tea.Model{
		overviewTab,
	}

	tabsModel := tabs.New(len(models))
	tabsModel.SetTabModels(models)
	tabsModel.SetTabTitles([]string{"Overview"})
	tabsModel.SetCurrentTab(0)
	tabsModel.SetSize(width, height)

	return Model{
		hash: hash,
		id:   id,
		ctx:  ctx,
		Tabs: tabsModel,
	}
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		return common.GenerateTorrentInfoMsg(m.ctx, m.id)
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg.(type) {
	case common.TorrentInfoMsg:
		cmds = append(cmds, common.TorrentInfoCmd(m.ctx, m.id))
	}

	m.Tabs, cmd = m.Tabs.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.Tabs.View()
}
