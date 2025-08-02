package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	list    list.Model
	player  *Player
	playing bool
	file    string
}

type item string
func (i item) Title() string       { return string(i) }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return string(i) }

func initialModel() model {
	files, _ := filepath.Glob("*.mp3")
	items := []list.Item{}
	for _, f := range files {
		items = append(items, item(f))
	}

	l := list.New(items, list.NewDefaultDelegate(), 40, 10)
	l.Title = "🎶 Live Remix Player"

	return model{
		list:   l,
		player: NewPlayer(),
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			selected := m.list.SelectedItem().(item)
			m.file = string(selected)
			m.player.Load(m.file)
			m.playing = true

		case "space":
			m.player.PlayPause()

		case "+":
			m.player.SetSpeed(min(2.0, m.player.speed+0.1))

		case "-":
			m.player.SetSpeed(max(0.5, m.player.speed-0.1))

		case "l":
			m.player.SetLoopStart()

		case "e":
			m.player.SetLoopEnd()

		case "c":
			m.player.ClearLoop()

		case "r":
			m.player.ToggleReverse()
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	loopInfo := "OFF"
	if m.player.looping {
		loopInfo = fmt.Sprintf("%v → %v",
			time.Duration(float64(m.player.loopStart)/float64(m.player.format.SampleRate))*time.Second,
			time.Duration(float64(m.player.loopEnd)/float64(m.player.format.SampleRate))*time.Second)
	}
	progress := ""
	if m.playing {
		ct := m.player.CurrentTime()
		tt := m.player.TotalTime()
		pct := int(float64(ct) / float64(tt) * 20)
		progress = fmt.Sprintf("[%s%s] %v / %v",
			strings.Repeat("█", pct),
			strings.Repeat("-", 20-pct),
			ct.Truncate(time.Second), tt.Truncate(time.Second))
	}
	return fmt.Sprintf(`%s

Speed: %.1fx   Loop: %s   Reverse: %v
%s
`, m.list.View(), m.player.speed, loopInfo, m.player.reverse, progress)
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
