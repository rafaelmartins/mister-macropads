package main

import (
	"embed"
	"path/filepath"

	"github.com/rafaelmartins/mister-macropads/internal/backends"
	"github.com/rafaelmartins/mister-macropads/internal/cleanup"
	"github.com/rafaelmartins/mister-macropads/internal/config"
	"github.com/rafaelmartins/mister-macropads/internal/misterscripts"
)

//go:embed octokeyz.ini
var configFS embed.FS

func main() {
	defer cleanup.Cleanup()

	misterscripts.SetMainApp(func(projectName string, configDir string, args []string) error {
		backend, err := backends.Get(projectName)
		if err != nil {
			return err
		}
		cleanup.Register(backend)

		if err := config.Load(filepath.Join(configDir, projectName+".ini"), projectName, backend); err != nil {
			return err
		}

		if err := backend.Open(); err != nil {
			return err
		}
		cleanup.Register(backend)

		if err := backend.ScreenRender(); err != nil {
			return err
		}

		if err := backend.KeypadSetup(); err != nil {
			return err
		}

		backend.Listen()
		return nil
	}, configFS)

	cleanup.Check(misterscripts.Dispatch())
}
