package listview

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/hekmon/transmissionrpc/v2"
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

type TorrentItem struct {
	item transmissionrpc.Torrent
}

func (t TorrentItem) Title() string {
	return fmt.Sprintf("%s      %s / %s     %s  %s",
		*t.item.Name, humanize.Bytes(uint64(*t.item.HaveValid)),
		humanize.Bytes(uint64((*t.item.SizeWhenDone).Byte())),
		humanize.Bytes(uint64(*t.item.RateDownload)),
		humanize.Bytes(uint64(*t.item.RateUpload)),
	)
}

func (t TorrentItem) Description() string {
	return fmt.Sprintf("%s      %s   %d peers connected     %d seeds  %d leeches",
		statusToString[*t.item.Status],
		humanize.Bytes(uint64(*t.item.UploadedEver)), *t.item.PeersConnected,
		t.maxSeeders(), t.maxLeechers(),
	)
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
