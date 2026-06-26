package app

import (
	"fmt"

	"termodoro/helper"

	"github.com/charmbracelet/lipgloss"
)

// Application constants
const (
	DotChar  = " • "
	AppWidth = 50
)

// Style definitions for the application UI
var (
	AppStyle          = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder(), true, true, true, true).Width(AppWidth)
	HeightStyle       = lipgloss.NewStyle().Height(21)
	PaddingLeftStyle  = lipgloss.NewStyle().PaddingLeft(2)
	TitleStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#49beaa")).Bold(true).SetString(helper.Center(`<TERMODORO>`, AppWidth-4)).AlignHorizontal(lipgloss.Center)
	ListTitleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#bfedc1")).PaddingLeft(-10)
	ItemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	SelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#CFF27E"))
	GreenColor        = lipgloss.NewStyle().Foreground(lipgloss.Color("#bfedc1")).PaddingLeft(2).Faint(true)
	DotStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(DotChar)
	BrownColor        = lipgloss.NewStyle().Foreground(lipgloss.Color("#967969"))

	// Todo List priority and done styles
	TodoHighPriorityStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff6b6b")).Bold(true)
	TodoMedPriorityStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffa94d")).Bold(true)
	TodoLowPriorityStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#74c0fc")).Bold(true)
	TodoDoneStyle         = lipgloss.NewStyle().Strikethrough(true).Faint(true)

	// Stats view styles
	StatLabelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#967969"))
	StatValueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#CFF27E")).Bold(true)
	StatBarStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#49beaa"))
	StatBarToday   = lipgloss.NewStyle().Foreground(lipgloss.Color("#CFF27E"))
	StatMutedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("236"))
)

// GitLink is the styled GitHub repository link
var GitLink = GreenColor.Render("https://github.com/hrushik98/termodoro")

// EndInfo is the message displayed when the application exits
var EndInfo = fmt.Sprintf("\n Thanks for using %s! \n Give a star %s \n Fork By     %s \n Original By %s\n",
	lipgloss.NewStyle().Foreground(lipgloss.Color("#49beaa")).Bold(true).Render("<TERMODORO>"),
	GitLink,
	GreenColor.Render("Hrushik Reddy"),
	GreenColor.Render("Sairash Gautam"))
