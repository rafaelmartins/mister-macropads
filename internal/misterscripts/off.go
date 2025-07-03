package misterscripts

import (
	"errors"
	"io/fs"
	"log"
	"os"

	"github.com/rafaelmartins/mister-macropads/internal/process"
)

func stopProcess() error {
	proc, err := process.New(pidFile)
	if err != nil {
		return err
	}

	return proc.Kill()
}

// offHandler is `{{ projectName }}_off.sh`. It is called from MiSTer Scripts
// menu to stop the macro keyboard and disable/remove the init script. It
// stops a main process running on background, and exits with a message.
func offHandler(_ []string) error {
	log.SetFlags(0)

	if isMister {
		if err := remountRW(); err != nil {
			return err
		}
	}

	if err := stopProcess(); err != nil {
		return err
	}

	if err := os.Remove(initd); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	log.Printf("%s is off and inactive at startup.\n", projectName)
	return nil
}
