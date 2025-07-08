package octokeyzimpl

import (
	"errors"
	"fmt"

	"github.com/rafaelmartins/mister-macropads/internal/dataproviders"
)

func (b *Backend) ScreenRender() error {
	if b.dev == nil {
		return errors.New("octokeyz: device not connected")
	}

	for line, data := range b.config.screen {
		switch data.action {
		case "string":
			if err := b.dev.DisplayLine(line, data.prefix+data.stringAction.str, data.align); err != nil {
				return err
			}

		case "corename":
			core, err := dataproviders.CoreName.Get(b)
			if err != nil {
				return err
			}
			if err := b.dev.DisplayLine(line, data.prefix+core, data.align); err != nil {
				return err
			}

		case "ipaddr":
			ip, err := dataproviders.IpAddr.Get(b, data.ipAddrAction.itf)
			if err != nil {
				return err
			}
			if err := b.dev.DisplayLine(line, data.prefix+ip, data.align); err != nil {
				return err
			}

		default:
			return fmt.Errorf("octokeyz: screen: line %d: invalid action: %s", line, data.action)
		}
	}
	return nil
}
