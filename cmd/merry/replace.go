package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func Replace(state *State) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "replace",
		Short: "Replace a specific service with a local one",
		RunE: func(c *cobra.Command, args []string) error {
			fmt.Println("replace any service")
			fmt.Println(formatVersion())
			return nil
		},
	}
	cmd.AddCommand(Orderbook(state), COBI(state), EVM(state))
	return cmd
}

func Orderbook(state *State) *cobra.Command {
	var path string
	var cmd = &cobra.Command{
		Use:   "orderbook",
		Short: "Replace orderbook service",
		RunE: func(c *cobra.Command, args []string) error {
			return replaceImage(state, path, "ghcr.io/catalogfi/orderbook:latest")
		},
	}
	cmd.Flags().StringVarP(&path, "path", "p", ".", "docker path")
	return cmd
}

func COBI(state *State) *cobra.Command {
	var path string
	var cmd = &cobra.Command{
		Use:   "cobi",
		Short: "Replace cobi service",
		RunE: func(c *cobra.Command, args []string) error {
			return replaceImage(state, path, "ghcr.io/catalogfi/cobi:latest")
		},
	}
	cmd.Flags().StringVarP(&path, "path", "p", ".", "docker path")
	return cmd
}

func EVM(state *State) *cobra.Command {
	var path string
	var cmd = &cobra.Command{
		Use:   "evm",
		Short: "Replace evm chain services",
		RunE: func(c *cobra.Command, args []string) error {
			return replaceImage(state, path, "ghcr.io/catalogfi/garden_sol:latest")
		},
	}
	cmd.Flags().StringVarP(&path, "path", "p", ".", "docker path")
	return cmd
}

func replaceImage(state *State, path, image string) error {
	if err := stopMerry(state, false); err != nil {
		return err
	}
	bashCmd := exec.Command("docker", "build", path, "-t", image)
	bashCmd.Stdout = os.Stdout
	bashCmd.Stderr = os.Stderr
	if err := bashCmd.Run(); err != nil {
		// log the error and move on
		fmt.Println("failed to build the image", err)
	}
	return startMerry(state)
}
