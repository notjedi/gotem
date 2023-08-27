package chunkstab

import (
	"encoding/base64"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hekmon/transmissionrpc/v2"
	"github.com/muesli/reflow/wordwrap"
	"github.com/notjedi/gotem/internal/ui/common"
	"github.com/notjedi/gotem/internal/ui/utils"
)

const HeightPadding = 2

type Model struct {
	hash       string
	id         int64
	width      int
	height     int
	yOffset    int
	pieceCount int
	pieces     []byte
	lines      []string
}

var (
	tabStyle   = lipgloss.NewStyle().Margin(1, 2, 1, 2)
	fillString = lipgloss.NewStyle().SetString(" ").Background(lipgloss.Color("#D9DCCF")).String()
)

func New(hash string, id int64, width int, height int) Model {
	h, v := tabStyle.GetFrameSize()
	return Model{
		hash:       hash,
		id:         id,
		yOffset:    0,
		width:      width - h,
		height:     height - v - HeightPadding,
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
				chunksString := m.getChunksString()
				m.lines = strings.Split(wordwrap.String(chunksString, m.width), "\n")
			}
		}

	case tea.KeyMsg:
		if msg.String() == "j" {
			m.yOffset = utils.Min(m.yOffset+1, utils.Abs(len(m.lines)-m.height))
		} else if msg.String() == "k" {
			m.yOffset = utils.Max(m.yOffset-1, 0)
		}

	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return tabStyle.Render(strings.Join(m.visibleLines(), "\n"))
}

func (m *Model) SetSize(width, height int) {
	h, v := tabStyle.GetFrameSize()
	m.width = width - h
	m.height = height - v - HeightPadding
	m.yOffset = utils.Min(m.yOffset, utils.Abs(len(m.lines)-m.height))

	chunksString := m.getChunksString()
	m.lines = strings.Split(wordwrap.String(chunksString, m.width), "\n")
}

func (m *Model) visibleLines() []string {
	if len(m.lines) > m.height {
		top := utils.Min(m.yOffset, len(m.lines)-m.height)
		bottom := utils.Clamp(m.yOffset+m.height, top, len(m.lines))
		return m.lines[top:bottom]
		// return []string{m.lines[len(m.lines)-1]}
		// return []string{m.lines[0]}
	}
	return m.lines
}

func (m *Model) getChunksString() string {
	b := strings.Builder{}
	// b.WriteString("\n\n")

	for i := 0; i < m.pieceCount; i++ {
		idx, shift := utils.DivMod(i, 8)
		havePiece := ((m.pieces[idx] >> (7 - shift)) & 1) == 1
		if havePiece {
			b.WriteString(fillString)
		} else {
			b.WriteString("-")
		}

		if (i+1)%m.width == 0 {
			b.WriteString("\n")
		}
	}
	return b.String()
}
