package merry

import (
	"fmt"
	"os"
	"os/exec"
)

func (m *Merry) Replace(path, image string) error {
	if err := m.Stop(false); err != nil {
		return err
	}
	bashCmd := exec.Command("docker", "build", path, "-t", image)
	bashCmd.Stdout = os.Stdout
	bashCmd.Stderr = os.Stderr
	if err := bashCmd.Run(); err != nil {
		// log the error and move on
		fmt.Println("failed to build the image", err)
	}
	return m.Start()
}
