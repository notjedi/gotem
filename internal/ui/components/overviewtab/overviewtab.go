package overviewtab

import (
	"fmt"
	"strings"

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

const HeightPadding = 2

type Model struct {
	height  int
	yOffset int

	id       int64
	hash     string
	lines    []string
	renderer *glamour.TermRenderer
}

func New(hash string, id int64, width int, height int) tea.Model {
	// https://stackoverflow.com/questions/49043292/error-template-is-an-incomplete-or-empty-template/49043639#49043639
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithPreservedNewLines(),
		glamour.WithWordWrap(width),
	)
	return Model{
		yOffset:  0,
		id:       id,
		hash:     hash,
		height:   height - HeightPadding,
		renderer: renderer,
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
			m.lines = strings.Split(m.getRenderedInfo(torrentInfo), "\n")
		}

	case tea.WindowSizeMsg:
		m.setWidth(msg.Width)
		m.setHeight(msg.Height)

	case tea.KeyMsg:
		if msg.String() == "j" {
			m.yOffset = utils.IntMin(m.yOffset+1, utils.Abs(len(m.lines)-m.height))
		} else if msg.String() == "k" {
			m.yOffset -= 1
			m.yOffset = utils.IntMax(m.yOffset, 0)
		}
	}
	return m, nil
}

func (m Model) View() string {
	// TODO: add status, peers connected to, downloading from, uploading to, seed limit, current
	// status, eta, percentDone, seeds and leeches
	// TODO: only update the content on new message (cache previous results and only update them if
	// we have received an TorrentInfoMsg message)

	return strings.Join(m.visibleLines(), "\n")
}

func (m *Model) setWidth(width int) {
	m.renderer, _ = glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithPreservedNewLines(),
		glamour.WithWordWrap(width),
	)
}

func (m *Model) setHeight(height int) {
	m.height = height - HeightPadding
	m.yOffset = utils.IntMin(m.yOffset, utils.Abs(len(m.lines)-m.height))
}

func (m *Model) visibleLines() []string {
	if len(m.lines) > m.height {
		top := utils.IntMin(len(m.lines)-m.height, m.yOffset)
		bottom := utils.Clamp(m.yOffset+m.height, top, len(m.lines))
		return m.lines[top:bottom]
	}
	return m.lines
}

func (m *Model) getRenderedInfo(t transmissionrpc.Torrent) string {
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
