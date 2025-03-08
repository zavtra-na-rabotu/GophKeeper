package style

import "github.com/charmbracelet/lipgloss"

const (
	Red   = lipgloss.Color("#ff0000")
	White = lipgloss.Color("#ffffff")

	Focused = lipgloss.Color("#ff5faf")
	Blurred = lipgloss.Color("#585858")
)

var (
	ErrorStyle = lipgloss.NewStyle().
		Background(Red).
		Foreground(White)

	FocusedStyle = lipgloss.NewStyle().Foreground(Focused)

	BlurredStyle = lipgloss.NewStyle().Foreground(Blurred)

	NoStyle = lipgloss.NewStyle()
)
