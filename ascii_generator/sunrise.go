package ascii_generator

import (
	"strings"

	"termodoro/helper"

	"github.com/charmbracelet/lipgloss"
)

// Sunrise represents a sun rising from the horizon as timer percentage increases (0%→100%)
type Sunrise struct {
	width  int
	height int
}

var (
	sunRedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4500"))
	sunOrangeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF8C00"))
	sunYellowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700"))
	skyDawnStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#4a3728"))
	skyDayStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#5dade2"))
	horizonLine    = lipgloss.NewStyle().Foreground(lipgloss.Color("#27AE60"))
	groundLine     = lipgloss.NewStyle().Foreground(lipgloss.Color("#795548"))
)

func (s *Sunrise) Width() int  { return s.width }
func (s *Sunrise) Height() int { return s.height }

func (s *Sunrise) NextAndString(percent int) string {
	horizonRow := s.height - 2
	groundRow := s.height - 1

	// sun rises from just above horizon (0%) to near top (100%)
	maxSunRow := horizonRow - 2 // lowest visible position (body row)
	minSunRow := 1              // highest visible position
	sunRow := maxSunRow - int(float64(percent)/100.0*float64(maxSunRow-minSunRow))

	// sun color shifts red→orange→yellow as it rises
	var sunStyle lipgloss.Style
	switch {
	case percent < 33:
		sunStyle = sunRedStyle
	case percent < 66:
		sunStyle = sunOrangeStyle
	default:
		sunStyle = sunYellowStyle
	}

	var skyStyle lipgloss.Style
	if percent < 50 {
		skyStyle = skyDawnStyle
	} else {
		skyStyle = skyDayStyle
	}

	w := s.width
	skyFill := helper.Center(strings.Repeat("-", w-4), w)

	var sb strings.Builder
	sb.WriteString("\n")

	for y := 0; y < s.height; y++ {
		var line string
		switch y {
		case sunRow - 1:
			line = sunStyle.Render(helper.Center(`\  |  /`, w))
		case sunRow:
			line = sunStyle.Render(helper.Center(`--(o)--`, w))
		case sunRow + 1:
			line = sunStyle.Render(helper.Center(`/  |  \`, w))
		case horizonRow:
			line = horizonLine.Render(helper.Center(strings.Repeat("~", w-4), w))
		case groundRow:
			line = groundLine.Render(helper.Center(strings.Repeat("=", w-4), w))
		default:
			line = skyStyle.Render(skyFill)
		}
		sb.WriteString(line + "\n")
	}

	sb.WriteString("\n")
	return sb.String()
}

// GenerateSunrise creates a new Sunrise animation with the given dimensions
func GenerateSunrise(width, height int) *Sunrise {
	return &Sunrise{width: width, height: height}
}
