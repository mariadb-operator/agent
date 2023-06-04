package mariadbd

import (
	"errors"
	"fmt"
	"syscall"
	"time"

	"github.com/mitchellh/go-ps"
)

const (
	mariadbdProcessName = "mariadbd"
	reloadSysCall       = syscall.SIGKILL
)

var (
	errProcessNotFound = fmt.Errorf("process '%s' not found", mariadbdProcessName)
)

func Reload() error {
	processes, err := ps.Processes()
	if err != nil {
		return fmt.Errorf("error getting processes: %v", err)
	}
	for _, p := range processes {
		if p.Executable() == mariadbdProcessName {
			if err := syscall.Kill(p.Pid(), reloadSysCall); err != nil {
				return fmt.Errorf("error sending kill signal to process '%s' with pid %d: %v", mariadbdProcessName, p.Pid(), err)
			}
			return nil
		}
	}
	return errProcessNotFound
}

type ReloadOptions struct {
	Retries   int
	WaitRetry time.Duration
}

func ReloadWithOptions(opts *ReloadOptions) error {
	for i := 0; i < opts.Retries; i++ {
		err := Reload()
		if err == nil {
			return nil
		}
		if errors.Is(err, errProcessNotFound) {
			time.Sleep(opts.WaitRetry)
			continue
		}
		return err
	}
	return fmt.Errorf("maximum retries (%d) reached attempting to reload '%s' process", opts.Retries, mariadbdProcessName)
}
