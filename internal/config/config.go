package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rafaelmartins/mister-macropads/internal/backends"
	"gopkg.in/ini.v1"
)

func Load(f string, project string, backend backends.Backend) error {
	cfg, err := ini.Load(f)
	if err != nil {
		return err
	}

	if len(backends.List(project)) > 1 {
		if !cfg.HasSection("device") {
			return errors.New("config: configuration file missing [device] section")
		}

		device := cfg.Section("device")
		if !device.HasKey("model") {
			return errors.New("config: configuration file missing `model` key from [device] section")
		}

		if model := device.Key("model").String(); model != backend.GetName() {
			return fmt.Errorf("config: configuration file is set for a different keypad model: %s != %s", model, backend.GetName())
		}
	}

	if cfg.HasSection("screen") {
		if err := backend.SetConfigScreenSection(cfg.Section("screen")); err != nil {
			return err
		}
	}

	if cfg.HasSection("keypad") {
		if err := backend.SetConfigKeypadSection(cfg.Section("keypad")); err != nil {
			return err
		}
	}

	for _, sect := range cfg.SectionStrings() {
		if strings.HasPrefix(sect, "screen:") {
			if err := backend.SetConfigScreenSection(cfg.Section(sect)); err != nil {
				return err
			}
		}
		if strings.HasPrefix(sect, "keypad:") {
			if err := backend.SetConfigKeypadSection(cfg.Section(sect)); err != nil {
				return err
			}
		}
	}
	return nil
}
