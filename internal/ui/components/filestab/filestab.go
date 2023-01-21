package filestab

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hekmon/transmissionrpc/v2"
	"github.com/notjedi/gotem/internal/ui/common"
)

type Model struct {
	hash        string
	id          int64
	fileTree    *Directory
	torrentInfo transmissionrpc.Torrent
}

// TODO: remove width and height?
func New(hash string, id int64, width int, height int) tea.Model {
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
		torrentInfo := transmissionrpc.Torrent(msg)
		if *torrentInfo.HashString == m.hash {
			m.torrentInfo = torrentInfo
			m.fileTree = buildFileTree(torrentInfo.Files, torrentInfo.FileStats)
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.fileTree == nil {
		return ""
	}

	return drawFileTree(m.fileTree, 0)
}

func drawFileTree(fileTree *Directory, depth int) string {
	var builder strings.Builder

	if depth != 0 {
		builder.WriteString(strings.Repeat("  ", depth))
		builder.WriteString("┣ ")
	}
	builder.WriteString(fileTree.name)
	builder.WriteString("\n")

	for _, directory := range fileTree.children {
		builder.WriteString(drawFileTree(directory, depth+1))
	}

	for idx, file := range fileTree.files {
		builder.WriteString(strings.Repeat("  ", depth+1))
		if idx == len(fileTree.files)-1 {
			builder.WriteString("┗ ")
		} else {
			builder.WriteString("┣ ")
		}
		builder.WriteString(file.name)
		builder.WriteString("\n")
	}
	return builder.String()
}
