package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/catalogfi/merry"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	state := merry.Default()
	state.Commit = commit
	state.Version = version
	state.Date = date
	var cmd = &cobra.Command{
		Use: "merry - catalog localnet",
		Run: func(c *cobra.Command, args []string) {
			c.HelpFunc()(c, args)
		},
		Version:           state.FormatVersion(),
		DisableAutoGenTag: true,
	}
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("failed to get user's home directory: %v", err))
	}
	if _, err := os.Stat(filepath.Join(home, ".merry", "merry.config.json")); !os.IsNotExist(err) {
		data, err := os.ReadFile(filepath.Join(home, ".merry", "merry.config.json"))
		if err != nil && err != os.ErrNotExist {
			panic(err)
		}

		if err := json.Unmarshal(data, &state); err != nil {
			panic(err)
		}
	}
	cmd.AddCommand(Start(&state), Stop(&state), Faucet(&state), Replace(&state), Logs(state), RPC(state), Update(), Version(state))
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}

func Faucet(merry *merry.Merry) *cobra.Command {
	var (
		to string
	)
	var cmd = &cobra.Command{
		Use:   "faucet",
		Short: "Generate and send supported assets to the given address",
		RunE: func(c *cobra.Command, args []string) error {
			return merry.Fund(to)
		},
	}
	cmd.Flags().StringVar(&to, "to", "", "user should pass the address they needs to be funded")
	return cmd
}

func Logs(merry merry.Merry) *cobra.Command {
	var (
		service string
	)
	var cmd = &cobra.Command{
		Use:   "logs",
		Short: "Check specific service logs",
		RunE: func(c *cobra.Command, args []string) error {
			return merry.Logs(service, os.Stdout, os.Stderr)
		},
	}
	cmd.Flags().StringVarP(&service, "service", "s", "", "name of the service")
	return cmd
}

func Replace(m *merry.Merry) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "replace",
		Short: "Replace a specific service with a local one",
		Run: func(c *cobra.Command, args []string) {
			c.HelpFunc()(c, args)
		},
	}
	cmd.AddCommand(Orderbook(m), COBI(m), EVM(m))
	return cmd
}

func Orderbook(m *merry.Merry) *cobra.Command {
	var path string
	var cmd = &cobra.Command{
		Use:   "orderbook",
		Short: "Replace orderbook service",
		RunE: func(c *cobra.Command, args []string) error {
			return m.Replace(path, "ghcr.io/catalogfi/orderbook:latest")
		},
	}
	cmd.Flags().StringVarP(&path, "path", "p", ".", "docker path")
	return cmd
}

func COBI(m *merry.Merry) *cobra.Command {
	var path string
	var cmd = &cobra.Command{
		Use:   "cobi",
		Short: "Replace cobi service",
		RunE: func(c *cobra.Command, args []string) error {
			return m.Replace(path, "ghcr.io/catalogfi/cobi:latest")
		},
	}
	cmd.Flags().StringVarP(&path, "path", "p", ".", "docker path")
	return cmd
}

func EVM(m *merry.Merry) *cobra.Command {
	var path string
	var cmd = &cobra.Command{
		Use:   "evm",
		Short: "Replace evm chain services",
		RunE: func(c *cobra.Command, args []string) error {
			return m.Replace(path, "ghcr.io/catalogfi/garden_sol:latest")
		},
	}
	cmd.Flags().StringVarP(&path, "path", "p", ".", "docker path")
	return cmd
}

func RPC(m merry.Merry) *cobra.Command {
	var (
		named     bool
		generate  int64
		rpcwallet string
	)
	var cmd = &cobra.Command{
		Use:   "rpc",
		Short: "Invoke bitcoin-cli commands",
		RunE: func(c *cobra.Command, args []string) error {
			return m.Relay(generate, named, rpcwallet, args...)
		},
	}
	cmd.Flags().BoolVar(&named, "named", false, "use named arguments")
	cmd.Flags().Int64VarP(&generate, "generate", "g", 0, "generate block")
	cmd.Flags().StringVarP(&rpcwallet, "rpcwallet", "w", "", "rpcwallet to be used for node JSONRPC commands")
	return cmd
}

func Start(m *merry.Merry) *cobra.Command {
	var isHeadless bool
	var isBare bool
	var cmd = &cobra.Command{
		Use:   "go",
		Short: "Start merry",
		RunE: func(c *cobra.Command, args []string) error {
			m.IsBare = isBare
			m.IsHeadless = isHeadless
			return m.Start()
		},
	}
	cmd.Flags().BoolVar(&isHeadless, "headless", false, "do not run UI services")
	cmd.Flags().BoolVar(&isBare, "bare", false, "deploy only blockchains")
	return cmd
}

func Stop(merry *merry.Merry) *cobra.Command {
	var (
		delete bool
	)
	var cmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop merry",
		RunE: func(c *cobra.Command, args []string) error {
			return merry.Stop(delete)
		},
	}
	cmd.Flags().BoolVarP(&delete, "delete", "d", false, "reset the blockchain data")
	return cmd
}

func Update() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "update",
		Short: "Check for updates and pull new docker images",
		RunE: func(c *cobra.Command, args []string) error {
			return merry.Update()
		},
	}
	return cmd
}

func Version(m merry.Merry) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "version",
		Short: "Get the merry version",
		RunE: func(c *cobra.Command, args []string) error {
			fmt.Println("merry CLI version")
			fmt.Println(m.FormatVersion())
			return nil
		},
	}
	return cmd
}
