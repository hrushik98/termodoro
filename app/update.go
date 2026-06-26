package app

import (
	"fmt"
	"time"

	"termodoro/helper"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles all state transitions based on messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	case tea.KeyMsg:
		// Open special views from TimerView
		if m.State == TimerView {
			switch msg.String() {
			case "t":
				m.PrevState = TimerView
				m.State = TodoView
				return m, nil
			case "o":
				m.PrevState = TimerView
				m.State = NotesView
				return m, m.Notes.Focus()
			case "T":
				m.PrevState = TimerView
				m.StatsData = LoadStats()
				m.State = StatsView
				return m, nil
			}
		}

		switch msg.Type {
		case tea.KeyCtrlC:
			m.Quitting = true
			return m, tea.Batch(tea.Quit, helper.BeepCmd())

		case tea.KeyCtrlB:
			switch m.State {
			case ConfigView:
				// No back from config screen
			case TimerView:
				m.State = ConfigView
				m.Timer.Stop()
			case TodoView:
				m.TodoAddMode = false
				m.TodoInput.Blur()
				m.TodoInput.Reset()
				m.State = m.PrevState
			case NotesView:
				m.Notes.Blur()
				m.State = m.PrevState
			case StatsView:
				m.State = m.PrevState
			}

		case tea.KeyCtrlT:
			if m.State == ConfigView || m.State == TimerView {
				m.PrevState = m.State
				m.StatsData = LoadStats()
				m.State = StatsView
				return m, nil
			}
		}
	}

	switch m.State {
	case LogoView:
		return updateLogo(msg, m)
	case ConfigView:
		return updateConfig(msg, m)
	case TimerView:
		return updateTimer(msg, m)
	case TodoView:
		return updateTodo(msg, m)
	case NotesView:
		return updateNotes(msg, m)
	case StatsView:
		return updateStats(msg, m)
	default:
		return m, nil
	}
}

// updateLogo handles updates for the logo/splash screen view
func updateLogo(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.LoadingTimer, cmd = m.LoadingTimer.Update(msg)
		return m, cmd
	case timer.TimeoutMsg:
		m.LoadingTimer.Stop()
		m.State = ConfigView
		return m, nil
	}
	return m, nil
}

// updateConfig handles configuration input navigation and updates
func updateConfig(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.ConfigCursor--
			if m.ConfigCursor < 0 {
				m.ConfigCursor = 5
			}

		case "down", "j", "tab":
			m.ConfigCursor++
			if m.ConfigCursor > 5 {
				m.ConfigCursor = 0
			}

		case "left", "h":
			switch m.ConfigCursor {
			case 0: // Preset
				m.SelectedPreset--
				if m.SelectedPreset < 0 {
					m.SelectedPreset = len(Presets) - 1
				}
				// Copy values from preset
				p := Presets[m.SelectedPreset]
				m.FocusMinutes = p.FocusMinutes
				m.BreakMinutes = p.BreakMinutes
				m.SelectedSound = p.Sound
				m.SelectedAnim = p.Animation

			case 1: // Focus Minutes
				if m.FocusMinutes > 1 {
					m.FocusMinutes--
					m.SelectedPreset = len(Presets) - 1 // Custom
				}

			case 2: // Break Minutes
				if m.BreakMinutes > 1 {
					m.BreakMinutes--
					m.SelectedPreset = len(Presets) - 1 // Custom
				}

			case 3: // Sound
				m.SelectedSound--
				if m.SelectedSound < 0 {
					m.SelectedSound = len(SoundNames) - 1
				}
				PlaySound(m.SelectedSound)
				m.SelectedPreset = len(Presets) - 1 // Custom

			case 4: // Animation
				m.SelectedAnim--
				if m.SelectedAnim < 0 {
					m.SelectedAnim = len(AnimNames) - 1
				}
				m.SelectedPreset = len(Presets) - 1 // Custom
			}

		case "right", "l":
			switch m.ConfigCursor {
			case 0: // Preset
				m.SelectedPreset++
				if m.SelectedPreset >= len(Presets) {
					m.SelectedPreset = 0
				}
				// Copy values from preset
				p := Presets[m.SelectedPreset]
				m.FocusMinutes = p.FocusMinutes
				m.BreakMinutes = p.BreakMinutes
				m.SelectedSound = p.Sound
				m.SelectedAnim = p.Animation

			case 1: // Focus Minutes
				if m.FocusMinutes < 120 {
					m.FocusMinutes++
					m.SelectedPreset = len(Presets) - 1 // Custom
				}

			case 2: // Break Minutes
				if m.BreakMinutes < 60 {
					m.BreakMinutes++
					m.SelectedPreset = len(Presets) - 1 // Custom
				}

			case 3: // Sound
				m.SelectedSound++
				if m.SelectedSound >= len(SoundNames) {
					m.SelectedSound = 0
				}
				PlaySound(m.SelectedSound)
				m.SelectedPreset = len(Presets) - 1 // Custom

			case 4: // Animation
				m.SelectedAnim++
				if m.SelectedAnim >= len(AnimNames) {
					m.SelectedAnim = 0
				}
				m.SelectedPreset = len(Presets) - 1 // Custom
			}

		case "enter":
			// Start Focus Timer
			m.IsBreak = false
			duration := time.Duration(m.FocusMinutes) * time.Minute
			m.Timer = timer.NewWithInterval(duration, time.Millisecond)
			m.State = TimerView
			m.AsciiArt = m.GenerateASCII()
			m.TimedOut = false

			body := fmt.Sprintf("Focus timer set for %d minutes.", m.FocusMinutes)
			SendNotification("Focus Session Started", body, m.SelectedSound)

			return m, m.Timer.Init()
		}
	}

	return m, nil
}

// updateTimer handles updates for the timer view
func updateTimer(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.Timer, cmd = m.Timer.Update(msg)
		return m, cmd

	case timer.TimeoutMsg:
		m.TimedOut = true
		if !m.IsBreak {
			RecordSession("Focus", m.FocusMinutes)
		}
		var title, body string
		if !m.IsBreak {
			title = "Focus Session Ended"
			body = fmt.Sprintf("Time for a %d minute break!", m.BreakMinutes)
		} else {
			title = "Break Session Ended"
			body = fmt.Sprintf("Time to focus for %d minutes!", m.FocusMinutes)
		}
		SendNotification(title, body, m.SelectedSound)
		return m, nil

	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.Timer, cmd = m.Timer.Update(msg)
		m.Keymap.Stop.SetEnabled(m.Timer.Running())
		m.Keymap.Start.SetEnabled(!m.Timer.Running())
		return m, tea.Batch(cmd, helper.BeepCmd())

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keymap.Quit):
			m.TimedOut = false
			m.Quitting = true
			return m, tea.Quit

		case key.Matches(msg, m.Keymap.Reset):
			m.TimedOut = false
			m.AsciiArt = m.GenerateASCII()
			duration := time.Duration(m.FocusMinutes) * time.Minute
			if m.IsBreak {
				duration = time.Duration(m.BreakMinutes) * time.Minute
			}
			m.Timer.Timeout = duration

			sessionName := "Focus"
			if m.IsBreak {
				sessionName = "Break"
			}
			body := fmt.Sprintf("%s session restarted.", sessionName)
			SendNotification("Timer Restarted", body, m.SelectedSound)

			return m, m.Timer.Start()

		case key.Matches(msg, m.Keymap.New):
			m.TimedOut = false
			m.State = ConfigView
			m.Timer.Stop()
			return m, nil

		case key.Matches(msg, m.Keymap.Start, m.Keymap.Stop):
			if m.TimedOut {
				// If timed out, Space key transitions to the next session (Focus -> Break or Break -> Focus)
				m.IsBreak = !m.IsBreak
				m.TimedOut = false
				m.AsciiArt = m.GenerateASCII()

				var duration time.Duration
				var title, body string

				if m.IsBreak {
					duration = time.Duration(m.BreakMinutes) * time.Minute
					title = "Break Session Started"
					body = fmt.Sprintf("Break timer set for %d minutes.", m.BreakMinutes)
				} else {
					duration = time.Duration(m.FocusMinutes) * time.Minute
					title = "Focus Session Started"
					body = fmt.Sprintf("Focus timer set for %d minutes.", m.FocusMinutes)
				}

				m.Timer = timer.NewWithInterval(duration, time.Millisecond)
				SendNotification(title, body, m.SelectedSound)
				return m, tea.Batch(m.Timer.Init(), m.Timer.Start())
			}

			// Otherwise toggle pause/play
			return m, m.Timer.Toggle()
		}
	}

	return m, nil
}

// handleBackgroundTimer processes timer ticks and timeout events while in sub-views
func handleBackgroundTimer(msg tea.Msg, m Model) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.Timer, cmd = m.Timer.Update(msg)
		return m, cmd

	case timer.TimeoutMsg:
		m.TimedOut = true
		if !m.IsBreak {
			RecordSession("Focus", m.FocusMinutes)
		}
		var title, body string
		if !m.IsBreak {
			title = "Focus Session Ended"
			body = fmt.Sprintf("Time for a %d minute break!", m.BreakMinutes)
		} else {
			title = "Break Session Ended"
			body = fmt.Sprintf("Time to focus for %d minutes!", m.FocusMinutes)
		}
		SendNotification(title, body, m.SelectedSound)
		m.State = TimerView
		return m, nil
	}
	return m, nil
}

// updateTodo handles updates for the todo list view
func updateTodo(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	// First, let background timer process if necessary
	var timerCmd tea.Cmd
	var newM Model
	newM, timerCmd = handleBackgroundTimer(msg, m)
	if newM.State != TodoView { // State changed (timed out)
		return newM, timerCmd
	}
	m = newM

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.TodoAddMode {
			switch msg.Type {
			case tea.KeyEnter:
				title := m.TodoInput.Value()
				if title != "" {
					m.Todos = append(m.Todos, TodoItem{
						Title:    title,
						Done:     false,
						Priority: PriorityNone,
					})
					SaveTodos(m.Username, m.Todos)
					m.TodoCursor = len(m.Todos) - 1
				}
				m.TodoInput.Reset()
				m.TodoAddMode = false
				m.TodoInput.Blur()
				return m, nil
			case tea.KeyEscape:
				m.TodoInput.Reset()
				m.TodoAddMode = false
				m.TodoInput.Blur()
				return m, nil
			}
			var cmd tea.Cmd
			m.TodoInput, cmd = m.TodoInput.Update(msg)
			return m, tea.Batch(cmd, timerCmd)
		}

		switch msg.String() {
		case "a":
			m.TodoAddMode = true
			return m, tea.Batch(m.TodoInput.Focus(), timerCmd)
		case "d":
			if len(m.Todos) > 0 {
				m.Todos = append(m.Todos[:m.TodoCursor], m.Todos[m.TodoCursor+1:]...)
				SaveTodos(m.Username, m.Todos)
				if m.TodoCursor >= len(m.Todos) && m.TodoCursor > 0 {
					m.TodoCursor--
				}
			}
		case " ":
			if len(m.Todos) > 0 {
				m.Todos[m.TodoCursor].Done = !m.Todos[m.TodoCursor].Done
				SaveTodos(m.Username, m.Todos)
			}
		case "1":
			if len(m.Todos) > 0 {
				m.Todos[m.TodoCursor].Priority = PriorityHigh
				SaveTodos(m.Username, m.Todos)
			}
		case "2":
			if len(m.Todos) > 0 {
				m.Todos[m.TodoCursor].Priority = PriorityMedium
				SaveTodos(m.Username, m.Todos)
			}
		case "3":
			if len(m.Todos) > 0 {
				m.Todos[m.TodoCursor].Priority = PriorityLow
				SaveTodos(m.Username, m.Todos)
			}
		case "0":
			if len(m.Todos) > 0 {
				m.Todos[m.TodoCursor].Priority = PriorityNone
				SaveTodos(m.Username, m.Todos)
			}
		}

		switch msg.Type {
		case tea.KeyUp, tea.KeyLeft:
			if m.TodoCursor > 0 {
				m.TodoCursor--
			}
		case tea.KeyDown, tea.KeyRight:
			if m.TodoCursor < len(m.Todos)-1 {
				m.TodoCursor++
			}
		}
	}

	return m, timerCmd
}

// updateNotes handles updates for the notes view
func updateNotes(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	// First, let background timer process if necessary
	var timerCmd tea.Cmd
	var newM Model
	newM, timerCmd = handleBackgroundTimer(msg, m)
	if newM.State != NotesView { // State changed (timed out)
		return newM, timerCmd
	}
	m = newM

	var cmd tea.Cmd
	m.Notes, cmd = m.Notes.Update(msg)
	return m, tea.Batch(cmd, timerCmd)
}

// updateStats handles updates for the stats view
func updateStats(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	// First, let background timer process if necessary
	var timerCmd tea.Cmd
	var newM Model
	newM, timerCmd = handleBackgroundTimer(msg, m)
	if newM.State != StatsView { // State changed (timed out)
		return newM, timerCmd
	}
	m = newM

	return m, timerCmd
}
