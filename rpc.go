package merry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/logrusorgru/aurora"
)

func (m *Merry) Relay(generate int64, named bool, rpcwallet string, args ...string) error {
	if !m.Running {
		return fmt.Errorf("merry is not running")
	}

	rpcArgs := []string{"exec", "bitcoin", "bitcoin-cli", "-rpcwallet=" + rpcwallet}
	if generate > 0 {
		rpcArgs = append(rpcArgs, "-generate", fmt.Sprint(generate))
	}
	if named {
		rpcArgs = append(rpcArgs, "-named")
	}
	cmdArgs := append(rpcArgs, args...)

	bashCmd := exec.Command("docker", cmdArgs...)

	// Create a pipe for the output of the "docker exec" command
	r, w := io.Pipe()
	bashCmd.Stdout = w
	bashCmd.Stderr = os.Stderr

	// Start a goroutine to run the "docker exec" command
	go func() {
		if err := bashCmd.Run(); err != nil {
			w.CloseWithError(err)
		} else {
			w.Close()
		}
	}()

	// Read the output of the "docker exec" command from the pipe
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	output := buf.Bytes()

	// Use the json.Unmarshal function to parse the output of the
	// "docker exec" command and check if it is a valid JSON object
	var v interface{}
	if err := json.Unmarshal(output, &v); err == nil {
		// Use the json.Marshal function to convert the parsed JSON object
		// to a byte slice
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return err
		}

		// Use the bytes.Buffer type to create a buffer that we can use
		// to write the indented JSON string to
		var buf bytes.Buffer

		// Use the json.Indent function to add indentation to the JSON byte slice
		// in the same way as if you were using the "jq" command
		if err := json.Indent(&buf, jsonBytes, "", "    "); err != nil {
			return err
		}

		// Split the indented JSON string into individual lines
		lines := strings.Split(buf.String(), "\n")

		// Loop through each line in the indented JSON string
		for _, line := range lines {
			// Check if the line starts with a "{" or a "["
			if strings.HasPrefix(line, "{") || strings.HasPrefix(line, "[") {
				// If the line starts with a "{" or a "[", it is the start of a JSON object
				// or array, so print it without any color
				fmt.Println(line)
			} else if strings.Contains(line, ":") {
				// If the line contains a ":", it is a key-value pair, so split the line
				// into the key and value parts and add the desired colors to each part
				parts := strings.SplitN(line, ":", 2)
				key := parts[0]
				value := parts[1]
				// Use the AnsiColorCode function from the "github.com/logrusorgru/aurora"
				// package to create ANSI escape codes for the "key" and "value" colors
				keyColor := aurora.BrightBlue(key)
				valueColor := aurora.BrightCyan(value)

				fmt.Printf("%s: %s\n", keyColor.String(), valueColor.String())
			} else {
				// If the line does not start with a "{" or a "[" and does not contain a ":",
				// it is not a JSON object or array and does not contain a key-value pair, so
				// print it without any color
				fmt.Println(line)
			}
		}
	} else {
		fmt.Println(string(output))
	}
	return nil

}

func (m *Merry) Proxy(req interface{}, service, method string) (Response, error) {
	var response Response

	data, err := json.Marshal(req)
	if err != nil {
		return response, err
	}
	request := Request{
		Version: "2.0",
		ID:      rand.Int(),
		Method:  method,
		Params:  data,
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(request); err != nil {
		return response, err
	}

	composePath, err := defaultComposePath()
	if err != nil {
		return response, err
	}
	endpoints, err := getServices(composePath)
	if err != nil {
		return response, err
	}

	resp, err := http.Post(endpoints[service], "application/json", buf)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, err
	}
	return response, err
}
