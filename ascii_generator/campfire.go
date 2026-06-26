package ascii_generator

import (
	"sync"
	"time"

	"termodoro/helper"

	"github.com/charmbracelet/lipgloss"
)

// Campfire represents a flickering ASCII campfire animation
type Campfire struct {
	cached   string
	curFrame int
	width    int
	height   int
	mu       sync.RWMutex
}

var (
	flameTipStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700"))
	flameMidStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B35"))
	flameBaseStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4500"))
	campLogStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#8B4513"))
	emberStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF8C00"))
)

// 4 frames: first 3 lines are flame layers, last 3 are the log structure
var campfireFrames = [][]string{
	{
		`    )    (    )  `,
		`  (   )  (  ) ( `,
		`   ) ( ) ( ) (  `,
		`  _______________`,
		`  |  === ===  |  `,
		`  |___________|  `,
	},
	{
		`      (    )   ( `,
		`    )   (   )    `,
		`   ( )   ) ( (   `,
		`  _______________`,
		`  |  === ===  |  `,
		`  |___________|  `,
	},
	{
		`    )    (   )   `,
		`  (   )  )  ) (  `,
		`   ) (  ) ( (    `,
		`  _______________`,
		`  |  === ===  |  `,
		`  |___________|  `,
	},
	{
		`     (   )   (   `,
		`   )   (  )  )   `,
		`   ( )  ) ( ) (  `,
		`  _______________`,
		`  |  === ===  |  `,
		`  |___________|  `,
	},
}

func (c *Campfire) Width() int  { return c.width }
func (c *Campfire) Height() int { return c.height }

func (c *Campfire) NextAndString(_ int) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cached
}

func (c *Campfire) advance() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.curFrame = (c.curFrame + 1) % len(campfireFrames)
	f := campfireFrames[c.curFrame]
	w := c.width

	c.cached = "\n\n\n" +
		flameTipStyle.Render(helper.Center(f[0], w)) + "\n" +
		flameMidStyle.Render(helper.Center(f[1], w)) + "\n" +
		flameBaseStyle.Render(helper.Center(f[2], w)) + "\n" +
		campLogStyle.Render(helper.Center(f[3], w)) + "\n" +
		emberStyle.Render(helper.Center(f[4], w)) + "\n" +
		campLogStyle.Render(helper.Center(f[5], w)) + "\n\n\n"
}

// GenerateCampfire creates a new Campfire animation
func GenerateCampfire() *Campfire {
	c := &Campfire{width: 40, height: 6, curFrame: -1}
	c.advance()

	go func() {
		t := time.NewTicker(250 * time.Millisecond)
		defer t.Stop()
		for range t.C {
			c.advance()
		}
	}()

	return c
}
