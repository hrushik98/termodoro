package ascii_generator

import (
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Rain represents falling raindrops scrolling down the screen
type Rain struct {
	rows   []string
	cached string
	width  int
	height int
	mu     sync.RWMutex
}

var (
	rainHeavy  = lipgloss.NewStyle().Foreground(lipgloss.Color("#74b9ff"))
	rainMedium = lipgloss.NewStyle().Foreground(lipgloss.Color("#0984e3"))
	rainLight  = lipgloss.NewStyle().Foreground(lipgloss.Color("#b2bec3"))
)

// dropChars are the characters used to render raindrops
var dropChars = []rune{'|', '/', '\'', '.'}

func newRainRow(width int) string {
	row := make([]rune, width)
	for i := range row {
		row[i] = ' '
	}
	count := rand.Intn(width/5) + 2
	for i := 0; i < count; i++ {
		x := rand.Intn(width)
		row[x] = dropChars[rand.Intn(len(dropChars))]
	}
	return string(row)
}

func (r *Rain) Width() int  { return r.width }
func (r *Rain) Height() int { return r.height }

func (r *Rain) NextAndString(_ int) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.cached
}

var rainStyles = []*lipgloss.Style{&rainHeavy, &rainMedium, &rainLight}

func (r *Rain) tick() {
	r.mu.Lock()
	defer r.mu.Unlock()

	// shift rows down: new row appears at top, last row drops off
	r.rows = append([]string{newRainRow(r.width)}, r.rows[:r.height-1]...)

	var sb strings.Builder
	sb.WriteString("\n")
	for i, row := range r.rows {
		style := rainStyles[i%len(rainStyles)]
		sb.WriteString(style.Render(row) + "\n")
	}
	sb.WriteString("\n")
	r.cached = sb.String()
}

// GenerateRain creates a new Rain animation with the given dimensions
func GenerateRain(width, height int) *Rain {
	r := &Rain{
		width:  width,
		height: height,
		rows:   make([]string, height),
	}
	for i := range r.rows {
		r.rows[i] = newRainRow(width)
	}
	r.tick()

	go func() {
		t := time.NewTicker(100 * time.Millisecond)
		defer t.Stop()
		for range t.C {
			r.tick()
		}
	}()

	return r
}
