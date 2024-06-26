package detailview

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/notjedi/gotem/internal/context"
	"github.com/notjedi/gotem/internal/ui/common"
	"github.com/notjedi/gotem/internal/ui/components/chunkstab"
	"github.com/notjedi/gotem/internal/ui/components/filestab"
	"github.com/notjedi/gotem/internal/ui/components/overviewtab"
	"github.com/notjedi/tabs"
)

type (
	Tab   int
	Model struct {
		ctx  *context.ProgramContext
		Tabs tabs.Model
		hash string
		id   int64
	}
)

var ignore = false

func New(hash string, id int64, width int, height int, ctx *context.ProgramContext) Model {
	overviewTab := overviewtab.New(hash, id, width, height)
	filesTab := filestab.New(hash, id, width, height)
	chunksTab := chunkstab.New(hash, id, width, height)

	models := []tea.Model{
		overviewTab,
		filesTab,
		chunksTab,
	}

	tabsModel := tabs.New(len(models))
	tabsModel.SetTabModels(models)
	tabsModel.SetTabTitles([]string{"Overview", "Files", "Chunks"})
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
	return common.TorrentInfoCmdInstant(m.ctx, m.id)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case common.TorrentInfoMsg:
		if ignore {
			ignore = false
		} else {
			cmds = append(cmds, common.TorrentInfoCmd(m.ctx, m.id))
		}

	// NOTE: generate new TorrentInfoMsg instantly on page change, update keys if keymap is updated
	// TODO: should i call init on page change in tabs library?
	// TODO: move this to each model's init once the tabs library supports it
	case tea.KeyMsg:
		if msg.Type == tea.KeyRight || msg.String() == "l" ||
			msg.Type == tea.KeyLeft || msg.String() == "h" {
			ignore = true
			cmds = append(cmds, func() tea.Msg {
				return common.GenerateTorrentInfoMsg(m.ctx, m.id)
			})
		}
	}

	m.Tabs, cmd = m.Tabs.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.Tabs.View()
}
