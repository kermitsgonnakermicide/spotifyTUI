package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func renderBars(level float64, width int) string {
	if level < 0 {
		level = 0
	}
	if level > 1 {
		level = 1
	}
	filled := int(level * float64(width))
	if filled > width {
		filled = width
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}

func renderCrossfade(v float64) string {
	w := 20
	f := int(v * float64(w))
	if f < 0 {
		f = 0
	}
	if f > w {
		f = w
	}
	return strings.Repeat("◼", f) + strings.Repeat("◻", w-f)
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func trimPath(p string, max int) string {
	b := filepath.Base(p)
	if len(b) > max {
		return "..." + b[len(b)-max+3:]
	}
	return b
}

func ternary(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}

// formatTime converts seconds to MM:SS format
func formatTime(seconds float64) string {
	if seconds < 0 {
		seconds = 0
	}
	
	minutes := int(seconds) / 60
	secs := int(seconds) % 60
	
	return fmt.Sprintf("%02d:%02d", minutes, secs)
}
