package merry

import (
	"fmt"
	"os"
	"path/filepath"
)

type Merry struct {
	Running    bool `json:"running"`
	Ready      bool `json:"ready"`
	IsBare     bool `json:"isBare"`
	IsHeadless bool `json:"isHeadless"`

	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`

	Services map[string]string `json:"-"`
}

func Default() Merry {
	merry := Merry{
		Version: "dev",
		Commit:  "none",
		Date:    "unknown",
	}
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("failed to get user's home directory: %v", err))
	}
	if err := merry.provisionResourcesToDatadir(filepath.Join(home, ".merry")); err != nil {
		panic(fmt.Errorf("failed to provision resources: %v", err))
	}
	services, err := getServices(filepath.Join(home, ".merry", "docker-compose.yml"))
	if err != nil {
		panic(fmt.Errorf("failed to get services from docker compose: %v", err))
	}
	merry.Services = services
	return merry
}
