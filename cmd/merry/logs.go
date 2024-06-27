package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func Logs(state State) *cobra.Command {
	var (
		service string
	)
	var cmd = &cobra.Command{
		Use:   "logs",
		Short: "Check specific service logs",
		RunE: func(c *cobra.Command, args []string) error {
			if !state.Running {
				return fmt.Errorf("merry is not running")
			}
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get user's home directory: %v", err)
			}
			composePath := filepath.Join(home, ".merry", DefaultCompose)

			bashCmd := runDockerCompose(composePath, "logs", service)
			bashCmd.Stdout = os.Stdout
			bashCmd.Stderr = os.Stderr
			if err := bashCmd.Run(); err != nil {
				return err
			}

			return nil
		},
	}
	cmd.Flags().StringVarP(&service, "service", "s", "", "name of the service")
	return cmd
}
