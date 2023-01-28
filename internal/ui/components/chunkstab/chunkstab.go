package chunkstab

import (
	"encoding/base64"
	// "fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hekmon/transmissionrpc/v2"
	"github.com/notjedi/gotem/internal/ui/common"
)

type Model struct {
	hash       string
	id         int64
	width      int
	height     int
	pieces     []byte
	pieceCount int64
}

func New(hash string, id int64, width int, height int) Model {
	return Model{
		hash:       hash,
		id:         id,
		width:      width,
		height:     height,
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
				m.pieceCount = *torrentInfo.PieceCount
				// fmt.Println(pieces, *torrentInfo.PieceCount, len(pieces))
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.pieceCount == 0 {
		return ""
	}

	b := strings.Builder{}
	b.WriteString("\n")
	b.WriteString("\n")

	for i := 0; i < len(m.pieces); i++ {
		temp := 1
		for shift := 7; shift >= 0; shift-- {
			if (m.pieces[i] & (1 << shift)) != 0 {
				b.WriteString("1")
			} else {
				b.WriteString("-")
			}
			if i != 0 && (((i*8)+temp)%m.width) == 0 {
				b.WriteString("\n")
			}
			temp += 1
		}
	}

	// for y := 0; y < int(m.pieceCount)/m.width; y++ {
	// 	for x := 0; x < m.width; x++ {
	// 	}
	// 	b.WriteString("\n")
	// }

	return b.String()
}
