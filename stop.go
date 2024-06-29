package merry

import (
	"fmt"
	"os"
	"path/filepath"
)

func (m *Merry) Stop(isDelete bool) error {
	composePath, err := defaultComposePath()
	if err != nil {
		return err
	}
	bashCmd := runDockerCompose(composePath, "stop")
	if isDelete {
		bashCmd = runDockerCompose(composePath, "down", "--volumes")
	}
	bashCmd.Stdout = os.Stdout
	bashCmd.Stderr = os.Stderr
	if err := bashCmd.Run(); err != nil {
		return err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user's home directory: %v", err)
	}
	if isDelete {
		fmt.Println("Removing data from volumes...")
		if err := os.RemoveAll(filepath.Join(home, ".merry")); err != nil {
			return err
		}

		if err := m.provisionResourcesToDatadir(filepath.Join(home, ".merry")); err != nil {
			return err
		}
		fmt.Println("Merry has been cleaned up successfully.")
	} else {
		m.Running = false
		m.IsBare = false
		m.IsHeadless = false
		m.Save()
	}
	return nil
}
