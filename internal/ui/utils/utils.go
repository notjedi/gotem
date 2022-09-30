package utils

import (
	"fmt"
	"math"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/hekmon/transmissionrpc/v2"
	"github.com/muesli/reflow/padding"
	"github.com/muesli/reflow/truncate"
	"golang.org/x/exp/constraints"
)

const (
	Ellipsis = "…"
)

type Number interface {
	constraints.Integer | constraints.Float
}

func TruncateText(text string, maxWidth uint, tail string) string {
	return truncate.StringWithTail(text, maxWidth, tail)
}

func LjustText(text string, maxWidth uint) string {
	return padding.String(TruncateText(text, maxWidth, Ellipsis), maxWidth)
}

// taken from - https://gist.github.com/harshavardhana/327e0577c4fed9211f65
func HumanizeDuration(duration time.Duration) string {
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

// https://stackoverflow.com/questions/67678331/how-to-write-a-generic-function-that-accepts-any-numerical-type
func HumanizeBytesGeneric[T Number](bytes T) string {
	return humanize.Bytes(uint64(bytes))
}

// FIXME: a better way to do this that is also compatible with FuncMaps
func HumanizeBytes(bytes interface{}) string {
	switch val := bytes.(type) {
	case *int64:
		return humanize.Bytes(uint64(*val))
	case int64:
		return humanize.Bytes(uint64(val))
	case *float64:
		return humanize.Bytes(uint64(*val))
	case float64:
		return humanize.Bytes(uint64(val))
	case *uint64:
		return humanize.Bytes(*val)
	case uint64:
		return humanize.Bytes(val)
	default:
		return fmt.Sprintf("%T", bytes)
	}
}

func HumanizeTime(torrentTime *time.Time) string {
	if torrentTime.Unix() == 0 {
		return "Never"
	} else {
		return fmt.Sprintf("%s (%s)", torrentTime.Format("02/01/2006 03:04:05 PM"), humanize.Time(*torrentTime))
	}
}

func HumanizeCorrupt(bytes *int64) string {
	if *bytes == 0 {
		return "Nothing corrupt"
	}
	return humanize.Bytes(uint64(*bytes))
}

func HumanizePrivary(isPrivate *bool) string {
	if *isPrivate {
		return "Private torrent"
	}
	return "Public torrent"
}

func HumanizeDownloadLimit(t transmissionrpc.Torrent) string {
	if *t.DownloadLimited {
		return string(*t.DownloadLimit)
	}
	return "No limit"
}

func HumanizeUploadLimit(t transmissionrpc.Torrent) string {
	if *t.UploadLimited {
		return string(*t.UploadLimit)
	}
	return "No limit"
}
