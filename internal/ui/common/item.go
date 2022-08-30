package common

import (
	"fmt"
	"math"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/hekmon/transmissionrpc/v2"
	"github.com/muesli/reflow/padding"
	"github.com/muesli/reflow/truncate"
	"github.com/notjedi/gotem/internal/context"
)

const (
	ellipsis = "…"
)

var (
	speedToStatus = map[bool]string{
		true:  "Downloading",
		false: "Idle",
	}
	/*
	   i can prolly use package specific fields to squeeze a tiny tiny amount of performance.
	   since the response to the request made comes from c, i can kinda assume that adding
	   more fields is basically free. so the performance gain comes down to json serialization and
	   deserialization? and as go is also kinda fast, ig the performance gain here is immeasurable?
	*/
	torrentFields = []string{"id", "hashString", "name", "status", "rateDownload", "rateUpload",
		"eta", "uploadRatio", "sizeWhenDone", "haveValid", "uploadedEver", "recheckProgress",
		"peersConnected", "uploadLimited", "downloadLimited", "bandwidthPriority",
		"peersSendingToUs", "peersGettingFromUs", "seedRatioLimit", "trackerStats", "magnetLink",
		"honorsSessionLimits", "metadataPercentComplete", "percentDone"}
)

type TorrentItem struct {
	item transmissionrpc.Torrent
	ctx  *context.ProgramContext
}

func (t TorrentItem) Title() string {
    titleSpacing := t.ctx.TitleSpacing()

	name := ljustText(*t.item.Name, titleSpacing[0])

	progress := ljustText(fmt.Sprintf("%s / %s", humanize.Bytes(uint64(*t.item.HaveValid)),
		humanize.Bytes(uint64(t.item.SizeWhenDone.Byte()))),
		titleSpacing[1])

	// NOTE: network speeds are in SI standards, we prolly want it in IEC standards
	// if we change this to IEC standards, then it makes sense
	// to change the file sizes to IEC standards too
	networkSpeed := truncateText(fmt.Sprintf("↓ %s  ↑ %s",
		humanize.Bytes(uint64(*t.item.RateDownload)),
		humanize.Bytes(uint64(*t.item.RateUpload))),
		titleSpacing[2], ellipsis)

	return fmt.Sprintf("%s%s%s", name, progress, networkSpeed)
}

func (t TorrentItem) Description() string {
    descSpacing := t.ctx.DescSpacing()

	status := ljustText(t.getStatus(), descSpacing[0])

	uploaded := ljustText(fmt.Sprintf("%s uploaded", humanize.Bytes(uint64(*t.item.UploadedEver))),
		descSpacing[1])

	peersConnected := ljustText(fmt.Sprintf("%d peers connected", *t.item.PeersConnected),
		descSpacing[2])

	seedsAndLeeches := ljustText(fmt.Sprintf("%d seeds %d leeches", t.maxSeeders(),
		t.maxLeechers()), descSpacing[3])

	etaAndRatio := truncateText(fmt.Sprintf("  %s   𢡄 %.2f",
		humanizeDuration(time.Second*time.Duration(*t.item.Eta)),
		math.Max(0.0, *t.item.UploadRatio)), descSpacing[4], ellipsis)

	return fmt.Sprintf("%s%s%s%s%s", status, uploaded, peersConnected, seedsAndLeeches, etaAndRatio)
}

func (t TorrentItem) FilterValue() string {
	return *t.item.Name
}

func (t *TorrentItem) getStatus() string {
	switch *t.item.Status {
	case transmissionrpc.TorrentStatusStopped:
		return "Paused"
	case transmissionrpc.TorrentStatusCheckWait:
		return "Queued to check files"
	case transmissionrpc.TorrentStatusCheck:
		return fmt.Sprintf("Checking files (%.2f%%)", *t.item.RecheckProgress*100)
	case transmissionrpc.TorrentStatusDownloadWait:
		return fmt.Sprintf("Queued to download (%s)", humanize.Ordinal(int(*t.item.QueuePosition)))
	case transmissionrpc.TorrentStatusDownload:
		if *t.item.MetadataPercentComplete == 1.0 {
			statusString := speedToStatus[*t.item.RateDownload > 0.0]
			return fmt.Sprintf("%s (%.2f%%)", statusString, *t.item.PercentDone*100)
		}
		return fmt.Sprintf("Getting metadata (%.2f%%)", *t.item.MetadataPercentComplete*100)
	case transmissionrpc.TorrentStatusSeedWait:
		return fmt.Sprintf("Queued to seed (%s)", humanize.Ordinal(int(*t.item.QueuePosition)))
	case transmissionrpc.TorrentStatusSeed:
		return "Seeding"
	case transmissionrpc.TorrentStatusIsolated:
		return "Isolated"
	default:
		return "Unknown state"
	}
}

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

// taken from - https://gist.github.com/harshavardhana/327e0577c4fed9211f65
func humanizeDuration(duration time.Duration) string {
	if duration.Seconds() < 0.0 {
		return ""
	} else if duration.Seconds() < 60.0 {
		return fmt.Sprintf("%ds", int64(duration.Seconds()))
	} else if duration.Minutes() < 60.0 {
		return fmt.Sprintf("%dm", int64(duration.Minutes()))
	} else if duration.Hours() < 24.0 {
		remainingMinutes := math.Mod(duration.Minutes(), 60)
		return fmt.Sprintf("%dh %dm",
			int64(duration.Hours()), int64(remainingMinutes))
	} else {
		remainingHours := math.Mod(duration.Hours(), 24)
		return fmt.Sprintf("%dd %dh", int64(duration.Hours()/24), int64(remainingHours))
	}
}
