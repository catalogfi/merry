package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func Update() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "update",
		Short: "Check for updates and pull new docker images",
		RunE: func(c *cobra.Command, args []string) error {
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get user's home directory: %v", err)
			}
			composePath := filepath.Join(home, ".merry", DefaultCompose)

			bashCmd := runDockerCompose(composePath, "pull")
			bashCmd.Stdout = os.Stdout
			bashCmd.Stderr = os.Stderr
			if err := bashCmd.Run(); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}
