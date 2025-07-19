package octokeyzimpl

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/ini.v1"
	"rafaelmartins.com/p/mister-macropads/internal/services/uinput"
	"rafaelmartins.com/p/octokeyz"
)

var (
	reOctokeyzLine = regexp.MustCompile(`^screen:line([0-9]+)$`)
	reOctokeyzKey  = regexp.MustCompile(`^keypad:key([0-9]+)$`)
)

type config struct {
	screen map[octokeyz.DisplayLine]*configScreenLine
	keypad map[octokeyz.ButtonID]*configKeypadKey
}

type configScreenLine struct {
	action       string
	prefix       string
	align        octokeyz.DisplayLineAlign
	stringAction *configScreenLineString
	ipAddrAction *configScreenLineIpAddr
}

type configScreenLineString struct {
	str string
}

type configScreenLineIpAddr struct {
	itf string
}

type configKeypadKey struct {
	action         string
	holdKeysAction *configKeypadKeyHoldKeys
}

type configKeypadKeyHoldKeys struct {
	dflt []uinput.Key
	mod  []uinput.Key
}

func (b *Backend) SetConfigScreenSection(section *ini.Section) error {
	name := section.Name()
	if name == "screen" {
		return nil
	}

	m := reOctokeyzLine.FindStringSubmatch(name)
	if len(m) != 2 {
		return fmt.Errorf("octokeyz: unknown screen line configuration section: %s", name)
	}

	l, err := strconv.ParseUint(m[1], 10, 8)
	if err != nil {
		return err
	}
	line := octokeyz.DisplayLine(l)

	prefix := ""
	if section.HasKey("prefix") {
		prefix = section.Key("prefix").String()
	}

	align := octokeyz.DisplayLineAlignLeft
	if section.HasKey("align") {
		switch strings.ToLower(section.Key("align").String()) {
		case "center":
			align = octokeyz.DisplayLineAlignCenter

		case "right":
			align = octokeyz.DisplayLineAlignRight
		}
	}

	if !section.HasKey("action") {
		return fmt.Errorf("octokeyz: configuration section [%s] missing `action` key", name)
	}
	action := strings.ToLower(section.Key("action").String())

	screenLine := &configScreenLine{
		action: action,
		prefix: prefix,
		align:  align,
	}

	switch action {
	case "string":
		if !section.HasKey("string") {
			return fmt.Errorf("octokeyz: configuration section [%s] missing `string` key", name)
		}
		screenLine.stringAction = &configScreenLineString{
			str: section.Key("string").String(),
		}

	case "ipaddr":
		if !section.HasKey("interface") {
			return fmt.Errorf("octokeyz: configuration section [%s] missing `interface` key", name)
		}
		screenLine.ipAddrAction = &configScreenLineIpAddr{
			itf: section.Key("interface").String(),
		}
	}

	if b.config.screen == nil {
		b.config.screen = map[octokeyz.DisplayLine]*configScreenLine{}
	}
	b.config.screen[line] = screenLine
	return nil
}

func stringToKeys(s string) ([]uinput.Key, error) {
	rv := []uinput.Key{}
	for p := range strings.SplitSeq(s, " ") {
		k := strings.TrimSpace(p)
		if k == "" {
			continue
		}

		key, ok := uinput.KeyMap[k]
		if !ok {
			return nil, fmt.Errorf("octokeyz: unknown key: %s", k)
		}
		rv = append(rv, key)
	}
	return rv, nil
}

func (b *Backend) SetConfigKeypadSection(section *ini.Section) error {
	name := section.Name()
	if name == "keypad" {
		return nil
	}

	m := reOctokeyzKey.FindStringSubmatch(name)
	if len(m) != 2 {
		return fmt.Errorf("octokeyz: unknown keypad line configuration section: %s", name)
	}

	k, err := strconv.ParseUint(m[1], 10, 8)
	if err != nil {
		return err
	}
	key := octokeyz.ButtonID(k)

	if !section.HasKey("action") {
		return fmt.Errorf("octokeyz: configuration section [%s] missing `action` key", name)
	}
	action := strings.ToLower(section.Key("action").String())

	keypadKey := &configKeypadKey{
		action: action,
	}

	switch action {
	case "hold_keys":
		if !section.HasKey("default") {
			return fmt.Errorf("octokeyz: configuration section [%s] missing `default` key", name)
		}
		dflt, err := stringToKeys(section.Key("default").String())
		if err != nil {
			return err
		}
		keypadKey.holdKeysAction = &configKeypadKeyHoldKeys{
			dflt: dflt,
		}
		if section.HasKey("modifier") {
			mod, err := stringToKeys(section.Key("modifier").String())
			if err != nil {
				return err
			}
			keypadKey.holdKeysAction.mod = mod
		}
	}

	if b.config.keypad == nil {
		b.config.keypad = map[octokeyz.ButtonID]*configKeypadKey{}
	}
	b.config.keypad[key] = keypadKey
	return nil
}
