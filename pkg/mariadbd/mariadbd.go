package mariadbd

import (
	"fmt"
	"syscall"

	"github.com/mitchellh/go-ps"
)

var (
	mariadbdProcessName = "mariadbd"
	reloadSysCall       = syscall.SIGKILL
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

	return fmt.Errorf("process '%s' not found", mariadbdProcessName)
}
