package misterscripts

import (
	"log"
	"os"
	"path/filepath"

	"rafaelmartins.com/p/mister-macropads/internal/process"
	"rafaelmartins.com/p/mister-macropads/internal/services"
)

// mainHandler is the main process run in background. It does all the heavy
// work to make the macro keyboard work, by calling some backend to do it.
func mainHandler(args []string) error {
	logfp, err := os.OpenFile(filepath.Join("/tmp", projectName+".log"), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	// FIXME: cleanup logfp somehow
	log.SetOutput(logfp)

	proc, err := process.New(pidFile)
	if err != nil {
		return err
	}
	if err := proc.Write(); err != nil {
		return err
	}

	if err := services.Start(projectName); err != nil {
		return err
	}

	configDir := filepath.Join("/media", "fat")
	if !isMister {
		cfg, err := os.Getwd()
		if err != nil {
			return err
		}
		configDir = cfg
	}

	log.Printf("Starting %s %s ...", projectName, projectVersion)
	return app(projectName, configDir, args)
}
