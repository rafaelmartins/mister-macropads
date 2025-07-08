package services

import (
	"errors"
	"net"

	"github.com/rafaelmartins/mister-macropads/internal/cleanup"
	"github.com/rafaelmartins/mister-macropads/internal/services/inotify"
	"github.com/rafaelmartins/mister-macropads/internal/services/ipaddr"
	"github.com/rafaelmartins/mister-macropads/internal/services/uinput"
)

var (
	svcInotify *inotify.Monitor
	svcIpAddr  *ipaddr.Monitor
	svcUInput  *uinput.Device

	ErrStarted    = errors.New("services: already started")
	ErrNotStarted = errors.New("services: not started")
)

func Start(projectName string) error {
	if svcInotify != nil && svcIpAddr != nil {
		return ErrStarted
	}

	in, err := inotify.New()
	if err != nil {
		return err
	}
	cleanup.Register(in)

	ip := ipaddr.New()
	cleanup.Register(ip)

	keys := []uinput.Key{}
	for _, k := range uinput.KeyMap {
		keys = append(keys, k)
	}
	ui, err := uinput.NewDevice(projectName, keys)
	if err != nil {
		return err
	}
	cleanup.Register(ui)

	go func() {
		cleanup.Check(in.Listen())
	}()

	go func() {
		cleanup.Check(ip.Listen())
	}()

	svcInotify = in
	svcIpAddr = ip
	svcUInput = ui

	return nil
}

func AddInotifyWatch(f string, ev inotify.InotifyMask, handler func(f string, ev inotify.InotifyMask) error) error {
	if svcInotify == nil {
		return ErrNotStarted
	}
	return svcInotify.AddWatch(f, ev, handler)
}

func AddIpAddrWatch(itf string, handler func(itf string, ip net.IP) error) error {
	if svcIpAddr == nil {
		return ErrNotStarted
	}
	return svcIpAddr.AddWatch(itf, handler)
}

func UInputPress(kc ...uinput.Key) error {
	if svcUInput == nil {
		return ErrNotStarted
	}
	return svcUInput.Press(kc...)
}

func UInputRelease(kc ...uinput.Key) error {
	if svcUInput == nil {
		return ErrNotStarted
	}
	return svcUInput.Release(kc...)
}
