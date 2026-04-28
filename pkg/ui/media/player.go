package media

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type MPVPlayer struct {
	Path string

	state    State
	position time.Duration
	duration time.Duration
	volume   float32
	err      error

	cmd       *exec.Cmd
	socket    string
	lastStart time.Time
}

func NewMPVPlayer() *MPVPlayer {
	return &MPVPlayer{
		state:  StateIdle,
		volume: 0.8,
	}
}

func (p *MPVPlayer) Load(path string) error {
	_ = p.Stop()

	p.Path = path
	p.state = StatePaused
	p.position = 0

	duration, err := probeDuration(context.Background(), path)
	if err == nil {
		p.duration = duration
	} else {
		p.duration = 0
	}
	p.err = nil

	return nil
}

func (p *MPVPlayer) Play() error {
	if p.Path == "" {
		return fmt.Errorf("no media loaded")
	}

	if p.state == StatePlaying {
		return nil
	}

	if p.cmd == nil || p.cmd.Process == nil {
		if err := p.startMPV(); err != nil {
			p.err = err
			p.state = StateError
			return err
		}
	} else {
		if err := p.command("set_property", "pause", false); err != nil {
			p.err = err
			p.state = StateError
			return err
		}
	}

	p.state = StatePlaying
	p.err = nil
	return nil
}

func (p *MPVPlayer) Pause() error {
	if p.cmd != nil && p.cmd.Process != nil {
		if err := p.command("set_property", "pause", true); err != nil {
			p.err = err
			p.state = StateError
			return err
		}
	}
	p.state = StatePaused
	return nil
}

func (p *MPVPlayer) Stop() error {
	if p.cmd != nil && p.cmd.Process != nil {
		_ = p.command("quit")
		_, _ = p.cmd.Process.Wait()
	}
	if p.socket != "" {
		_ = os.Remove(p.socket)
	}
	p.cmd = nil
	p.socket = ""
	p.state = StateStopped
	p.position = 0
	p.err = nil
	return nil
}

func (p *MPVPlayer) Seek(pos time.Duration) error {
	if p.cmd != nil && p.cmd.Process != nil {
		if err := p.command("seek", pos.Seconds(), "absolute"); err != nil {
			p.err = err
			p.state = StateError
			return err
		}
	}
	p.position = pos
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
	if p.cmd != nil && p.cmd.Process != nil {
		if err := p.command("set_property", "volume", v*100); err != nil {
			p.err = err
			p.state = StateError
			return err
		}
	}
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

func (p *MPVPlayer) startMPV() error {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}

	if err := os.MkdirAll(filepath.Join(cacheDir, "bare"), 0o755); err != nil {
		return err
	}

	p.lastStart = time.Now()
	p.socket = filepath.Join(cacheDir, "bare", fmt.Sprintf("mpv-%d.sock", p.lastStart.UnixNano()))
	args := []string{
		"--input-ipc-server=" + p.socket,
		"--force-window=yes",
		"--pause=no",
		"--idle=no",
		p.Path,
	}
	cmd := exec.Command("mpv", args...)
	if err := cmd.Start(); err != nil {
		return err
	}
	p.cmd = cmd

	for range 30 {
		if _, err := os.Stat(p.socket); err == nil {
			return p.command("set_property", "volume", p.volume*100)
		}
		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("mpv IPC socket did not become ready")
}

func (p *MPVPlayer) command(args ...any) error {
	if p.socket == "" {
		return fmt.Errorf("mpv socket not initialized")
	}

	conn, err := net.Dial("unix", p.socket)
	if err != nil {
		return err
	}
	defer conn.Close()

	payload := map[string]any{
		"command": args,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	data = append(data, '\n')

	_, err = conn.Write(data)
	return err
}
