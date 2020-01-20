package godooj

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"os"
	"path"
)

type OdooConfig struct {
	ServerAddress string
	Database      string
	Login         string
	Password      string
}

func ClientConnect(serverName string) (*Client, error) {
	config, err := LoadOdooConfig(serverName)
	if err != nil {
		return nil, err
	}
	return config.CreateOdooClient()
}

func (conf *OdooConfig) CreateOdooClient() (*Client, error) {
	return Connect(
		conf.ServerAddress,
		conf.Database,
		conf.Login,
		conf.Password)
}

func fileExists(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func lastFileExists(fileNames []string) string {
	result := ""
	for _, name := range fileNames {
		if fileExists(name) {
			result = name
		}
	}
	return result
}

func LoadOdooConfig(serverName string) (*OdooConfig, error) {
	dirs := []string{"odoo.toml"}
	homeDir, err := os.UserHomeDir()
	if err == nil {
		dirs = append(dirs, path.Join(homeDir, ".odoo.toml"))
	}
	configFile := lastFileExists(dirs)
	if configFile == "" {
		return nil, fmt.Errorf("Could not find odoo.toml from %v", dirs)
	}
	config, err := toml.LoadFile(configFile)
	if err != nil {
		return nil, err
	}

	serverSection := "servers." + serverName
	serverTree := config.Get(serverSection).(*toml.Tree)
	if serverTree == nil {
		return nil, fmt.Errorf("Server configuration not found for %s", serverName)
	}

	for _, key := range []string{"server", "database", "user", "password"} {
		if !serverTree.Has(key) {
			return nil, fmt.Errorf("No `%s` configuration for %s", key, serverName)
		}
	}
	return &OdooConfig{
		ServerAddress: serverTree.Get("server").(string),
		Database:      serverTree.Get("database").(string),
		Login:         serverTree.Get("user").(string),
		Password:      serverTree.Get("password").(string),
	}, nil
}
