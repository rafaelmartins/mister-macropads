package process

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

var (
	ErrRunning    = errors.New("process already running")
	ErrNotRunning = errors.New("process not running")
)

type Process struct {
	pidFile string
	pid     int
}

func New(pidFile string) (*Process, error) {
	if _, err := os.Stat(pidFile); err == nil {
		pidB, err := os.ReadFile(pidFile)
		if err != nil {
			return nil, err
		}

		pid, err := strconv.Atoi(strings.TrimSpace(string(pidB)))
		if err != nil {
			return nil, err
		}

		proc, err := os.FindProcess(pid)
		if err != nil {
			return nil, err
		}

		if err := proc.Signal(syscall.Signal(0)); err == nil {
			return &Process{
				pidFile: pidFile,
				pid:     pid,
			}, nil
		}
	}
	return &Process{
		pidFile: pidFile,
		pid:     -1,
	}, nil
}

func (p *Process) Write() error {
	if p.pid > 0 {
		return ErrRunning
	}
	return os.WriteFile(p.pidFile, fmt.Appendf(nil, "%d", os.Getpid()), 0666)
}

func (p *Process) Kill() error {
	if p.pid <= 0 {
		return ErrNotRunning
	}

	proc, err := os.FindProcess(p.pid)
	if err != nil {
		return err
	}

	return proc.Signal(syscall.SIGTERM)
}

func (p *Process) IsRunning() error {
	if p.pid > 0 {
		return ErrRunning
	}
	return nil
}
