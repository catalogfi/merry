package merry

import (
	"fmt"
	"os"
)

func Update() error {
	composePath, err := defaultComposePath()
	if err != nil {
		return fmt.Errorf("failed to get user's home directory: %v", err)
	}
	bashCmd := runDockerCompose(composePath, "pull")
	bashCmd.Stdout = os.Stdout
	bashCmd.Stderr = os.Stderr
	if err := bashCmd.Run(); err != nil {
		return err
	}
	return nil
}
