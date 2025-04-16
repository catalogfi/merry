package merry

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/compose-spec/compose-go/loader"
)

//go:embed resources/docker-compose.yml
//go:embed resources/bitcoin.conf
//go:embed resources/config/*
var f embed.FS

var DefaultComposePath string

func defaultComposePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user's home directory: %v", err)
	}
	return filepath.Join(home, ".merry", "docker-compose.yml"), nil
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

// getServices extracts the services from the docker compose file
func getServices(composePath string) (map[string]string, error) {
	composeBytes, err := os.ReadFile(composePath)
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

	services := map[string]string{}
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
			services[k] = endpoint
		}
	}

	return services, nil
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

// Provisioning Nigiri reosurces
func (m *Merry) provisionResourcesToDatadir(datadir string) error {
	if m.Ready {
		return nil
	}

	// create folders in volumes/{bitcoin,elements} for node datadirs
	if err := makeDirectoryIfNotExists(filepath.Join(datadir, "volumes", "bitcoin")); err != nil {
		return err
	}

	// copy docker compose into the Nigiri data directory
	if err := copyFromResourcesToDatadir(
		filepath.Join("resources", "docker-compose.yml"),
		filepath.Join(datadir, "docker-compose.yml"),
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

	// if err := makeDirectoryIfNotExists(filepath.Join(datadir, "starknet")); err != nil {
	// 	return err
	// }

	// // copy dump.json into the starknet data directory
	// if err := copyFromResourcesToDatadir(
	// 	filepath.Join("resources", "starknet", "dump.json"),
	// 	filepath.Join(datadir, "starknet", "dump.json"),
	// ); err != nil {
	// 	return err
	// }

	// Entire config directory
	if err := copyDirectoryFromResourcesToDatadir("resources/config", filepath.Join(datadir, "config")); err != nil {
		return err
	}

	m.Ready = true
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(datadir, "merry.config.json"), data, 0777); err != nil {
		return err
	}
	return nil
}

func copyDirectoryFromResourcesToDatadir(srcDir, destDir string) error {
	entries, err := f.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("failed to read embedded directory %s: %w", srcDir, err)
	}

	// if not exists, make it
	if err := makeDirectoryIfNotExists(destDir); err != nil {
		return fmt.Errorf("failed to create destination directory %s: %w", destDir, err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(destDir, entry.Name())

		// if directory then copy dir if not copy file (recursion)
		if entry.IsDir() {
			if err := copyDirectoryFromResourcesToDatadir(srcPath, destPath); err != nil {
				return err
			}
		} else {
			if err := copyFromResourcesToDatadir(srcPath, destPath); err != nil {
				return err
			}
		}
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

func getPortsForService(composeFile string, serviceName string) ([]string, error) {
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
