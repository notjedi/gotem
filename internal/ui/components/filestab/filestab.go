package filestab

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/hekmon/transmissionrpc/v2"
	"github.com/notjedi/gotem/internal/ui/common"
	"github.com/notjedi/gotem/internal/ui/utils"
)

const (
	tabSpacing = "    "
	entry      = "  "
	lastEntry  = "  "
	// entry      = "┣ "
	// lastEntry  = "┗ "

	columnKeyNumber   = "fileNumebr"
	columnKeyProgress = "progress"
	columnKeySize     = "size"
	columnKeyFilename = "fileName"

	columnNumberName   = "#"
	columnProgressName = "Progress"
	columnSizeName     = "Size"
	columnFilenameName = "Filename"

	paddingRight = " "
	paddingLeft  = " "
)

var customBorder = table.Border{
	Top:    "─",
	Left:   "│",
	Right:  "│",
	Bottom: "─",

	TopRight:    "╮",
	TopLeft:     "╭",
	BottomRight: "╯",
	BottomLeft:  "╰",

	TopJunction:    "╥",
	LeftJunction:   "├",
	RightJunction:  "┤",
	BottomJunction: "╨",
	InnerJunction:  "╫",

	InnerDivider: "║",
}

type Model struct {
	hash        string
	id          int64
	fileTree    *Directory
	table       table.Model
	torrentInfo transmissionrpc.Torrent
}

// TODO: remove width and height?
func New(hash string, id int64, width int, height int) tea.Model {
	table := table.New([]table.Column{
		table.NewColumn(columnKeyNumber, columnNumberName, 5).WithStyle(lipgloss.NewStyle().
			Align(lipgloss.Center)),
		table.NewColumn(columnKeyProgress, columnProgressName, 10).WithStyle(lipgloss.NewStyle().
			Align(lipgloss.Left)),
		table.NewColumn(columnKeySize, columnSizeName, 10).WithStyle(lipgloss.NewStyle().
			Align(lipgloss.Left)),
		table.NewFlexColumn(columnKeyFilename, columnFilenameName, 1).WithStyle(lipgloss.NewStyle().
			Align(lipgloss.Left)),
	}).WithTargetWidth(width).
		WithPaginationWrapping(true).
		WithPageSize(height - 8).
		HeaderStyle(lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).
			Foreground(lipgloss.Color("#A550DF"))).
		Focused(true)
		// Border(customBorder).

		// HeaderStyle(lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).
		// 	Foreground(lipgloss.Color("#A550DF"))).

		// HeaderStyle(lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).
		// 	Background(lipgloss.Color("#6124DF")).Foreground(lipgloss.Color("#ffffff")).
		// 	ColorWhitespace(false)).

	return Model{
		hash:  hash,
		id:    id,
		table: table,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	// var cmds []tea.Cmd

	switch msg := msg.(type) {
	case common.TorrentInfoMsg:
		torrentInfo := transmissionrpc.Torrent(msg)
		if *torrentInfo.HashString == m.hash {
			fileNumber := 1
			m.torrentInfo = torrentInfo
			m.fileTree = buildFileTree(torrentInfo.Files, torrentInfo.FileStats)
			rows := buildTable(m.fileTree, 0, &fileNumber)
			// rows = append([]table.Row{table.NewRow(table.RowData{})}, rows...)
			m.table = m.table.WithRows(rows)
		}
	}
	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	if m.fileTree == nil {
		return ""
	}

	return fmt.Sprintf("\n%v", m.table.View())
}

func buildTable(fileTree *Directory, depth int, fileNumber *int) []table.Row {
	rows := []table.Row{}

	// name := fileTree.name
	// if depth != 0 {
	// name = fmt.Sprintf("%v%v%v", strings.Repeat(tabSpacing, depth), entry, name)
	// }
	name := fmt.Sprintf(" %v%v", strings.Repeat(tabSpacing, depth), fileTree.name)
	// if depth == 0 && len(fileTree.children) == 0 && len(fileTree.files) == 0 {
	// }

	row := table.NewRow(table.RowData{
		columnKeyNumber:   "",
		columnKeyProgress: "",
		columnKeySize:     "",
		columnKeyFilename: name,
	})
	rows = append(rows, row)

	for _, directory := range fileTree.children {
		childRows := buildTable(directory, depth+1, fileNumber)
		rows = append(rows, childRows...)
	}

	for _, file := range fileTree.files {
		fileName := fmt.Sprintf(" %v%v", strings.Repeat(tabSpacing, depth+1), file.name)
		// fileName := strings.Repeat(tabSpacing, depth+1)
		// if idx == len(fileTree.files)-1 {
		// 	fileName = fmt.Sprintf("%v%v%v", fileName, lastEntry, file.name)
		// } else {
		// 	fileName = fmt.Sprintf("%v%v%v", fileName, entry, file.name)
		// }

		row := table.NewRow(table.RowData{
			columnKeyNumber:   fmt.Sprintf("%v", *fileNumber),
			columnKeyProgress: fmt.Sprintf(" %0.2f", file.percentDone) + "%",
			columnKeySize:     fmt.Sprintf(" %v", utils.HumanizeBytes(file.bytesTotal)),
			columnKeyFilename: fileName,
		})
		rows = append(rows, row)
		*fileNumber = *fileNumber + 1
		// *fileNumber += 1     // TODO: will this work?
	}

	return rows
}
