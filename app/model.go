package app

import (
	"time"

	"termodoro/ascii_generator"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

// SessionState represents the current view/state of the application
type SessionState int

const (
	LogoView SessionState = iota
	ConfigView
	TimerView
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
	"Sunrise",
	"BigClock",
}

// noopArt satisfies AsciiArt for animations handled entirely in the view layer
type noopArt struct{}

func (noopArt) Width() int            { return 0 }
func (noopArt) Height() int           { return 0 }
func (noopArt) NextAndString(_ int) string { return "" }

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
}

// NewModel creates and returns a new Model with default values
func NewModel() Model {
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
	}
}

// Init implements tea.Model interface
func (m Model) Init() tea.Cmd {
	return tea.Batch(m.LoadingTimer.Init(), m.LoadingTimer.Start())
}

// GenerateASCII creates the appropriate ASCII art based on selected animation
func (m Model) GenerateASCII() ascii_generator.AsciiArt {
	switch AnimNames[m.SelectedAnim] {
	case "Coffee":
		return ascii_generator.GenerateCoffee()
	case "Tree":
		return ascii_generator.GenerateTree(40, 18)
	case "Campfire":
		return ascii_generator.GenerateCampfire()
	case "Rain":
		return ascii_generator.GenerateRain(40, 15)
	case "Sunrise":
		return ascii_generator.GenerateSunrise(40, 15)
	case "BigClock":
		return noopArt{}
	default: // Flow
		return ascii_generator.GenerateRow(40, 17)
	}
}
