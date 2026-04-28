package media

import "time"

type MPVPlayer struct {
	Path string

	state    State
	position time.Duration
	duration time.Duration
	volume   float32
	err      error
}

func NewMPVPlayer() *MPVPlayer {
	return &MPVPlayer{
		state:  StateIdle,
		volume: 0.8,
	}
}

func (p *MPVPlayer) Load(path string) error {
	p.Path = path
	p.state = StatePaused

	// later: send mpv IPC command:
	// {"command":["loadfile", path, "replace"]}

	return nil
}

func (p *MPVPlayer) Play() error {
	p.state = StatePlaying

	// later:
	// {"command":["set_property","pause",false]}

	return nil
}

func (p *MPVPlayer) Pause() error {
	p.state = StatePaused

	// later:
	// {"command":["set_property","pause",true]}

	return nil
}

func (p *MPVPlayer) Stop() error {
	p.state = StateStopped

	// later:
	// {"command":["stop"]}

	return nil
}

func (p *MPVPlayer) Seek(pos time.Duration) error {
	p.position = pos

	// later:
	// {"command":["seek", seconds, "absolute"]}

	return nil
}

func (p *MPVPlayer) SetVolume(v float32) error {
	if v < 0 {
		v = 0
	}
	if v > 1 {
		v = 1
	}

	p.volume = v

	// later:
	// {"command":["set_property","volume", v*100]}

	return nil
}

func (p *MPVPlayer) Position() time.Duration {
	return p.position
}

func (p *MPVPlayer) Duration() time.Duration {
	return p.duration
}

func (p *MPVPlayer) State() State {
	return p.state
}

func (p *MPVPlayer) Error() error {
	return p.err
}
