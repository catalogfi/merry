package main

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/compose-spec/compose-go/loader"
	"github.com/spf13/cobra"
)

//go:embed resources/docker-compose.yml
//go:embed resources/bitcoin.conf
var f embed.FS

type State struct {
	Running    bool `json:"running"`
	Ready      bool `json:"ready"`
	IsBare     bool `json:"isBare"`
	IsHeadless bool `json:"isHeadless"`
}

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	var cmd = &cobra.Command{
		Use: "Merry - catalog localnet",
		Run: func(c *cobra.Command, args []string) {
			c.HelpFunc()(c, args)
		},
		Version:           formatVersion(),
		DisableAutoGenTag: true,
	}

	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("failed to get user's home directory: %v", err))
	}
	state := State{}
	if _, err := os.Stat(filepath.Join(home, ".merry", "merry.config.json")); !os.IsNotExist(err) {
		data, err := os.ReadFile(filepath.Join(home, ".merry", "merry.config.json"))
		if err != nil && err != os.ErrNotExist {
			panic(err)
		}

		if err := json.Unmarshal(data, &state); err != nil {
			panic(err)
		}
	}

	if err := provisionResourcesToDatadir(&state, filepath.Join(home, ".merry")); err != nil {
		panic(fmt.Errorf("failed to provision resources: %v", err))
	}

	cmd.AddCommand(Start(&state), Stop(&state), Faucet(&state), Logs(state), RPC(state), Update(), Version())
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}

var DefaultCompose = "docker-compose.yml"

// Provisioning Nigiri reosurces
func provisionResourcesToDatadir(state *State, datadir string) error {
	if state.Ready {
		return nil
	}

	// create folders in volumes/{bitcoin,elements} for node datadirs
	if err := makeDirectoryIfNotExists(filepath.Join(datadir, "volumes", "bitcoin")); err != nil {
		return err
	}

	// copy docker compose into the Nigiri data directory
	if err := copyFromResourcesToDatadir(
		filepath.Join("resources", DefaultCompose),
		filepath.Join(datadir, DefaultCompose),
	); err != nil {
		return err
	}

	// copy bitcoin.conf into the Nigiri data directory
	if err := copyFromResourcesToDatadir(
		filepath.Join("resources", "bitcoin.conf"),
		filepath.Join(datadir, "volumes", "bitcoin", "bitcoin.conf"),
	); err != nil {
		return err
	}

	state.Ready = true

	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(datadir, "merry.config.json"), data, 0777); err != nil {
		return err
	}
	return nil
}

func makeDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModeDir|0755)
	}
	return nil
}

func copyFromResourcesToDatadir(src string, dest string) error {
	data, err := f.ReadFile(src)
	if err != nil {
		return fmt.Errorf("read embed: %w", err)
	}
	err = os.WriteFile(dest, data, 0777)
	if err != nil {
		return fmt.Errorf("write %s to %s: %w", src, dest, err)
	}
	return nil
}

func GetServices(composeFile string) ([][]string, error) {
	composeBytes, err := os.ReadFile(composeFile)
	if err != nil {
		return nil, err
	}

	parsed, err := loader.ParseYAML(composeBytes)
	if err != nil {
		return nil, err
	}

	if _, ok := parsed["services"]; !ok {
		return nil, errors.New("missing services in compose")
	}

	serviceMap := parsed["services"].(map[string]interface{})

	var services [][]string
	for k, v := range serviceMap {
		m := v.(map[string]interface{})
		i, ok := m["ports"].([]interface{})
		if !ok {
			continue
		}
		for _, j := range i {
			port := j.(string)
			exposedPorts := strings.Split(port, ":")
			endpoint := "localhost:" + exposedPorts[0]
			services = append(services, []string{k, endpoint})
		}
	}

	return services, nil
}

func GetPortsForService(composeFile string, serviceName string) ([]string, error) {
	composeBytes, err := os.ReadFile(composeFile)
	if err != nil {
		return nil, err
	}

	parsed, err := loader.ParseYAML(composeBytes)
	if err != nil {
		return nil, err
	}

	if _, ok := parsed["services"]; !ok {
		return nil, errors.New("missing services in compose")
	}

	serviceMap := parsed["services"].(map[string]interface{})

	var ports []string
	for k, v := range serviceMap {
		if k == serviceName {
			m := v.(map[string]interface{})
			i := m["ports"].([]interface{})
			for _, j := range i {
				port := j.(string)
				ports = append(ports, port)
			}
		}
	}

	return ports, nil
}

func formatVersion() string {
	return fmt.Sprintf(
		"\nVersion: %s\nCommit: %s\nDate: %s",
		version, commit, date,
	)
}
