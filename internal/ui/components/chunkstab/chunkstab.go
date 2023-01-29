package chunkstab

import (
	"encoding/base64"
	// "fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hekmon/transmissionrpc/v2"
	"github.com/notjedi/gotem/internal/ui/common"
	"github.com/notjedi/gotem/internal/ui/utils"
)

type Model struct {
	hash       string
	id         int64
	width      int
	height     int
	pieces     []byte
	pieceCount int
}

var tabStyle = lipgloss.NewStyle().Margin(1, 2, 1, 2)

func New(hash string, id int64, width int, height int) Model {
	h, v := tabStyle.GetFrameSize()
	return Model{
		hash:       hash,
		id:         id,
		width:      width - h,
		height:     height - v,
		pieceCount: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case common.TorrentInfoMsg:
		torrentInfo := transmissionrpc.Torrent(msg)
		if *torrentInfo.HashString == m.hash {
			pieces, err := base64.StdEncoding.DecodeString(*torrentInfo.Pieces)
			if err == nil {
				m.pieces = pieces
				m.pieceCount = int(*torrentInfo.PieceCount)
			}
		}

	case tea.WindowSizeMsg:
		h, v := tabStyle.GetFrameSize()
		m.width = msg.Width - h
		m.height = msg.Height - v
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.pieceCount == 0 {
		return ""
	}

	b := strings.Builder{}
	b.WriteString("\n\n")

	fill := lipgloss.NewStyle().SetString(" ").Background(lipgloss.Color("#D9DCCF"))
	for i := 0; i < m.pieceCount; i++ {
		idx, shift := utils.DivMod(i, 8)
		havePiece := ((m.pieces[idx] >> (7 - shift)) & 1) == 1
		if havePiece {
			b.WriteString(fill.String())
		} else {
			b.WriteString("-")
		}

		if (i+1)%m.width == 0 {
			b.WriteString("\n")
		}
	}

	return tabStyle.Render(b.String())
}
