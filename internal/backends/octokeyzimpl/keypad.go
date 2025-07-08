package octokeyzimpl

import (
	"errors"

	"github.com/rafaelmartins/mister-macropads/internal/services"
	"rafaelmartins.com/p/octokeyz"
)

func (b *Backend) KeypadSetup() error {
	if b.dev == nil {
		return errors.New("octokeyz: device not connected")
	}

	for btn, data := range b.config.keypad {
		switch data.action {
		case "modifier":
			if err := b.dev.AddHandler(btn, func(bt *octokeyz.Button) error {
				return b.dev.Led(octokeyz.LedFlash)
			}); err != nil {
				return err
			}

			m := &octokeyz.Modifier{}
			if err := b.dev.AddHandler(btn, m.Handler); err != nil {
				return err
			}
			b.mod = append(b.mod, m)

		case "hold_keys":
			if err := b.dev.AddHandler(btn, func(bt *octokeyz.Button) error {
				modPressed := false
				for _, p := range b.mod {
					if p.Pressed() {
						modPressed = true
						break
					}
				}

				c := data.holdKeysAction.dflt
				if modPressed {
					c = data.holdKeysAction.mod
				}

				if err := services.UInputPress(c...); err != nil {
					return err
				}
				bt.WaitForRelease()
				return services.UInputRelease(c...)
			}); err != nil {
				return err
			}
		}
	}
	return nil
}
