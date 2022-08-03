package listview

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/hekmon/transmissionrpc/v2"
	"github.com/muesli/reflow/padding"
	"github.com/muesli/reflow/truncate"
)

var statusToString = map[transmissionrpc.TorrentStatus]string{
	transmissionrpc.TorrentStatusStopped:      "Paused",
	transmissionrpc.TorrentStatusCheckWait:    "Queued to verify",
	transmissionrpc.TorrentStatusCheck:        "Verifying",
	transmissionrpc.TorrentStatusDownloadWait: "Queued to download",
	transmissionrpc.TorrentStatusDownload:     "Downloading",
	transmissionrpc.TorrentStatusSeedWait:     "Queued to seed",
	transmissionrpc.TorrentStatusSeed:         "Seeding",
	transmissionrpc.TorrentStatusIsolated:     "Isolated",
}

// TODO: take arg to global context instead of copying all spacing values to item object
type TorrentItem struct {
	item         transmissionrpc.Torrent
	titleSpacing [3]uint
	descSpacing  [4]uint
}

func (t TorrentItem) Title() string {
	name := ljustText(*t.item.Name, t.titleSpacing[0])

	progress := ljustText(fmt.Sprintf("%s / %s", humanize.Bytes(uint64(*t.item.HaveValid)),
		humanize.Bytes(uint64((*t.item.SizeWhenDone).Byte()))),
		t.titleSpacing[1])

	// TODO: network speeds are in SI standards, we prolly want it in IEC standards
	// if we change this to IEC standards, then it makes sense
	// to change the file sizes to IEC standards too
	networkSpeed := truncateText(fmt.Sprintf("↓ %s  ↑ %s",
		humanize.Bytes(uint64(*t.item.RateDownload)),
		humanize.Bytes(uint64(*t.item.RateUpload))),
		t.titleSpacing[2], ellipsis)

	return fmt.Sprintf("%s%s%s", name, progress, networkSpeed)
}

// TODO: display eta and seeding ratio
func (t TorrentItem) Description() string {
	var statusString string
	// BUG: change progress to `recheckProgress` if state == verifying
	if *t.item.Status == transmissionrpc.TorrentStatusDownload ||
		*t.item.Status == transmissionrpc.TorrentStatusCheck {
		statusString = fmt.Sprintf("%s (%.2f%%)",
			statusToString[*t.item.Status], *t.item.PercentDone*100)
	} else {
		statusString = statusToString[*t.item.Status]
	}
	status := ljustText(statusString, t.descSpacing[0])

	uploaded := ljustText(fmt.Sprintf("%s uploaded", humanize.Bytes(uint64(*t.item.UploadedEver))),
		t.descSpacing[1])

	peersConnected := ljustText(fmt.Sprintf("%d peers connected", *t.item.PeersConnected),
		t.descSpacing[2])

	seedsAndLeeches := truncateText(fmt.Sprintf("%d seeds %d leeches", t.maxSeeders(),
		t.maxLeechers()), t.descSpacing[3], ellipsis)

	return fmt.Sprintf("%s%s%s%s", status, uploaded, peersConnected, seedsAndLeeches)
}

func (t TorrentItem) FilterValue() string {
	return *t.item.Name
}

/*
should i replace maxSeeders and maxLeechers with this?
https://stackoverflow.com/questions/18930910/access-struct-property-by-name/18931036#18931036
TODO: set seeds and leeches while instantiating the struct? using a New() method prolly?
*/
func (t *TorrentItem) maxSeeders() int64 {
	var max int64
	for _, value := range (*t).item.TrackerStats {
		if max < (*value).SeederCount {
			max = (*value).SeederCount
		}
	}
	return max
}

func (t *TorrentItem) maxLeechers() int64 {
	var max int64
	for _, value := range (*t).item.TrackerStats {
		if max < (*value).LeecherCount {
			max = (*value).LeecherCount
		}
	}
	return max
}

func truncateText(text string, maxWidth uint, tail string) string {
	return truncate.StringWithTail(text, maxWidth, tail)
}

func ljustText(text string, maxWidth uint) string {
	return padding.String(truncateText(text, maxWidth, ellipsis), maxWidth)
}
