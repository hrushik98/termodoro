package app

import (
	"fmt"
	"time"

	"termodoro/helper"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

// View renders the current state of the application
func (m Model) View() string {
	if m.Quitting {
		return ""
	}

	var view string
	switch m.State {
	case LogoView:
		view = fmt.Sprintf("\n%s  \n\n\n                  Loading...\n",
			TitleStyle.SetString(Logo).Render())

	case ConfigView:
		// Render rows
		var rows [6]string
		labels := []string{
			"Preset:     ",
			"Focus Time: ",
			"Break Time: ",
			"Sound:      ",
			"Animation:  ",
			"",
		}

		// Values
		presetVal := Presets[m.SelectedPreset].Name
		focusVal := fmt.Sprintf("%d mins", m.FocusMinutes)
		breakVal := fmt.Sprintf("%d mins", m.BreakMinutes)
		soundVal := SoundNames[m.SelectedSound]
		animVal := AnimNames[m.SelectedAnim]

		values := []string{
			presetVal,
			focusVal,
			breakVal,
			soundVal,
			animVal,
			"Start Session",
		}

		for i := 0; i < 6; i++ {
			isSelected := (m.ConfigCursor == i)
			var rowStr string

			if i == 5 {
				// Start Button
				if isSelected {
					rowStr = SelectedItemStyle.Render("   ▶ [ START SESSION ] ◀")
				} else {
					rowStr = lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Render("     [ START SESSION ]")
				}
			} else {
				// Label styling
				label := labels[i]
				var valStr string
				if isSelected {
					valStr = SelectedItemStyle.Render(fmt.Sprintf("◀  %s  ▶", values[i]))
					rowStr = fmt.Sprintf("  %s %s%s", SelectedItemStyle.Render("▶"), SelectedItemStyle.Render(label), valStr)
				} else {
					valStr = lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Render(fmt.Sprintf("[ %s ]", values[i]))
					rowStr = fmt.Sprintf("    %s%s", lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Render(label), valStr)
				}
			}
			rows[i] = rowStr
		}

		content := fmt.Sprintf(
			"\n  %s\n\n%s\n%s\n%s\n%s\n%s\n\n%s\n\n\n  %s\n",
			lipgloss.NewStyle().Foreground(lipgloss.Color("#bfedc1")).Bold(true).Render("CONFIGURATION & PRESETS"),
			rows[0],
			rows[1],
			rows[2],
			rows[3],
			rows[4],
			rows[5],
			lipgloss.NewStyle().Faint(true).Render("Use ↑↓ to navigate • ←→ to adjust • Enter to start"),
		)

		view = fmt.Sprintf(
			"%s \n\n%s",
			TitleStyle.Render(),
			PaddingLeftStyle.Render(content),
		)

	case TimerView:
		if AnimNames[m.SelectedAnim] == "BigClock" {
			return renderBigClockView(m)
		}

		totalDuration := time.Duration(m.FocusMinutes) * time.Minute
		if m.IsBreak {
			totalDuration = time.Duration(m.BreakMinutes) * time.Minute
		}

		sessionLabel := "Focus Session"
		if m.IsBreak {
			sessionLabel = "Break Session"
		}

		var timerText string
		var actionText string

		if m.TimedOut {
			timerText = "00 : 00"
			if !m.IsBreak {
				actionText = "Focus Ended! Press [Space] for Break"
			} else {
				actionText = "Break Ended! Press [Space] for Focus"
			}
		} else {
			timerText = formatTime(int(m.Timer.Timeout.Minutes())) + " : " + formatTime(int(m.Timer.Timeout.Seconds())%60)
			actionText = "Session: " + sessionLabel
		}

		view = fmt.Sprintf(
			"%s \n\n %s",
			TitleStyle.Render(),
			PaddingLeftStyle.Render(fmt.Sprintf("%s\n%s%s\n%s",
				helper.Center(timerText, AppWidth-10),
				m.DrawASCII(totalDuration, m.Timer.Timeout),
				helper.Center(actionText, AppWidth-8),
				m.helpView())))
	}

	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, AppStyle.Render(view))
}

// helpView returns the help text for the current state
func (m Model) helpView() string {
	switch m.State {
	case ConfigView:
		return ""
	}

	return "\n" + m.Help.ShortHelpView([]key.Binding{
		m.Keymap.Start,
		m.Keymap.Stop,
		m.Keymap.Reset,
		m.Keymap.Quit,
		m.Keymap.New,
	})
}

// DrawASCII renders the ASCII art animation based on timer progress
func (m Model) DrawASCII(total, remaining time.Duration) string {
	n := m.AsciiArt.Height()
	if n == 0 {
		return ""
	}
	if AnimNames[m.SelectedAnim] != "Flow" {
		return m.AsciiArt.NextAndString(int(percentageDifference(total, remaining)))
	}
	return "\n" + GreenColor.Render(m.AsciiArt.NextAndString(0)) + "\n"
}

// formatTime formats a number as a two-digit string with leading zero
func formatTime(n int) string {
	var b [2]byte
	if n < 10 {
		b[0], b[1] = '0', byte(n)+'0'
	} else {
		b[0], b[1] = byte(n/10)+'0', byte(n%10)+'0'
	}
	return string(b[:])
}

// percentageDifference calculates the percentage of time elapsed
func percentageDifference(total, remaining time.Duration) float64 {
	if total == 0 && remaining == 0 {
		return 0.0
	}
	return ((total.Seconds() - remaining.Seconds()) / total.Seconds()) * 100
}
