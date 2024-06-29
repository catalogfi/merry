package merry

import (
	"fmt"
	"io"
	"os"
)

func (m *Merry) Logs(service string, outWriter, errWriter io.Writer) error {
	if !m.Running {
		return fmt.Errorf("merry is not running")
	}
	composePath, err := defaultComposePath()
	if err != nil {
		return err
	}
	bashCmd := runDockerCompose(composePath, "logs", service)
	bashCmd.Stdout = os.Stdout
	bashCmd.Stderr = os.Stderr
	return bashCmd.Run()
}
