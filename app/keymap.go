package app

import "github.com/charmbracelet/bubbles/key"

// Keymap defines all key bindings for the application
type Keymap struct {
	Start      key.Binding
	Stop       key.Binding
	Reset      key.Binding
	Quit       key.Binding
	New        key.Binding
	Back       key.Binding
	Enter      key.Binding
	CtrlC      key.Binding
	OpenTodo   key.Binding
	TodoAdd    key.Binding
	TodoDelete key.Binding
	TodoToggle key.Binding
	Notes      key.Binding
	Stats      key.Binding
}

// NewKeymap creates and returns a new Keymap with default key bindings
func NewKeymap() Keymap {
	return Keymap{
		Start: key.NewBinding(
			key.WithKeys(" ", "s"),
			key.WithHelp("space", "start"),
		),
		Stop: key.NewBinding(
			key.WithKeys(" ", "s"),
			key.WithHelp("space", "stop"),
		),
		Reset: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "restart"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		New: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new"),
		),
		Back: key.NewBinding(
			key.WithKeys("ctrl+b"),
			key.WithHelp("ctrl+b", "back"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("↩", "enter"),
		),
		CtrlC: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
		OpenTodo: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "todos"),
		),
		TodoAdd: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add"),
		),
		TodoDelete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
		TodoToggle: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle"),
		),
		Notes: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "notes"),
		),
		Stats: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "stats"),
		),
	}
}
