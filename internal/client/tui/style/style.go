package style

import "github.com/charmbracelet/lipgloss"

const (
	Black   = lipgloss.Color("#000000")
	Red     = lipgloss.Color("#ff0000")
	White   = lipgloss.Color("#ffffff")
	Focused = lipgloss.Color("#ff5faf")
	Blurred = lipgloss.Color("#585858")
)

var (
	ErrorStyle = lipgloss.NewStyle().
			Background(Red).
			Foreground(White)
	HintStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#ffadad")).
			Foreground(Black)
	FocusedStyle = lipgloss.NewStyle().Foreground(Focused)
	BlurredStyle = lipgloss.NewStyle().Foreground(Blurred)
	NoStyle      = lipgloss.NewStyle()
)
