package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type ServerConfig struct {
	LogRequests bool `toml:"LogRequests"`
}

type DatabaseConfig struct {
	Host     string `toml:"Host"`
	Port     int    `toml:"Port"`
	User     string `toml:"User"`
	Password string `toml:"Password"`
	DBName   string `toml:"DBName"`
}

type GameConfig struct {
	IP   string `toml:"IP"`
	Port int    `toml:"Port"`
}

type FreeroamConfig struct {
	Port int `toml:"Port"`
}

type XMPPConfig struct {
	Port int `toml:"Port"`
}

type Config struct {
	Server   ServerConfig   `toml:"Server"`
	Database DatabaseConfig `toml:"Database"`
	Game     GameConfig     `toml:"Game"`
	Freeroam FreeroamConfig `toml:"Freeroam"`
	XMPP     XMPPConfig     `toml:"XMPP"`
}

var AppConfig Config

func LoadConfig() error {
	configFilePath := "config/config.toml"
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return fmt.Errorf("config file not found at %s", configFilePath)
	}
	_, err := toml.DecodeFile(configFilePath, &AppConfig)
	if err != nil {
		return fmt.Errorf("failed to load config file: %s: %w", configFilePath, err)
	}
	return nil
}
