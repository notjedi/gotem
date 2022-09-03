package utils

import (
	"fmt"
	"time"
    "math"

	"github.com/muesli/reflow/padding"
	"github.com/muesli/reflow/truncate"
)

const (
	Ellipsis = "…"
)

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
