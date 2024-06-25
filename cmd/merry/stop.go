package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func Stop(state *State) *cobra.Command {
	var (
		delete bool
	)
	var cmd = &cobra.Command{
		Use:   "stop",
		Short: "stop merry",
		RunE: func(c *cobra.Command, args []string) error {
			return stopMerry(state, delete)
		},
	}
	cmd.Flags().BoolVarP(&delete, "delete", "d", false, "reset the blockchain data")
	return cmd
}

func stopMerry(state *State, isDelete bool) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user's home directory: %v", err)
	}
	composePath := filepath.Join(home, ".merry", DefaultCompose)
	bashCmd := runDockerCompose(composePath, "stop")
	if isDelete {
		bashCmd = runDockerCompose(composePath, "down", "--volumes")
	}
	bashCmd.Stdout = os.Stdout
	bashCmd.Stderr = os.Stderr
	if err := bashCmd.Run(); err != nil {
		return err
	}
	if isDelete {
		fmt.Println("Removing data from volumes...")
		if err := os.RemoveAll(filepath.Join(home, ".merry")); err != nil {
			return err
		}

		if err := provisionResourcesToDatadir(state, filepath.Join(home, ".merry")); err != nil {
			return err
		}
		fmt.Println("Merry has been cleaned up successfully.")
	} else {
		state.Running = false
		state.IsBare = false
		state.IsHeadless = false
		data, err := json.Marshal(state)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(home, ".merry", "merry.config.json"), data, 0777); err != nil {
			return err
		}
	}
	return nil
}
