package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	engine *Engine
	msg    string
	ampA   float64
	ampB   float64
}

func NewModel(engine *Engine) Model {
	return Model{engine: engine}
}

func (m Model) Init() tea.Cmd {
	// UI ticker
	return tea.Tick(time.Millisecond*80, func(t time.Time) tea.Msg { return t })
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case time.Time:
		if m.engine != nil {
			if m.engine.deckA != nil {
				m.ampA = m.engine.deckA.RMS()
			}
			if m.engine.deckB != nil {
				m.ampB = m.engine.deckB.RMS()
			}
		}
		return m, tea.Tick(time.Millisecond*80, func(t time.Time) tea.Msg { return t })
	case tea.KeyMsg:
		switch v.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		// deck A
		case "1":
			if m.engine != nil && m.engine.deckA != nil {
				m.engine.deckA.Toggle()
			}
		case "a":
			if m.engine != nil && m.engine.deckA != nil {
				// Speed changes are now just for display since beep doesn't support it easily
				m.msg = "Speed control: Limited with current audio backend"
			}
		case "s":
			if m.engine != nil && m.engine.deckA != nil {
				// Speed changes are now just for display since beep doesn't support it easily
				m.msg = "Speed control: Limited with current audio backend"
			}
		case "z":
			if m.engine != nil && m.engine.deckA != nil {
				// Volume control now works with the new audio engine
				m.engine.deckA.SetVolume(0.8)
				m.engine.updateGains()
				m.msg = "Deck A: Volume decreased"
			}
		case "x":
			if m.engine != nil && m.engine.deckA != nil {
				// Volume control now works with the new audio engine
				m.engine.deckA.SetVolume(1.2)
				m.engine.updateGains()
				m.msg = "Deck A: Volume increased"
			}

		// deck B
		case "2":
			if m.engine != nil && m.engine.deckB != nil {
				m.engine.deckB.Toggle()
			}
		case "k":
			if m.engine != nil && m.engine.deckB != nil {
				// Speed changes are now just for display since beep doesn't support it easily
				m.msg = "Speed control: Limited with current audio backend"
			}
		case "l":
			if m.engine != nil && m.engine.deckB != nil {
				// Speed changes are now just for display since beep doesn't support it easily
				m.msg = "Speed control: Limited with current audio backend"
			}
		case "n":
			if m.engine != nil && m.engine.deckB != nil {
				// Volume control now works with the new audio engine
				m.engine.deckB.SetVolume(0.8)
				m.engine.updateGains()
				m.msg = "Deck B: Volume decreased"
			}
		case "m":
			if m.engine != nil && m.engine.deckB != nil {
				// Volume control now works with the new audio engine
				m.engine.deckB.SetVolume(1.2)
				m.engine.updateGains()
				m.msg = "Deck B: Volume increased"
			}

		// crossfade
		case "f":
			if m.engine != nil {
				m.engine.SetCrossfade(clamp(m.engine.crossfade-0.05, 0.0, 1.0))
			}
		case "h":
			if m.engine != nil {
				m.engine.SetCrossfade(clamp(m.engine.crossfade+0.05, 0.0, 1.0))
			}

		// seek
		case "left":
			if m.engine != nil && m.engine.deckA != nil {
				m.engine.deckA.SeekSeconds(-5)
				m.msg = "Deck A: Rewind 5s"
			}
		case "right":
			if m.engine != nil && m.engine.deckA != nil {
				m.engine.deckA.SeekSeconds(5)
				m.msg = "Deck A: Forward 5s"
			}
		case "shift+left":
			if m.engine != nil && m.engine.deckB != nil {
				m.engine.deckB.SeekSeconds(-5)
				m.msg = "Deck B: Rewind 5s"
			}
		case "shift+right":
			if m.engine != nil && m.engine.deckB != nil {
				m.engine.deckB.SeekSeconds(5)
				m.msg = "Deck B: Forward 5s"
			}
		
		// Fine seeking (1 second increments)
		case "a+left":
			if m.engine != nil && m.engine.deckA != nil {
				m.engine.deckA.SeekSeconds(-1)
				m.msg = "Deck A: Rewind 1s"
			}
		case "a+right":
			if m.engine != nil && m.engine.deckA != nil {
				m.engine.deckA.SeekSeconds(1)
				m.msg = "Deck A: Forward 1s"
			}
		case "k+left":
			if m.engine != nil && m.engine.deckB != nil {
				m.engine.deckB.SeekSeconds(-1)
				m.msg = "Deck B: Rewind 1s"
			}
		case "k+right":
			if m.engine != nil && m.engine.deckB != nil {
				m.engine.deckB.SeekSeconds(1)
				m.msg = "Deck B: Forward 1s"
			}
		
		// Coarse seeking (10 second increments)
		case "shift+a+left":
			if m.engine != nil && m.engine.deckA != nil {
				m.engine.deckA.SeekSeconds(-10)
				m.msg = "Deck A: Rewind 10s"
			}
		case "shift+a+right":
			if m.engine != nil && m.engine.deckA != nil {
				m.engine.deckA.SeekSeconds(10)
				m.msg = "Deck A: Forward 10s"
			}
		case "shift+k+left":
			if m.engine != nil && m.engine.deckB != nil {
				m.engine.deckB.SeekSeconds(-10)
				m.msg = "Deck B: Rewind 10s"
			}
		case "shift+k+right":
			if m.engine != nil && m.engine.deckB != nil {
				m.engine.deckB.SeekSeconds(10)
				m.msg = "Deck B: Forward 10s"
			}
		
		// Jump to start/end
		case "home":
			if m.engine != nil && m.engine.deckA != nil {
				m.engine.deckA.SeekToStart()
				m.msg = "Deck A: Jump to start"
			}
		case "end":
			if m.engine != nil && m.engine.deckA != nil {
				m.engine.deckA.SeekToEnd()
				m.msg = "Deck A: Jump to end"
			}
		case "shift+home":
			if m.engine != nil && m.engine.deckB != nil {
				m.engine.deckB.SeekToStart()
				m.msg = "Deck B: Jump to start"
			}
		case "shift+end":
			if m.engine != nil && m.engine.deckB != nil {
				m.engine.deckB.SeekToEnd()
				m.msg = "Deck B: Jump to end"
			}

		// beatmatch naive
		case "b":
			if m.engine != nil && m.engine.deckA != nil && m.engine.deckB != nil {
				m.msg = "Beatmatching: Limited with current audio backend"
			}

		// EQ presets
		case "F1":
			if m.engine != nil {
				m.engine.ApplyEQPreset("A", "bass")
				m.msg = "Deck A: BASS preset"
			}
		case "F2":
			if m.engine != nil {
				m.engine.ApplyEQPreset("B", "vocals")
				m.msg = "Deck B: VOCALS preset"
			}
		default:
			// ignore other keys
		}
	}
	return m, nil
}

func (m Model) View() string {
	sb := &strings.Builder{}
	// header with DJ ASCII art
sb.WriteString(`
	____      _   _____ _   _ ___ 
	|  _ \    | | |_   _| | | |_ _| Controls: 1/2 Play Toggle	
	| | | |_  | |   | | | | | || | z/x n/m Volume +/-
	| |_| | |_| |   | | | |_| || | f/h Crossfade
	|____/ \___/    |_|  \___/|___| ←/→ Seek A 5s   Shift+←/Shift+→ Seek B 5s
	|____/ \___/    |_|  \___/|___| a+←/→ Fine A 1s   k+←/→ Fine B 1s
	|____/ \___/    |_|  \___/|___| Shift+a+←/→ Coarse A 10s   Shift+k+←/→ Coarse B 10s
	|____/ \___/    |_|  \___/|___| Home/End A   Shift+Home/End B
	|____/ \___/    |_|  \___/|___| F1/F2 EQ presets
`)

	// deck info
	var aPath, bPath string
	var aPlaying, bPlaying bool
	var aPos, aDur, bPos, bDur float64
	if m.engine != nil {
		if m.engine.deckA != nil {
			aPath = m.engine.deckA.path
			aPlaying = m.engine.deckA.playing
			aPos = m.engine.deckA.GetCurrentPosition()
			aDur = m.engine.deckA.GetDuration()
		}
		if m.engine.deckB != nil {
			bPath = m.engine.deckB.path
			bPlaying = m.engine.deckB.playing
			bPos = m.engine.deckB.GetCurrentPosition()
			bDur = m.engine.deckB.GetDuration()
		}
	}

	sb.WriteString(fmt.Sprintf(" Deck A: %s [%s]\n",
		trimPath(aPath, 40), ternary(aPlaying, "PLAYING", "PAUSED")))
	sb.WriteString(fmt.Sprintf("  Position: %s / %s\n", formatTime(aPos), formatTime(aDur)))
	sb.WriteString("  " + renderBars(m.ampA, 40) + "\n\n")

	sb.WriteString(fmt.Sprintf(" Deck B: %s [%s]\n",
		trimPath(bPath, 40), ternary(bPlaying, "PLAYING", "PAUSED")))
	sb.WriteString(fmt.Sprintf("  Position: %s / %s\n", formatTime(bPos), formatTime(bDur)))
	sb.WriteString("  " + renderBars(m.ampB, 40) + "\n\n")

	// crossfader
	if m.engine != nil {
		sb.WriteString(fmt.Sprintf(" Crossfade: [%s] %.2f\n\n", renderCrossfade(m.engine.crossfade), m.engine.crossfade))
	}

	// status line
	if m.msg != "" {
		sb.WriteString(" Status: " + m.msg + "\n\n")
	} else if m.engine != nil && m.engine.lastErr != nil {
		sb.WriteString(" Warning: audio init issues. Check files.\n\n")
	}

	sb.WriteString(" Press q to quit\n")

	return sb.String()
}
	