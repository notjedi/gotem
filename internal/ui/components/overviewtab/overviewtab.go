package overviewtab

import (
	"bytes"
	"log"
	"text/template"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/hekmon/transmissionrpc/v2"
	"github.com/notjedi/gotem/internal/ui/common"
)

type Model struct {
	hash            string
	id              int64
	template        *template.Template
	prevTorrentInfo transmissionrpc.Torrent
	renderer        *glamour.TermRenderer
}

func New(hash string, id int64) tea.Model {
	template := template.Must(template.ParseFiles("internal/ui/components/overviewtab/overviewTabTemplate.txt"))
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
		m.prevTorrentInfo = transmissionrpc.Torrent(msg)
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	var outputBuffer bytes.Buffer
	err := m.template.Execute(&outputBuffer, m.prevTorrentInfo)
	if err != nil {
		log.Fatalln(err)
	}
	bufferString := outputBuffer.String()
	out, _ := m.renderer.Render(bufferString)
	return out
}

func (m *Model) SetPrevTorrentInfo(torrentInfo transmissionrpc.Torrent) {
	m.prevTorrentInfo = torrentInfo
}
