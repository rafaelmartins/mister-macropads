package misterscripts

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/rafaelmartins/mister-macropads/internal/cleanup"
)

var (
	lexe string
	exe  string

	projectName string
	pidFile     string

	initd string
	off   string

	app func(projectName string, args []string) error

	isMister bool
)

func init() {
	var err error
	lexe, err = filepath.Abs(os.Args[0])
	if err != nil {
		panic(err)
	}

	exe, err = filepath.EvalSymlinks(lexe)
	if err != nil {
		panic(err)
	}

	projectName = strings.TrimSuffix(filepath.Base(exe), "_on.sh")
	pidFile = filepath.Join("/run", projectName+".pid")

	off = filepath.Join(filepath.Dir(exe), projectName+"_off.sh")
	initd = "S98" + projectName

	if strings.HasSuffix(filepath.Base(exe), "_on.sh") {
		initd = filepath.Join("/etc", "init.d", initd)
		isMister = true
	} else {
		initd = filepath.Join(filepath.Dir(exe), initd)
	}
}

// Dispatch detects how the binary was called and dispatches execution to a
// proper handler function.
func Dispatch() error {
	if len(os.Args) >= 2 && os.Args[1] == "-v" {
		if bi, ok := debug.ReadBuildInfo(); ok {
			fmt.Println(bi.Main.Version)
			return nil
		}
		fmt.Println("UNKNOWN")
		cleanup.Exit(1)
		return nil
	}

	if os.Getuid() != 0 {
		return errors.New("must run as root")
	}

	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "-init":
			return initscriptHandler(os.Args[2:])

		case "-main":
			return mainHandler(os.Args[2:])
		}
	}

	if filepath.Base(lexe) == filepath.Base(off) {
		return offHandler(os.Args[1:])
	}
	return onHandler(os.Args[1:])
}

func SetApp(fn func(projectName string, args []string) error) {
	app = fn
}

func remountRW() error {
	cmd := exec.Command("bash", "-c", "if mount | grep \"on / .*[(,]ro[,$]\"; then mount / -o remount,rw; fi")
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}
