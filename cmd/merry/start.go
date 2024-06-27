package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/catalogfi/blockchain/localnet"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func Start(state *State) *cobra.Command {
	var isHeadless bool
	var isBare bool
	var cmd = &cobra.Command{
		Use:   "go",
		Short: "Start merry",
		RunE: func(c *cobra.Command, args []string) error {
			state.IsBare = isBare
			state.IsHeadless = isHeadless
			return startMerry(state)
		},
	}
	cmd.Flags().BoolVar(&isHeadless, "headless", false, "do not run UI services")
	cmd.Flags().BoolVar(&isBare, "bare", false, "deploy only blockchains")
	return cmd
}

// runDockerCompose runs docker-compose with the given arguments
func runDockerCompose(composePath string, args ...string) *exec.Cmd {
	var cmd *exec.Cmd
	_, err := exec.LookPath("docker-compose")
	if err != nil {
		cmd = exec.Command("docker", append([]string{"compose", "-f", composePath}, args...)...)
	} else {
		cmd = exec.Command("docker-compose", append([]string{"-f", composePath}, args...)...)
	}
	return cmd
}

func startMerry(state *State) error {
	if state.Running {
		fmt.Println(fmt.Errorf("merry is already running"))
		return fmt.Errorf("merry is already running")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(fmt.Errorf("failed to get user's home directory: %v", err))
		return fmt.Errorf("failed to get user's home directory: %v", err)
	}

	composePath := filepath.Join(home, ".merry", DefaultCompose)

	bashCmd := runDockerCompose(composePath, "up", "-d", "cobi", "esplora", "ethereum-explorer", "arbitrum-explorer")
	if state.IsHeadless {
		bashCmd = runDockerCompose(composePath, "up", "-d", "cobi")
	}
	if state.IsBare {
		bashCmd = runDockerCompose(composePath, "up", "-d", "chopsticks", "ethereum", "arbitrum")
	}
	bashCmd.Stdout = os.Stdout
	bashCmd.Stderr = os.Stderr
	if err := bashCmd.Run(); err != nil {
		fmt.Println("failed to run docker compose command", err)
		return err
	}

	services, err := GetServices(composePath)
	if err != nil {
		fmt.Println("failed to get services", err)
		return err
	}

	fmt.Println()
	fmt.Println("ENDPOINTS")
	for _, nameAndEndpoint := range services {
		name := nameAndEndpoint[0]
		endpoint := nameAndEndpoint[1]

		fmt.Println(name + " " + endpoint)
	}

	state.Running = true
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(home, ".merry", "merry.config.json"), data, 0777); err != nil {
		return err
	}
	retry(func() error { return fundBTC("bcrt1q5428vq2uzwhm3taey9sr9x5vm6tk78ew8pf2xw") })
	retry(func() error { return hasEVMClientStarted() })
	return nil
}

func retry(f func() error) {
	for {
		err := f()
		if err == nil {
			return
		}
		fmt.Printf("failed with %v, retrying after 5 seconds\n", err)
		time.Sleep(5 * time.Second)
	}
}

func hasEVMClientStarted() error {
	bal1, err := localnet.EVMClient().Balance(context.Background(), localnet.WBTC(), common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"), nil)
	if err != nil {
		return err
	}
	bal2, err := localnet.EVMClient().Balance(context.Background(), localnet.ArbitrumWBTC(), common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"), nil)
	if err != nil {
		return err
	}
	bal, ok := new(big.Int).SetString("2100000000000000", 10)
	if !ok {
		return fmt.Errorf("constraint violation")
	}
	if bal1.Cmp(bal) != 0 || bal2.Cmp(bal) != 0 {
		return fmt.Errorf("wbtc tokens have not been deployed yet")
	}
	return nil
}
