package main

import (
	_ "embed"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"os"
	"strings"
	"time"

	"github.com/rafaelmartins/mister-macropads/internal/cleanup"
	"github.com/rafaelmartins/mister-macropads/internal/inotify"
	"github.com/rafaelmartins/mister-macropads/internal/ipaddr"
	"github.com/rafaelmartins/mister-macropads/internal/misterscripts"
	"github.com/rafaelmartins/mister-macropads/internal/uinput"
	"github.com/rafaelmartins/mister-macropads/internal/vkbd"
	"rafaelmartins.com/p/octokeyz"
)

func waitForOctokeyz(sn string) (*octokeyz.Device, error) {
	tick := time.NewTicker(time.Second)
	tim := time.NewTimer(2 * time.Minute)

	for {
		dev, err := octokeyz.GetDevice(sn)
		if err == nil {
			return dev, nil
		}
		if errors.Is(err, octokeyz.ErrDeviceLocked) {
			return nil, err
		}

		select {
		case <-tick.C:
			continue
		case <-tim.C:
			return nil, err
		}
	}
}

func updateCore(dev *octokeyz.Device) error {
	v, err := os.ReadFile("/tmp/CORENAME")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
	}
	return dev.DisplayLine(octokeyz.DisplayLine4, fmt.Sprintf("Core: %s", strings.TrimSpace(string(v))), octokeyz.DisplayLineAlignLeft)
}

func app(projectName string, args []string) error {
	dev, err := waitForOctokeyz("")
	if err != nil {
		return err
	}

	if err := dev.Open(); err != nil {
		return err
	}
	cleanup.Register(dev)

	kbd, err := vkbd.New(dev, octokeyz.BUTTON_5, map[octokeyz.ButtonID]vkbd.ButtonMapping{
		octokeyz.BUTTON_1: {
			Normal: []uinput.Key{uinput.KEY_F12},
			Mod:    []uinput.Key{uinput.KEY_LEFTALT, uinput.KEY_F12},
		},
		octokeyz.BUTTON_2: {
			Normal: []uinput.Key{uinput.KEY_ESC},
		},
		octokeyz.BUTTON_3: {
			Normal: []uinput.Key{uinput.KEY_UP},
		},
		octokeyz.BUTTON_4: {
			Normal: []uinput.Key{uinput.KEY_ENTER},
		},
		octokeyz.BUTTON_6: {
			Normal: []uinput.Key{uinput.KEY_LEFT},
			Mod:    []uinput.Key{uinput.KEY_LEFTCTRL, uinput.KEY_LEFTALT, uinput.KEY_RIGHTALT},
		},
		octokeyz.BUTTON_7: {
			Normal: []uinput.Key{uinput.KEY_DOWN},
			Mod:    []uinput.Key{uinput.KEY_LEFTSHIFT, uinput.KEY_LEFTCTRL, uinput.KEY_LEFTALT, uinput.KEY_RIGHTALT},
		},
		octokeyz.BUTTON_8: {
			Normal: []uinput.Key{uinput.KEY_RIGHT},
		},
	})
	if err != nil {
		return err
	}
	cleanup.Register(kbd)

	dev.AddHandler(octokeyz.BUTTON_5, func(b *octokeyz.Button) error {
		return dev.Led(octokeyz.LedFlash)
	})

	if err := dev.DisplayLine(octokeyz.DisplayLine1, "MiSTer FPGA", octokeyz.DisplayLineAlignCenter); err != nil {
		return err
	}

	in, err := inotify.New()
	if err != nil {
		return err
	}
	cleanup.Register(in)

	if err := in.AddWatch("/tmp/CORENAME", inotify.IN_CLOSE_WRITE); err != nil {
		return err
	}
	if err := updateCore(dev); err != nil {
		return err
	}

	go func() {
		cleanup.Check(in.Listen(func(f string, ev inotify.InotifyEvent) error {
			return updateCore(dev)
		}))
	}()

	ip, err := ipaddr.NewMonitor("eth0", "wlan0")
	if err != nil {
		return err
	}
	cleanup.Register(ip)

	go func() {
		cleanup.Check(ip.Run(func(itf string, ip net.IP) error {
			line := octokeyz.DisplayLine6
			if itf == "eth0" {
				line = octokeyz.DisplayLine7
			}
			if ip == nil {
				return dev.DisplayClearLine(line)
			}
			return dev.DisplayLine(line, fmt.Sprintf("%s: %s", strings.ToUpper(itf[:len(itf)-1]), ip), octokeyz.DisplayLineAlignLeft)
		}))
	}()

	return dev.Listen(nil)
}

func main() {
	misterscripts.SetApp(app)
	cleanup.Check(misterscripts.Dispatch())
}
