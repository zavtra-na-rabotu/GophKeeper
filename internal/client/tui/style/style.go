package style

import "github.com/charmbracelet/lipgloss"

const (
	Red   = lipgloss.Color("#ff0000")
	White = lipgloss.Color("#ffffff")

	Focused = lipgloss.Color("205")
	Blurred = lipgloss.Color("240")
)

var (
	ErrorStyle = lipgloss.NewStyle().
		Background(Red).
		Foreground(White)

	FocusedStyle = lipgloss.NewStyle().Foreground(Focused)

	BlurredStyle = lipgloss.NewStyle().Foreground(Blurred)

	NoStyle = lipgloss.NewStyle()
)
