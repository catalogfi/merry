package merry

import (
	"encoding/json"
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
	merry := Merry{}
	if err := merry.Load(); err == nil {
		return merry
	}
	merry.Version = "dev"
	merry.Commit = "none"
	merry.Date = "unknown"
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

func (m *Merry) Load() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(filepath.Join(home, ".merry", "merry.config.json"))
	if err != nil {
		return err
	}
	return json.Unmarshal(data, m)
}

func (m *Merry) Save() error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(home, ".merry", "merry.config.json"), data, 0777); err != nil {
		return err
	}
	return nil
}
