package misterscripts

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"rafaelmartins.com/p/mister-macropads/internal/process"
)

func startProcess() error {
	oproc, err := process.New(pidFile)
	if err != nil {
		return err
	}
	if err := oproc.IsRunning(); err != nil {
		return err
	}

	syscall.Setsid()
	cmd := exec.Command(exe, "-main")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}

// onHandler is `{{ projectName }}_on.sh`. It is called from MiSTer Scripts menu to
// start the macro keyboard and enable/create the init script. It starts a new
// main process on background, and exits with a message.
func onHandler(_ []string) error {
	log.SetFlags(0)

	if isMister {
		if err := remountRW(); err != nil {
			return err
		}
	}

	if err := startProcess(); err != nil {
		return err
	}

	if _, err := os.Stat(off); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		if err := os.Symlink(filepath.Base(exe), off); err != nil {
			return err
		}
	}

	if _, err := os.Stat(initd); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		if err := os.WriteFile(initd, fmt.Appendf(nil, "#!/bin/bash\nexec %q -init \"${@}\"\n", exe), 0777); err != nil {
			return err
		}
	}

	log.Printf("%s is on and active at startup.\n", projectName)
	return nil
}
