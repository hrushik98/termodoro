package app

import (
	"time"

	"github.com/gen2brain/beeep"
)

// Sound option constants
const (
	SoundDefault = iota
	SoundHighBeep
	SoundLowBeep
	SoundDoubleBeep
	SoundMelody
	SoundSilent
)

// SoundNames maps sound index to human-readable names
var SoundNames = []string{
	"Default Alert",
	"High Beep",
	"Low Beep",
	"Double Beep",
	"Melody",
	"Silent",
}

// PlaySound plays a custom sound based on the sound index in a background goroutine
func PlaySound(soundType int) {
	go func() {
		switch soundType {
		case SoundDefault:
			// Default alert/beep
			_ = beeep.Beep(800, 200)
		case SoundHighBeep:
			_ = beeep.Beep(1200, 300)
		case SoundLowBeep:
			_ = beeep.Beep(350, 600)
		case SoundDoubleBeep:
			_ = beeep.Beep(800, 150)
			time.Sleep(100 * time.Millisecond)
			_ = beeep.Beep(800, 150)
		case SoundMelody:
			notes := []float64{523.25, 659.25, 784.00} // C5, E5, G5
			for i, note := range notes {
				_ = beeep.Beep(note, 150)
				if i < len(notes)-1 {
					time.Sleep(80 * time.Millisecond)
				}
			}
		case SoundSilent:
			// Do nothing
		}
	}()
}

// SendNotification sends a desktop notification and plays the configured sound
func SendNotification(title, body string, soundType int) {
	// If silent, just use beeep.Notify (does not force sound)
	if soundType == SoundSilent {
		_ = beeep.Notify(title, body, "assets/logo.png")
		return
	}

	// Send desktop notification
	_ = beeep.Notify(title, body, "assets/logo.png")
	// Play sound
	PlaySound(soundType)
}
