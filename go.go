package merry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/catalogfi/blockchain/localnet"
)

func (m *Merry) Start() error {
	if m.Running {
		fmt.Println(fmt.Errorf("merry is already running"))
		return fmt.Errorf("merry is already running")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(fmt.Errorf("failed to get user's home directory: %v", err))
		return fmt.Errorf("failed to get user's home directory: %v", err)
	}
	composePath := filepath.Join(home, ".merry", "docker-compose.yml")

	bashCmd := runDockerCompose(composePath, "up", "-d", "cobi", "esplora", "ethereum-explorer", "arbitrum-explorer")
	if m.IsHeadless && m.IsBare {
		bashCmd = runDockerCompose(composePath, "up", "-d", "chopsticks", "ethereum", "arbitrum")
	} else if m.IsHeadless {
		bashCmd = runDockerCompose(composePath, "up", "-d", "chopsticks", "cobi")
	} else if m.IsBare {
		bashCmd = runDockerCompose(composePath, "up", "-d", "chopsticks", "ethereum-explorer", "arbitrum-explorer")
	}
	bashCmd.Stdout = os.Stdout
	bashCmd.Stderr = os.Stderr
	if err := bashCmd.Run(); err != nil {
		fmt.Println("failed to run docker compose command", err)
		return err
	}

	fmt.Println()
	fmt.Println("ENDPOINTS")
	for name, endpoint := range m.Services {
		if m.IsBare {
			if name == "cobi" || name == "redis" || name == "orderbook" || name == "postgres" {
				continue
			}
		}
		if m.IsHeadless {
			if name == "esplora" || name == "ethereum-explorer" || name == "arbitrum-explorer" {
				continue
			}
		}
		fmt.Println(name + " " + endpoint)
	}

	m.Running = true
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(home, ".merry", "merry.config.json"), data, 0777); err != nil {
		return err
	}

	retry(func() error {
		_, err := localnet.FundBTC("bcrt1q5428vq2uzwhm3taey9sr9x5vm6tk78ew8pf2xw")
		return err
	})
	retry(func() error {
		// try establishing connection with the ethereum clients
		_, err := localnet.EVMClient()
		return err
	})
	return nil
}
