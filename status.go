package merry

import (
	"bytes"
	"fmt"
	"os"
	"strings"

)

func displayContainerStatus(composepath string) error {
	// Run docker ps with custom format
	cmd := runDockerCompose(composepath, "ps", "--format", "\"{{.Names}}|{{.Ports}}|{{.Status}}\"")

	// Capture the output
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	// Print table header
	fmt.Printf("%-20s | %-30s | %-20s\n", "SERVICE NAME", "PORT", "STATUS")
	fmt.Println(strings.Repeat("-", 75))

	// Process each line which is already in our desired format
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) < 3 {
			continue
		}

		serviceName := parts[0]

		// Extract the port from the port string
		port := "N/A"
		if parts[1] != "" {
			portInfo := parts[1]
			// Extract first port only for cleaner display
			if strings.Contains(portInfo, "->") {
				portsParts := strings.Split(portInfo, ",")[0]
				if serviceName == "bitcoin" {
					portsParts = strings.Split(portInfo, ",")[2]
				}
				portParts := strings.Split(portsParts, "->")[0]
				if strings.Contains(portParts, ":") {
					port = strings.Split(portParts, ":")[1]
					port = strings.Split(port, "-")[0]
				} else {
					port = portParts
				}
			} else {
				port = portInfo
			}
		}

		status := parts[2]

		fmt.Printf("%-20s | %-30s | %-20s\n", serviceName, port, status)
	}

	return nil
}

func (m *Merry) Status() error {
	if !m.Running {
		return fmt.Errorf("merry is not running")
	}
	fmt.Println("ACTIVE ENDPOINTS")
	composePath, err := defaultComposePath()
	if err != nil {
		return err
	}
	return displayContainerStatus(composePath)
}
