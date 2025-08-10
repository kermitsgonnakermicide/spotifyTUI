package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// ensure music files exist; warn rather than crash
	if _, err := os.Stat("music/trackA.mp3"); os.IsNotExist(err) {
		log.Println("Warning: put music/trackA.mp3 in place")
	}
	if _, err := os.Stat("music/trackB.mp3"); os.IsNotExist(err) {
		log.Println("Warning: put music/trackB.mp3 in place")
	}

	// create DJ engine
	engine, err := NewEngine("music/trackA.mp3", "music/trackB.mp3")
	if err != nil {
		log.Println("Engine init warning:", err)
		// continue — UI will show error
	}

	model := NewModel(engine)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal("Program error:", err)
	}
}
