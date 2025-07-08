package dataproviders

import (
	"errors"
	"io/fs"
	"os"
	"strings"

	"github.com/rafaelmartins/mister-macropads/internal/services"
	"github.com/rafaelmartins/mister-macropads/internal/services/inotify"
)

var CoreName CoreNameType

type CoreNameType struct {
	initialized bool
	value       string
}

func (c *CoreNameType) read() error {
	v, err := os.ReadFile("/tmp/CORENAME")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			c.value = "UNKNOWN"
			return nil
		}
		return err
	}

	c.value = strings.TrimSpace(string(v))
	return nil
}

func (c *CoreNameType) Get(backend Backend) (string, error) {
	if !c.initialized {
		if err := services.AddInotifyWatch("/tmp/CORENAME", inotify.IN_CLOSE_WRITE, func(f string, ev inotify.InotifyMask) error {
			if err := c.read(); err != nil {
				return err
			}
			return backend.ScreenRender()
		}); err != nil {
			return "", err
		}
		if err := c.read(); err != nil {
			return "", err
		}
		c.initialized = true
	}
	return c.value, nil
}
