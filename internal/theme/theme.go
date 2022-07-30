package theme

import "github.com/charmbracelet/lipgloss"

// Theme represents the properties that make up a theme.
type Theme struct {
	StatusbarSelectedFileForegroundColor lipgloss.AdaptiveColor
	StatusbarSelectedFileBackgroundColor lipgloss.AdaptiveColor
	StatusbarBarForegroundColor          lipgloss.AdaptiveColor
	StatusbarBarBackgroundColor          lipgloss.AdaptiveColor
	StatusbarTotalFilesForegroundColor   lipgloss.AdaptiveColor
	StatusbarTotalFilesBackgroundColor   lipgloss.AdaptiveColor
	StatusbarLogoForegroundColor         lipgloss.AdaptiveColor
	StatusbarLogoBackgroundColor         lipgloss.AdaptiveColor
}

var themeMap = map[string]Theme{
	"default": {
		StatusbarSelectedFileForegroundColor: lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
		StatusbarSelectedFileBackgroundColor: lipgloss.AdaptiveColor{Dark: "#F25D94", Light: "#F25D94"},
		StatusbarBarForegroundColor:          lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
		StatusbarBarBackgroundColor:          lipgloss.AdaptiveColor{Dark: "#3c3836", Light: "#3c3836"},
		StatusbarTotalFilesForegroundColor:   lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
		StatusbarTotalFilesBackgroundColor:   lipgloss.AdaptiveColor{Dark: "#A550DF", Light: "#A550DF"},
		StatusbarLogoForegroundColor:         lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
		StatusbarLogoBackgroundColor:         lipgloss.AdaptiveColor{Dark: "#6124DF", Light: "#6124DF"},
	},
}

func GetTheme(theme string) Theme {
	switch theme {
	case "default":
		return themeMap["default"]
	default:
		return themeMap["default"]
	}
}
