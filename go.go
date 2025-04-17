package merry

import (
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

	bashCmd := runDockerCompose(composePath, "up", "-d", "esplora", "ethereum-explorer", "arbitrum-explorer", "nginx", "garden-evm-watcher", "garden-db", "quote", "bit-ponder", "cobiv2", "relayer", "bit-indexer", "authenticator", "orderbookV2")
	if m.IsHeadless && m.IsBare {
		bashCmd = runDockerCompose(composePath, "up", "-d", "chopsticks", "ethereum", "arbitrum", "cosigner")
	} else if m.IsHeadless {
		bashCmd = runDockerCompose(composePath, "up", "-d", "chopsticks", "cosigner")
	} else if m.IsBare {
		bashCmd = runDockerCompose(composePath, "up", "-d", "chopsticks", "ethereum-explorer", "arbitrum-explorer", "cosigner")

	}
	bashCmd.Stdout = os.Stdout
	bashCmd.Stderr = os.Stderr
	if err := bashCmd.Run(); err != nil {
		fmt.Println("failed to run docker compose command", err)
		return err
	}

	m.Running = true
	if err := m.Save(); err != nil {
		return err
	}

	// Funding
	fundAddresses := []string{
		"bcrt1qgyf47wrtnr9gsr06gn62ft6m4lzylcnllrf9cf", // cobi btc address
		"0x70997970c51812dc3a010c7d01b50e0d17dc79c8",   // cobi evm address
	}

	for _, addr := range fundAddresses {
		retry(func() error { return m.Fund(addr) })
	}

	retry(func() error {
		// try establishing connection with the ethereum clients
		_, err := localnet.EVMClient()
		return err
	})

	// display endpoints
	if err := m.Status(); err != nil {
		return err
	}

	return nil
}
