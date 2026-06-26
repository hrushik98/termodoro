package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	bigClockOn    = lipgloss.NewStyle().Foreground(lipgloss.Color("#7C6AE8"))
	bigClockDim   = lipgloss.NewStyle().Foreground(lipgloss.Color("#3d3878"))
	bigClockFill  = lipgloss.NewStyle().Foreground(lipgloss.Color("#7C6AE8"))
	bigClockEmpty = lipgloss.NewStyle().Foreground(lipgloss.Color("#2d2b55"))
	bigClockMeta  = lipgloss.NewStyle().Foreground(lipgloss.Color("#9090a0"))
	bigClockHint  = lipgloss.NewStyle().Foreground(lipgloss.Color("#555566"))
)

// 5×7 pixel bitmaps: '#' = filled, ' ' = off
var pixelFont = [10][7]string{
	{"#####", "#   #", "#   #", "#   #", "#   #", "#   #", "#####"}, // 0
	{"  #  ", " ##  ", "  #  ", "  #  ", "  #  ", "  #  ", "#####"}, // 1
	{"#####", "    #", "    #", "#####", "#    ", "#    ", "#####"}, // 2
	{"#####", "    #", "    #", "#####", "    #", "    #", "#####"}, // 3
	{"#   #", "#   #", "#   #", "#####", "    #", "    #", "    #"}, // 4
	{"#####", "#    ", "#    ", "#####", "    #", "    #", "#####"}, // 5
	{"#####", "#    ", "#    ", "#####", "#   #", "#   #", "#####"}, // 6
	{"#####", "    #", "   # ", "  #  ", " #   ", " #   ", " #   "}, // 7
	{"#####", "#   #", "#   #", "#####", "#   #", "#   #", "#####"}, // 8
	{"#####", "#   #", "#   #", "#####", "    #", "    #", "#####"}, // 9
}

// 3×7 colon bitmap
var pixelColon = [7]string{"   ", " # ", " # ", "   ", " # ", " # ", "   "}

// pixelRowToTerm doubles each pixel: '#' → "██", ' ' → "  "
func pixelRowToTerm(row string, style lipgloss.Style) string {
	var sb strings.Builder
	for _, ch := range row {
		if ch == '#' {
			sb.WriteString(style.Render("██"))
		} else {
			sb.WriteString("  ")
		}
	}
	return sb.String()
}

// buildClockRows returns 7 rows for MM:SS in pixel font.
// Visual width per row: (5+1+5)×2 + 4 + 3×2 + 4 = 58 chars.
func buildClockRows(mins, secs int, style lipgloss.Style) [7]string {
	d := [4]int{mins / 10, mins % 10, secs / 10, secs % 10}
	var rows [7]string
	for r := 0; r < 7; r++ {
		rows[r] =
			pixelRowToTerm(pixelFont[d[0]][r], style) + "  " +
				pixelRowToTerm(pixelFont[d[1]][r], style) + "    " +
				pixelRowToTerm(pixelColon[r], style) + "    " +
				pixelRowToTerm(pixelFont[d[2]][r], style) + "  " +
				pixelRowToTerm(pixelFont[d[3]][r], style)
	}
	return rows
}

// renderBigClockView renders a full-screen timer, bypassing the boxed AppStyle.
func renderBigClockView(m Model) string {
	w := m.Width
	if w < 70 {
		w = 70
	}

	// Remaining time
	remaining := m.Timer.Timeout
	if m.TimedOut {
		remaining = 0
	}
	mins := int(remaining.Minutes())
	secs := int(remaining.Seconds()) % 60
	if mins > 99 {
		mins, secs = 99, 59
	}

	// Dim digits when paused
	digitStyle := bigClockOn
	if !m.Timer.Running() && !m.TimedOut {
		digitStyle = bigClockDim
	}

	clockRows := buildClockRows(mins, secs, digitStyle)

	// center uses lipgloss Width so ANSI codes are handled correctly
	center := func(s string) string {
		return lipgloss.NewStyle().Width(w).AlignHorizontal(lipgloss.Center).Render(s)
	}

	var sb strings.Builder

	// ── Big digit rows ────────────────────────────────────────────────
	for _, row := range clockRows {
		sb.WriteString(center(row) + "\n")
	}

	// ── Session label ─────────────────────────────────────────────────
	sb.WriteString("\n\n")
	label := "focus session"
	if m.IsBreak {
		label = "break session"
	}
	switch {
	case m.TimedOut && !m.IsBreak:
		label = "focus ended  ·  press space to start break"
	case m.TimedOut && m.IsBreak:
		label = "break ended  ·  press space to start focus"
	case !m.Timer.Running():
		label = "paused"
	}
	sb.WriteString(bigClockMeta.Width(w).AlignHorizontal(lipgloss.Center).Render(label) + "\n")

	// ── Progress bar ──────────────────────────────────────────────────
	totalDuration := time.Duration(m.FocusMinutes) * time.Minute
	if m.IsBreak {
		totalDuration = time.Duration(m.BreakMinutes) * time.Minute
	}
	pct := int(percentageDifference(totalDuration, remaining))
	if pct < 0 {
		pct = 0
	} else if pct > 100 {
		pct = 100
	}

	barW := w - 10
	if barW < 10 {
		barW = 10
	}
	filled := pct * barW / 100
	bar := bigClockFill.Render(strings.Repeat("█", filled)) +
		bigClockEmpty.Render(strings.Repeat("░", barW-filled)) +
		fmt.Sprintf("  %3d%%", pct)
	sb.WriteString(center(bar) + "\n")

	// ── Key hints ─────────────────────────────────────────────────────
	sb.WriteString("\n")
	hints := "space pause/resume  ·  r reset  ·  n new session  ·  q quit"
	sb.WriteString(bigClockHint.Width(w).AlignHorizontal(lipgloss.Center).Render(hints) + "\n")

	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, sb.String())
}
