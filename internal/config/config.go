package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, configFileName), nil
}

func Read() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	configFile, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}
	data := Config{}
	err = json.Unmarshal(configFile, &data)
	if err != nil {
		return Config{}, err
	}
	return data, nil
}

func (c *Config) SetUser(username string) error {
	c.Current_user_name = username
	err := write(*c)
	if err != nil {
		return err
	}
	return nil
}

func write(config Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	os.WriteFile(filePath, data, 0600)
	return nil
}
