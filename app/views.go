package app

import (
	"fmt"
	"strings"
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
		view = fmt.Sprintf(
			"%s \n\n\n%s",
			TitleStyle.Render(),
			lipgloss.NewStyle().Faint(true).Render(helper.Center("Loading...", AppWidth-8)),
		)

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

		selectedColor := lipgloss.NewStyle().Foreground(lipgloss.Color("#CFF27E"))

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
					valStr = selectedColor.Render(fmt.Sprintf("◀  %s  ▶", values[i]))
					rowStr = fmt.Sprintf("  ▶ %s%s", selectedColor.Render(label), valStr)
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
			lipgloss.NewStyle().Faint(true).Render("Use ↑↓ to navigate • ←→ to adjust • Enter to start • Ctrl+T stats"),
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

	case TodoView:
		view = fmt.Sprintf(
			"%s \n\n%s",
			TitleStyle.Render(),
			PaddingLeftStyle.Render(
				fmt.Sprintf("%s\n%s",
					HeightStyle.Render(m.todoListView()),
					m.helpView())))

	case NotesView:
		sessionLabel := "Focus"
		if m.IsBreak {
			sessionLabel = "Break"
		}
		view = fmt.Sprintf(
			"%s \n\n%s",
			TitleStyle.Render(),
			PaddingLeftStyle.Render(fmt.Sprintf("%s\n\n%s\n\n%s",
				BrownColor.Render("Notes — "+sessionLabel+":"),
				m.Notes.View(),
				m.helpView(),
			)))

	case StatsView:
		view = fmt.Sprintf(
			"%s \n\n%s",
			TitleStyle.Render(),
			PaddingLeftStyle.Render(fmt.Sprintf("%s\n\n%s",
				m.statsBody(),
				m.helpView(),
			)))
	}

	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, AppStyle.Render(view))
}

// todoListView renders the todo list content
func (m Model) todoListView() string {
	taskCount := fmt.Sprintf("  %d task", len(m.Todos))
	if len(m.Todos) != 1 {
		taskCount += "s"
	}
	header := ListTitleStyle.Render("Todo List") + " " + BrownColor.Render(taskCount) + "\n\n"

	var items string
	if len(m.Todos) == 0 {
		items = GreenColor.Render("No tasks yet. Press 'a' to add one.") + "\n"
	} else {
		for i, todo := range m.Todos {
			checkbox := "[ ]"
			if todo.Done {
				checkbox = "[✓]"
			}

			title := todo.Title
			if todo.Done {
				title = TodoDoneStyle.Render(title)
			}

			var prio string
			switch todo.Priority {
			case PriorityHigh:
				prio = " " + TodoHighPriorityStyle.Render("P1")
			case PriorityMedium:
				prio = " " + TodoMedPriorityStyle.Render("P2")
			case PriorityLow:
				prio = " " + TodoLowPriorityStyle.Render("P3")
			}

			line := fmt.Sprintf("%s %s%s", checkbox, title, prio)
			if i == m.TodoCursor {
				line = SelectedItemStyle.Render(line)
			} else {
				line = ItemStyle.Render(line)
			}
			items += line + "\n"
		}
	}

	var bottom string
	if m.TodoAddMode {
		bottom = "\n" + BrownColor.Render("New task:") + "\n" + m.TodoInput.View() + "\n"
	} else {
		bottom = "\n" + GreenColor.Render("1/2/3 set priority · 0 clear") + "\n"
	}

	return header + items + bottom
}

// statsBody renders the stats dashboard: headline metrics plus a 7-day bar chart.
func (m Model) statsBody() string {
	s := m.StatsData
	if s.TotalSessions() == 0 {
		return BrownColor.Render("📊 Stats\n\n") +
			StatMutedStyle.Render("No sessions yet — finish a timer\nto start tracking your focus!")
	}

	var b strings.Builder

	b.WriteString(BrownColor.Render("📊 Your Focus Stats"))
	b.WriteString("\n\n")

	b.WriteString(statLine("🔥 Streak", pluralize(s.Streak(), "day", "days")))
	b.WriteString(statLine("⏱  Total", formatMinutes(s.TotalMinutes())))
	b.WriteString(statLine("✅ Sessions", fmt.Sprintf("%d", s.TotalSessions())))
	b.WriteString(statLine("📅 Today",
		fmt.Sprintf("%d · %s", s.TodaySessions(), formatMinutes(s.TodayMinutes()))))

	b.WriteString("\n")
	b.WriteString(BrownColor.Render("Last 7 days"))
	b.WriteString("\n")
	b.WriteString(renderBarChart(s.Last7Days()))

	return b.String()
}

// statLine formats a single label/value metric row.
func statLine(label, value string) string {
	return StatLabelStyle.Render(fmt.Sprintf("%-12s", label)) +
		StatValueStyle.Render(value) + "\n"
}

// renderBarChart renders a horizontal bar chart for the trailing week.
func renderBarChart(days []DayStat) string {
	const maxWidth = 16

	max := 0
	for _, d := range days {
		if d.Minutes > max {
			max = d.Minutes
		}
	}

	var b strings.Builder
	for _, d := range days {
		bar := ""
		if max > 0 && d.Minutes > 0 {
			filled := (d.Minutes*maxWidth + max/2) / max
			if filled == 0 {
				filled = 1
			}
			bar = strings.Repeat("█", filled)
		}

		barStyle := StatBarStyle
		if d.Today {
			barStyle = StatBarToday
		}

		label := StatLabelStyle.Render(fmt.Sprintf("%-4s", d.Label))
		value := ""
		if d.Minutes > 0 {
			value = " " + StatMutedStyle.Render(formatMinutes(d.Minutes))
		}
		b.WriteString(label + barStyle.Render(bar) + value + "\n")
	}
	return b.String()
}

// formatMinutes renders a minute count as "Xh Ym" or "Ym".
func formatMinutes(total int) string {
	h := total / 60
	mm := total % 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, mm)
	}
	return fmt.Sprintf("%dm", mm)
}

// pluralize formats a count with the right singular/plural unit.
func pluralize(n int, singular, plural string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, singular)
	}
	return fmt.Sprintf("%d %s", n, plural)
}

// helpView returns the help text for the current state
func (m Model) helpView() string {
	switch m.State {
	case ConfigView:
		return ""
	case TodoView:
		return "\n" + m.Help.ShortHelpView([]key.Binding{
			m.Keymap.TodoAdd,
			m.Keymap.TodoToggle,
			m.Keymap.TodoDelete,
			m.Keymap.Back,
		})
	case NotesView, StatsView:
		return "\n" + m.Help.ShortHelpView([]key.Binding{
			m.Keymap.Back,
		})
	}

	return "\n" + m.Help.ShortHelpView([]key.Binding{
		m.Keymap.Start,
		m.Keymap.Stop,
		m.Keymap.Reset,
		m.Keymap.Quit,
		m.Keymap.New,
		m.Keymap.OpenTodo,
		m.Keymap.Notes,
		m.Keymap.Stats,
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
