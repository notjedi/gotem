package detailview

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/notjedi/gotem/internal/context"
	"github.com/notjedi/gotem/internal/ui/components/tabs"
)

type Tab int
type torrentInfoMsg []list.Item
type Model struct {
	hash string
	id   int64
	ctx  context.Context
	tabs tabs.Model
}

const (
	OverviewTab Tab = iota + 1
	FilesTab
	ChunksTab
	// PeersTab
	// TrackersTab
)

func New(hash string, id int64, ctx context.Context) Model {
	tabsModel := tabs.New(3)
	return Model{
		hash: hash,
		id:   id,
		ctx:  ctx,
		tabs: tabsModel,
	}
}

func (m Model) Init() tea.Cmd {
	// since the default starting section is listview, we can return nil here
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {

	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.tabs, cmd = m.tabs.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return ""
}
