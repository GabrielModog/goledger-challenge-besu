package config

import (
	"errors"
	"fmt"
	"os"
)

type Config struct {
	RPCURL           string
	ContractAddress  string
	PrivateKey       string
	ABIPath          string
	ConnectionString string
	ServerPort       string
}

func LoadConfig() (*Config, error) {
	required := []struct {
		name string
		val  string
	}{
		{name: "BESU_NODE_URL", val: os.Getenv("BESU_NODE_URL")},
		{name: "CONTRACT_ADDRESS", val: os.Getenv("CONTRACT_ADDRESS")},
		{name: "SIGNER_PRIVATE_KEY", val: os.Getenv("SIGNER_PRIVATE_KEY")},
		{name: "CONTRACT_ABI_PATH", val: os.Getenv("CONTRACT_ABI_PATH")},
		{name: "CONNECTION_STRING", val: os.Getenv("CONNECTION_STRING")},
	}

	var missing []string
	for _, r := range required {
		if r.val == "" {
			missing = append(missing, r.name)
		}
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required env vars: %v", missing)
	}

	return &Config{
		RPCURL:           os.Getenv("BESU_NODE_URL"),
		ContractAddress:  os.Getenv("CONTRACT_ADDRESS"),
		PrivateKey:       os.Getenv("SIGNER_PRIVATE_KEY"),
		ABIPath:          os.Getenv("CONTRACT_ABI_PATH"),
		ConnectionString: os.Getenv("CONNECTION_STRING"),
		ServerPort:       getEnvOrDefault("SERVER_PORT", "8080"),
	}, nil
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func (c *Config) Validate() error {
	if c.RPCURL == "" {
		return errors.New("BESU_NODE_URL cannot be empty")
	}
	if c.ContractAddress == "" {
		return errors.New("CONTRACT_ADDRESS cannot be empty")
	}
	if c.PrivateKey == "" {
		return errors.New("SIGNER_PRIVATE_KEY cannot be empty")
	}
	return nil
}
