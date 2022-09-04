package overviewtab

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hekmon/transmissionrpc/v2"
	"github.com/notjedi/gotem/internal/ui/common"
)

type Model struct {
	hash            string
	id              int64
	prevTorrentInfo transmissionrpc.Torrent
}

func New(hash string, id int64) tea.Model {
	return Model{
		hash: hash,
		id:   id,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.TorrentInfoMsg:
		// TODO: should i check if the hash is same?
		m.prevTorrentInfo = transmissionrpc.Torrent(msg)
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	return ""
}

func (m *Model) SetPrevTorrentInfo(torrentInfo transmissionrpc.Torrent) {
	m.prevTorrentInfo = torrentInfo
}
