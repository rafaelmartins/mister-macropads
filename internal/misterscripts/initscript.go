package misterscripts

import (
	"errors"
	"log"
	"time"

	"github.com/rafaelmartins/mister-macropads/internal/cleanup"
	"github.com/rafaelmartins/mister-macropads/internal/process"
)

// initscriptHandler is `S98{{ projectName }}`. It is called by MiSTer init
// system to start/stop/restart the macro keyboard. It is the same as calling
// onHandler or offHandler as appropriate but without the messages and
// creation/deletion of scripts.
func initscriptHandler(args []string) error {
	log.SetFlags(0)

	action := ""
	if len(args) > 0 {
		action = args[0]
	}

	switch action {
	case "start":
		return startProcess()

	case "stop":
		return stopProcess()

	case "restart":
		err := stopProcess()
		notRunning := errors.Is(err, process.ErrNotRunning)
		if err != nil && !notRunning {
			return err
		}
		if !notRunning {
			time.Sleep(2 * time.Second)
		}
		return startProcess()

	default:
		log.Printf("Usage: %s {start|stop|restart}\n", initd)
		cleanup.Exit(1)
		return nil
	}
}
