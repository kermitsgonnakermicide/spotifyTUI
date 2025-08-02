package main

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type Player struct {
	stream    beep.StreamSeekCloser
	format    beep.Format
	ctrl      *beep.Ctrl
	resample  *beep.Resampler
	reverse   bool
	speed     float64
	looping   bool
	loopStart int
	loopEnd   int
	totalLen  int
}

func NewPlayer() *Player {
	return &Player{
		speed:   1.0,
		looping: false,
	}
}

func (p *Player) Load(file string) error {
	if p.stream != nil {
		p.stream.Close()
	}
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return err
	}
	p.format = format
	p.stream = streamer
	p.ctrl = &beep.Ctrl{Streamer: streamer, Paused: false}
	p.resample = beep.ResampleRatio(4, 1.0, p.ctrl)
	p.totalLen = streamer.Len()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	go speaker.Play(p)

	return nil
}

func (p *Player) Stream(samples [][2]float64) (n int, ok bool) {
	if p.reverse {
		// Reverse: manually read backwards
		if p.stream.Position() <= 0 {
			return 0, false
		}
		p.stream.Seek(p.stream.Position() - len(samples))
	}
	n, ok = p.resample.Stream(samples)

	// Loop check
	if p.looping && p.stream.Position() >= p.loopEnd {
		p.stream.Seek(p.loopStart)
	}
	return
}

func (p *Player) Err() error { return nil }

func (p *Player) PlayPause() {
	if p.ctrl != nil {
		p.ctrl.Paused = !p.ctrl.Paused
	}
}

func (p *Player) SetSpeed(ratio float64) {
	if p.resample != nil {
		p.resample.SetRatio(ratio)
		p.speed = ratio
	}
}

func (p *Player) ToggleReverse() {
	p.reverse = !p.reverse
}

func (p *Player) SetLoopStart() {
	if p.stream != nil {
		p.loopStart = p.stream.Position()
	}
}

func (p *Player) SetLoopEnd() {
	if p.stream != nil {
		p.loopEnd = p.stream.Position()
		p.looping = true
	}
}

func (p *Player) ClearLoop() {
	p.looping = false
	p.loopStart = 0
	p.loopEnd = 0
}

func (p *Player) CurrentTime() time.Duration {
	if p.stream != nil {
		return time.Duration(float64(p.stream.Position())/float64(p.format.SampleRate)) * time.Second
	}
	return 0
}

func (p *Player) TotalTime() time.Duration {
	if p.stream != nil && p.totalLen > 0 {
		return time.Duration(float64(p.totalLen)/float64(p.format.SampleRate)) * time.Second
	}
	return 0
}
