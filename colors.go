package main

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[37m"
	ColorBold   = "\033[1m"
)

func colorize(color, text string) string {
	return color + text + ColorReset
}

func getStatusColor(remaining float64) string {
	if remaining < 0 {
		return ColorRed
	} else if remaining < 10 {
		return ColorYellow
	}
	return ColorGreen
}

func colorStatus(remaining float64, text string) string {
	return colorize(getStatusColor(remaining), text)
}
