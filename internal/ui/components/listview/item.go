package listview

import (
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
	return *t.item.Name
}

func (t TorrentItem) Description() string {
	return statusToString[*t.item.Status]
}

func (t TorrentItem) FilterValue() string {
	return *t.item.Name
}
