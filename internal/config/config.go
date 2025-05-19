package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	homeDir, err := getConfigFilePath()
	cfg := Config{}
	if err != nil {
		return cfg, err
	}
	full_path := filepath.Join(homeDir, configFileName)
	file, err := os.Open(full_path)
	if err != nil {
		return cfg, err
	}
	defer file.Close()
	data := make([]byte, 500)
	count, bytes_err := file.Read(data)
	if bytes_err != nil {
		fmt.Println("bytes error")
		return cfg, bytes_err
	}
	// We need to resize data[] to exact size to prevent invalid character '\x00' error
	json_err := json.Unmarshal(data[:count], &cfg)
	if json_err != nil {
		fmt.Printf("json error: %s\n", json_err)
		return cfg, json_err
	}
	return cfg, nil

}

func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUserName = username
	cfg.write()
	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir, nil
}

func (cfg *Config) write() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	full_path := filepath.Join(homeDir, configFileName)
	file, err := os.Create(full_path)
	if err != nil {
		return err
	}

	data, json_err := json.Marshal(cfg)
	if json_err != nil {
		return json_err
	}
	file.Write(data)
	return nil
}
