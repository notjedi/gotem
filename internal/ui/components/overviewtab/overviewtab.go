package overviewtab

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/hekmon/transmissionrpc/v2"
	"github.com/notjedi/gotem/internal/ui/common"
	"github.com/notjedi/gotem/internal/ui/utils"
)

const (
	generalInfoTemplate = `
## General Info

Name:               %v
Hash:               %v
ID:                 %v
Location:           %v
Files:              %v
Chunks:             %v;  %v each
***
    `

	sizeInfoTemplate = `
## Size Info

Size:               %v
Downloaded:         %v
Uploaded:           %v
Left until done:    %v
Verified:           %v
Corrupt:            %v
Ratio:              %v
***
    `

	bandwidthInfoTemplate = `
## Bandwidth Info

Download limit:     %v
Upload limit:       %v
Comment:            %v
Creator:            %v
Privacy:            %v
***
    `

	timeInfoTemplate = `
## Time Info

Created at:         %v
Added at:           %v
Started at:         %v
Last activity at:   %v
Completed at:       %v
    `
)

type Model struct {
	hash        string
	id          int64
	torrentInfo transmissionrpc.Torrent
	renderer    *glamour.TermRenderer
}

func New(hash string, id int64) tea.Model {
	// https://stackoverflow.com/questions/49043292/error-template-is-an-incomplete-or-empty-template/49043639#49043639
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithPreservedNewLines(),
	)
	return Model{
		hash:     hash,
		id:       id,
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
	if m.torrentInfo.SizeWhenDone == nil {
		return ""
	}
	// TODO: remove title `Overview`
	// TODO: add status, peers connected to, downloading from, uploading to, seed limit, current
	// status, eta, percentDone, seeds and leeches
	// TODO: only update the content on new message
	t := m.torrentInfo

	generalInfoText := fmt.Sprintf(generalInfoTemplate, *t.Name, *t.HashString, *t.ID,
		*t.DownloadDir, len(t.Files), *t.PieceCount, *t.PieceSize)

	sizeInfoText := fmt.Sprintf(sizeInfoTemplate,
		utils.HumanizeBytes(uint64(t.SizeWhenDone.Byte())),
		utils.HumanizeBytesGeneric(*t.HaveValid),
		utils.HumanizeBytesGeneric(*t.UploadedEver),
		utils.HumanizeBytesGeneric(*t.LeftUntilDone),
		utils.HumanizeBytesGeneric(*t.HaveValid),
		utils.HumanizeBytesGeneric(*t.CorruptEver),
		*t.UploadRatio,
	)

	bandwidthInfoText := fmt.Sprintf(bandwidthInfoTemplate,
		utils.HumanizeLimit(*t.DownloadLimit, *t.DownloadLimited),
		utils.HumanizeLimit(*t.UploadLimit, *t.UploadLimited),
		*t.Comment,
		*t.Creator,
		utils.HumanizePrivary(*t.IsPrivate),
	)

	timeInfoText := fmt.Sprintf(timeInfoTemplate,
		utils.HumanizeTime(*t.DateCreated),
		utils.HumanizeTime(*t.AddedDate),
		utils.HumanizeTime(*t.StartDate),
		utils.HumanizeTime(*t.ActivityDate),
		utils.HumanizeTime(*t.DoneDate),
	)

	out, _ := m.renderer.Render(fmt.Sprintf("%s\n%s\n%s\n%s",
		generalInfoText, sizeInfoText, bandwidthInfoText, timeInfoText))
	return out
}

func (m *Model) SetPrevTorrentInfo(torrentInfo transmissionrpc.Torrent) {
	m.torrentInfo = torrentInfo
}
