package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Version() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "version",
		Short: "check for updates and pull new docker images",
		RunE: func(c *cobra.Command, args []string) error {
			fmt.Println("merry CLI version")
			fmt.Println(formatVersion())
			return nil
		},
	}
	return cmd
}
