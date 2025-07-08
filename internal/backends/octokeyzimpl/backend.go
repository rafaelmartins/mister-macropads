package octokeyzimpl

import (
	"errors"
	"time"

	"github.com/rafaelmartins/mister-macropads/internal/cleanup"
	"rafaelmartins.com/p/octokeyz"
)

type Backend struct {
	Name   string
	dev    *octokeyz.Device
	mod    []*octokeyz.Modifier
	config config
}

func (b *Backend) GetName() string {
	return b.Name
}

func (b *Backend) Open() error {
	if b.dev != nil {
		return nil
	}

	for {
		dev, err := octokeyz.GetDevice("")
		if err == nil {
			b.dev = dev
			break
		}
		if errors.Is(err, octokeyz.ErrDeviceLocked) {
			return err
		}
		time.Sleep(time.Second)
	}

	return b.dev.Open()
}

func (b *Backend) Close() error {
	if b.dev == nil {
		return errors.New("octokeyz: device not connected")
	}
	return b.dev.Close()
}

func (b *Backend) Listen() {
	if b.dev != nil {
		cleanup.Check(b.dev.Listen(nil))
	}
}
