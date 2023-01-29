package utils

import (
	"fmt"
	"math"
	"time"

	"github.com/dustin/go-humanize"
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
	// TODO:  should i move the - uint(len(tail)) to LjustText? cause that can
	// lead to compiler optimizations as Ellipsis is a const and as a result
	// len(Ellipsis) is also a const. but it makes sense for it to be here.
	return truncate.StringWithTail(text, maxWidth-uint(len(tail)), tail)
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
// https://stackoverflow.com/questions/71274361/go-error-cannot-use-generic-type-without-instantiation
func HumanizeBytesGeneric[T Number](bytes T) string {
	return HumanizeBytes(uint64(bytes))
}

// FIXME: a better way to do this that is also compatible with FuncMaps
func HumanizeBytes(bytes uint64) string {
	return humanize.Bytes(bytes)
}

func HumanizeTime(torrentTime time.Time) string {
	if torrentTime.Unix() == 0 {
		return "Never"
	} else {
		return fmt.Sprintf("%s (%s)", torrentTime.Format("02/01/2006 03:04:05 PM"), humanize.Time(torrentTime))
	}
}

func HumanizeCorrupt(bytes int64) string {
	if bytes == 0 {
		return "Nothing corrupt"
	}
	return HumanizeBytes(uint64(bytes))
}

func HumanizePrivary(isPrivate bool) string {
	if isPrivate {
		return "Private torrent"
	}
	return "Public torrent"
}

func HumanizeLimit(limit int64, isLimited bool) string {
	if isLimited {
		return HumanizeBytes(uint64(limit * 1024))
	}
	return "No limit"
}

func IntMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func IntMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func IntCeil(a, b int) int {
	return int(math.Ceil(float64(a) / float64(b)))
}

func DivMod(numerator, denominator int) (quotient, remainder int) {
	quotient = numerator / denominator
	remainder = numerator % denominator
	return
}
