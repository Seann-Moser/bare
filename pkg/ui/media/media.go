package media

import "time"

type Kind string

const (
	KindImage Kind = "image"
	KindAudio Kind = "audio"
	KindVideo Kind = "video"
)

type State string

const (
	StateIdle    State = "idle"
	StatePlaying State = "playing"
	StatePaused  State = "paused"
	StateStopped State = "stopped"
	StateError   State = "error"
)

type Player interface {
	Load(path string) error
	Play() error
	Pause() error
	Stop() error
	Seek(pos time.Duration) error
	SetVolume(v float32) error

	Position() time.Duration
	Duration() time.Duration
	State() State
	Error() error
}
