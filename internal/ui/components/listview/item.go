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
var titleSpacingRatios = [...]float32{0.75, 0.15, 0.10}
var descSpacingRatios = [...]float32{0.25, 0.25, 0.25, 0.25}

// TODO: make this ratio global like context or something? instead of recomputing it over and over
type TorrentItem struct {
	item  transmissionrpc.Torrent
	width float32
}

// TODO: create wrapper func to pad and truncate string
func (t TorrentItem) Title() string {
	name := padding.String(truncate.StringWithTail(*t.item.Name,
		uint(titleSpacingRatios[0]*t.width), ellipsis), uint(titleSpacingRatios[0]*t.width))

	progress := padding.String(fmt.Sprintf("%s / %s", humanize.Bytes(uint64(*t.item.HaveValid)),
		humanize.Bytes(uint64((*t.item.SizeWhenDone).Byte()))), uint(titleSpacingRatios[1]*t.width))

	networkSpeed := padding.String(fmt.Sprintf("↓ %s  ↑ %s",
		humanize.Bytes(uint64(*t.item.RateDownload)),
		humanize.Bytes(uint64(*t.item.RateUpload))), uint(titleSpacingRatios[2]*t.width))

	return fmt.Sprintf("%s%s%s", name, progress, networkSpeed)
}

func (t TorrentItem) Description() string {
	status := padding.String(truncate.StringWithTail(statusToString[*t.item.Status],
		uint(descSpacingRatios[0]*t.width), ellipsis), uint(descSpacingRatios[0]*t.width))

	uploaded := padding.String(
		fmt.Sprintf("%s uploaded", humanize.Bytes(uint64(*t.item.UploadedEver))),
		uint(descSpacingRatios[1]*t.width))

	peersConnected := padding.String(fmt.Sprintf("%d peers connected", *t.item.PeersConnected),
		uint(descSpacingRatios[2]*t.width))

	seedsAndLeeches := padding.String(fmt.Sprintf("%d seeds %d leeches", t.maxSeeders(),
		t.maxLeechers()), uint(descSpacingRatios[3]*t.width))

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
