package misterscripts

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/rafaelmartins/mister-macropads/internal/process"
)

// mainHandler is the main process run in background. It does all the heavy
// work to make the macro keyboard work, by calling some backend to do it.
func mainHandler(args []string) error {
	logfp, err := os.OpenFile(filepath.Join("/tmp", projectName+".log"), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
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

	log.Printf("Starting %s ...", projectName)
	if app == nil {
		return fmt.Errorf("no app set for %s", projectName)
	}
	return app(projectName, args)
}
