package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

func Faucet(state *State) *cobra.Command {
	var (
		to string
	)
	var cmd = &cobra.Command{
		Use:   "faucet",
		Short: "Generate and send bitcoin to given address",
		RunE: func(c *cobra.Command, args []string) error {
			if !state.Running {
				return fmt.Errorf("merry is not running")
			}

			payload, err := json.Marshal(map[string]string{
				"address": to,
			})
			if err != nil {
				return fmt.Errorf("failed to marshal address: %v", err)
			}

			res, err := http.Post("http://127.0.0.1:3000/faucet", "application/json", bytes.NewBuffer(payload))
			if err != nil {
				return fmt.Errorf("failed to get funds from faucet: %v", err)
			}
			data, err := io.ReadAll(res.Body)
			if err != nil {
				return err
			}
			if res.StatusCode != http.StatusOK {
				return errors.New(string(data))
			}
			var dat map[string]string
			if err := json.Unmarshal([]byte(data), &dat); err != nil {
				return errors.New("internal error, please try again")
			}
			if dat["txId"] == "" {
				return errors.New("not successful")
			}
			fmt.Println("Successfully submitted at http://localhost:5000/tx/" + dat["txId"])
			return nil
		},
	}
	cmd.Flags().StringVar(&to, "to", "", "user should pass the address they needs to be funded")
	return cmd
}
