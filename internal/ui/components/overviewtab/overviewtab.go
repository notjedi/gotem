package overviewtab

import (
	"bytes"
	"log"
	"text/template"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/hekmon/transmissionrpc/v2"
	"github.com/notjedi/gotem/internal/ui/common"
	"github.com/notjedi/gotem/internal/ui/utils"
)

const (
	templateName     = "overviewTabTemplate.md"
	templateFilePath = "internal/ui/components/overviewtab/overviewTabTemplate.md"
)

// https://stackoverflow.com/questions/71274361/go-error-cannot-use-generic-type-without-instantiation
var funcMap = template.FuncMap{
	"humanizeBytes":   utils.HumanizeBytes,
	"humanizeTime":    utils.HumanizeTime,
	"humanizeCorrupt": utils.HumanizeCorrupt,
}

type Model struct {
	hash        string
	id          int64
	template    *template.Template
	torrentInfo transmissionrpc.Torrent
	renderer    *glamour.TermRenderer
}

func New(hash string, id int64) tea.Model {
	// https://stackoverflow.com/questions/49043292/error-template-is-an-incomplete-or-empty-template/49043639#49043639
	template := template.Must(template.New(templateName).Funcs(funcMap).ParseFiles(templateFilePath))
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithPreservedNewLines(),
	)
	return Model{
		hash:     hash,
		id:       id,
		template: template,
		renderer: renderer,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.TorrentInfoMsg:
		// TODO: should i check if the hash is same?
		m.torrentInfo = transmissionrpc.Torrent(msg)
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	var outputBuffer bytes.Buffer
	if m.torrentInfo.SizeWhenDone == nil {
		return ""
	}
	err := m.template.Execute(&outputBuffer, m.torrentInfo)
	if err != nil {
		log.Fatalln(err)
	}
	bufferString := outputBuffer.String()
	out, _ := m.renderer.Render(bufferString)
	return out
}

func (m *Model) SetPrevTorrentInfo(torrentInfo transmissionrpc.Torrent) {
	m.torrentInfo = torrentInfo
}
