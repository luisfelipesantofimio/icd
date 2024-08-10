package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

// GetConfig returns clientId and clientSecret from config file. Panics if not found
func GetConfig(path string) Config {
	var config Config
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.WriteFile(path, []byte("{"+"\"clientId\":\"YOUR_CLIENT_ID\","+"\"clientSecret\":\"YOUR_CLIENT_SECRET\""+"}"), 0644)
	}
	configFile, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(configFile, &config)
	if err != nil {
		panic(err)
	}

	if config.ClientID == "" || config.ClientSecret == "" {
		panic("Missing clientId or clientSecret")
	}

	return config
}
