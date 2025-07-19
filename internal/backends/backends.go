package backends

import (
	"errors"
	"strings"

	"gopkg.in/ini.v1"
	"rafaelmartins.com/p/mister-macropads/internal/backends/octokeyzimpl"
)

var reg = []struct {
	name string
	new  func(name string) Backend
}{
	{
		name: "octokeyz",
		new:  func(name string) Backend { return &octokeyzimpl.Backend{Name: name} },
	},
}

type Backend interface {
	GetName() string
	Open() error
	Close() error
	Listen()

	GetModel() (string, error)

	SetConfigScreenSection(section *ini.Section) error
	SetConfigKeypadSection(section *ini.Section) error

	ScreenRender() error

	KeypadSetup() error
}

func Get(name string) (Backend, error) {
	for _, be := range reg {
		if be.name == name {
			return be.new(be.name), nil
		}
	}
	return nil, errors.New("backend not found")
}

func List(prefix string) []string {
	rv := []string{}
	for _, be := range reg {
		if strings.HasPrefix(be.name, prefix) {
			rv = append(rv, be.name)
		}
	}
	return rv
}
