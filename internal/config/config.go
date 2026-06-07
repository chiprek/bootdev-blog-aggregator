package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// read the config file from home directory
func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	f, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	var cfg Config
	err = json.NewDecoder(f).Decode(&cfg)
	return cfg, err
}

// sets username
func (cfg Config) SetUser(name string) error {
	cfg.CurrentUserName = name
	return write(cfg)
}

// write json data to config
func write(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, jsonData, 0644)
}

// get file path to config
func getConfigFilePath() (string, error) {
	base, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", base, configFileName), nil
}
