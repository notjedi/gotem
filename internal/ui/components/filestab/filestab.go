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
	tabSpacing = "    " // 4 spaces

	columnKeyNumber   = "fileNumebr"
	columnKeyProgress = "progress"
	columnKeySize     = "size"
	columnKeyPriority = "priority"
	columnKeyFilename = "fileName"

	columnNumberName   = "#"
	columnProgressName = "Progress"
	columnSizeName     = "Size"
	columnPriorityName = "Priority"
	columnFilenameName = "Filename"
)

var (
	tabWidth  int
	tabHeight int

	priorityHighStyle = lipgloss.NewStyle().Background(lipgloss.Color("#f38ba8")).
				Foreground(lipgloss.Color("#1e1e2e"))
	priorityNormalStyle = lipgloss.NewStyle().Background(lipgloss.Color("#cdd6f4")). // should i use  #bac2de for background?
				Foreground(lipgloss.Color("#1e1e2e"))
	priorityLowStyle = lipgloss.NewStyle().Background(lipgloss.Color("#f9e2af")).
				Foreground(lipgloss.Color("#1e1e2e"))

	priorityMap = map[int64]string{
		1:  priorityHighStyle.Render("High"),
		0:  priorityNormalStyle.Render("Normal"),
		-1: priorityLowStyle.Render("Low"),
	}
)

type Model struct {
	hash        string
	id          int64
	fileTree    *Directory
	table       table.Model
	torrentInfo transmissionrpc.Torrent
}

func New(hash string, id int64, width int, height int) tea.Model {
	// headerStyle := lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).
	// 	Background(lipgloss.Color("#6124DF")).Foreground(lipgloss.Color("#ffffff")).
	// 	ColorWhitespace(false)
	headerStyle := lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).Foreground(lipgloss.Color("#A550DF"))
	tabWidth = width
	tabHeight = height

	// TODO: figure out a way to align columns to right with padding
	table := table.New([]table.Column{
		table.NewColumn(columnKeyNumber, columnNumberName, 5).WithStyle(lipgloss.NewStyle().
			Align(lipgloss.Center)),
		table.NewColumn(columnKeyProgress, columnProgressName, 10).WithStyle(lipgloss.NewStyle().
			Align(lipgloss.Left)),
		table.NewColumn(columnKeySize, columnSizeName, 10).WithStyle(lipgloss.NewStyle().
			Align(lipgloss.Left)),
		table.NewColumn(columnKeyPriority, columnPriorityName, 15).WithStyle(lipgloss.NewStyle().
			Align(lipgloss.Center)),
		table.NewFlexColumn(columnKeyFilename, columnFilenameName, 1).WithStyle(lipgloss.NewStyle().
			Align(lipgloss.Left)),
	}).WithTargetWidth(width).
		HeaderStyle(headerStyle).
		Focused(true)

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

	// TODO: handle KeyMsg so only files can be selected instead of both directories and files
	switch msg := msg.(type) {
	case common.TorrentInfoMsg:
		torrentInfo := transmissionrpc.Torrent(msg)
		if *torrentInfo.HashString == m.hash {
			fileNumber := 1
			m.torrentInfo = torrentInfo
			m.fileTree = buildFileTree(torrentInfo.Files, torrentInfo.FileStats)
			rows := buildFilesTable(m.fileTree, 0, &fileNumber)
			// rows = append([]table.Row{table.NewRow(table.RowData{})}, rows...)   // Append empty 1st row
			m.table = m.table.WithRows(rows)
			if m.table.TotalRows()-3 > tabHeight && m.table.PageSize() == 0 {
				m.table = m.table.WithPageSize(tabHeight - 8).WithPaginationWrapping(true)
			}
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

func buildFilesTable(fileTree *Directory, depth int, fileNumber *int) []table.Row {
	rows := []table.Row{}
	if depth != 0 { // Ignore the first directory as it's does not belong to the torrent and is there for implementation convenience
		name := fmt.Sprintf(" %v%v", strings.Repeat(tabSpacing, depth-1), fileTree.name)
		row := table.NewRow(table.RowData{
			columnKeyNumber:   "",
			columnKeyProgress: "",
			columnKeySize:     "",
			columnKeyPriority: "",
			columnKeyFilename: name,
		})
		rows = append(rows, row)
	}

	for _, directory := range fileTree.children {
		childRows := buildFilesTable(directory, depth+1, fileNumber)
		rows = append(rows, childRows...)
	}

	for _, file := range fileTree.files {
		fileName := fmt.Sprintf(" %v%v", strings.Repeat(tabSpacing, depth), file.name)
		row := table.NewRow(table.RowData{
			columnKeyNumber:   fmt.Sprintf("%v", *fileNumber),
			columnKeyProgress: fmt.Sprintf(" %0.1f", file.percentDone) + "%",
			columnKeySize:     fmt.Sprintf(" %v", utils.HumanizeBytes(file.bytesTotal)),
			columnKeyPriority: fmt.Sprintf(" %v", priorityMap[file.priority]),
			columnKeyFilename: fileName,
		})
		rows = append(rows, row)
		*fileNumber += 1
	}

	return rows
}
