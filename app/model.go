package app

import (
	"os"
	"time"

	"termodoro/ascii_generator"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SessionState represents the current view/state of the application
type SessionState int

const (
	LogoView SessionState = iota
	ConfigView
	TimerView
	TodoView
	NotesView
	StatsView
)

// Preset represents a configuration preset
type Preset struct {
	Name         string
	FocusMinutes int
	BreakMinutes int
	Sound        int
	Animation    int
}

// Built-in presets
var Presets = []Preset{
	{
		Name:         "Classic Pomodoro",
		FocusMinutes: 25,
		BreakMinutes: 5,
		Sound:        SoundMelody,
		Animation:    0, // Tree
	},
	{
		Name:         "Short Focus",
		FocusMinutes: 15,
		BreakMinutes: 3,
		Sound:        SoundHighBeep,
		Animation:    2, // Coffee
	},
	{
		Name:         "Long Focus",
		FocusMinutes: 50,
		BreakMinutes: 10,
		Sound:        SoundDoubleBeep,
		Animation:    1, // Flow
	},
	{
		Name:         "Custom",
		FocusMinutes: 25,
		BreakMinutes: 5,
		Sound:        SoundDefault,
		Animation:    0, // Tree
	},
}

// AnimNames defines the names of animations
var AnimNames = []string{
	"Tree",
	"Flow",
	"Coffee",
	"Campfire",
	"Rain",
	"BigClock",
}

// Model represents the application state
type Model struct {
	State          SessionState
	SelectedPreset int
	FocusMinutes   int
	BreakMinutes   int
	SelectedSound  int
	SelectedAnim   int
	ConfigCursor   int // 0: Preset, 1: Focus, 2: Break, 3: Sound, 4: Anim, 5: Start

	// Timer state
	IsBreak      bool // true if currently in break session, false if in focus session
	Timer        timer.Model
	LoadingTimer timer.Model
	Keymap       Keymap
	Help         help.Model
	Err          error
	Width        int
	Height       int
	AsciiArt     ascii_generator.AsciiArt
	Quitting     bool
	TimedOut     bool

	// Todo List
	Todos       []TodoItem
	TodoInput   textinput.Model
	TodoCursor  int
	TodoAddMode bool
	Username    string

	// Notes
	Notes     textarea.Model
	PrevState SessionState

	// Stats
	StatsData Stats
}

// NewModel creates and returns a new Model with default values
func NewModel() Model {
	// Initialize Todo Input
	todoI := textinput.New()
	todoI.PlaceholderStyle = lipgloss.NewStyle().Faint(true)
	todoI.Placeholder = "Task title..."
	todoI.CharLimit = 100
	todoI.Width = 38
	todoI.Prompt = "> "

	// Initialize Notes textarea
	ta := textarea.New()
	ta.Placeholder = "Jot down your thoughts..."
	ta.CharLimit = 0
	ta.SetWidth(AppWidth - 6)
	ta.SetHeight(14)
	ta.ShowLineNumbers = false

	username := os.Getenv("USER")

	return Model{
		State:          LogoView,
		SelectedPreset: 0, // Classic Pomodoro
		FocusMinutes:   25,
		BreakMinutes:   5,
		SelectedSound:  SoundMelody,
		SelectedAnim:   0, // Tree
		ConfigCursor:   0,

		LoadingTimer: timer.NewWithInterval(800*time.Millisecond, time.Millisecond),
		Err:          nil,
		TimedOut:     false,
		Keymap:       NewKeymap(),
		Help:         help.New(),

		// Todo list
		TodoInput:   todoI,
		Username:    username,
		Todos:       LoadTodos(username),

		// Notes
		Notes: ta,

		// Stats
		StatsData: LoadStats(),
	}
}

// Init implements tea.Model interface
func (m Model) Init() tea.Cmd {
	return tea.Batch(m.LoadingTimer.Init(), m.LoadingTimer.Start())
}

// GenerateASCII creates the appropriate ASCII art based on selected item
func (m Model) GenerateASCII() ascii_generator.AsciiArt {
	animName := AnimNames[m.SelectedAnim]
	switch animName {
	case "Coffee":
		return ascii_generator.GenerateCoffee()
	case "Tree":
		return ascii_generator.GenerateTree(40, 18)
	case "Campfire":
		return ascii_generator.GenerateCampfire()
	case "Rain":
		return ascii_generator.GenerateRain(40, 17)
	default:
		return ascii_generator.GenerateRow(40, 17)
	}
}
