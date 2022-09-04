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
		tabs tabs.Model
	}
)

const (
	OverviewTab Tab = iota + 1
	FilesTab
	ChunksTab
	// PeersTab
	// TrackersTab
)

func New(hash string, id int64, ctx *context.ProgramContext) Model {
	overviewTab := overviewtab.New(hash, id)
	var models []tea.Model = []tea.Model{
		overviewTab,
	}

	tabsModel := tabs.New(len(models))
	tabsModel.SetTabModels(models)
	tabsModel.SetTabTitles([]string{"Overview"})
	tabsModel.SetCurrentTab(0)

	return Model{
		hash: hash,
		id:   id,
		ctx:  ctx,
		tabs: tabsModel,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.tabs, cmd = m.tabs.Update(msg)
	cmds = append(cmds, cmd)

	// TODO: append TorrentInfoCmd to cmds
	// TODO: add this Cmd on page change

	cmds = append(cmds, func() tea.Msg {
		return common.GenerateTorrentInfoMsg(m.ctx, m.id)
	})
	// cmds = append(cmds, common.GenerateTorrentInfoMsg(m.ctx, m.id))

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.tabs.View()
}
