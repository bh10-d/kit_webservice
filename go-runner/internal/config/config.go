// package config

// import (
// 	"fmt"
// 	"net"
// 	"os"
// 	"strings"
// 	// "log"
// 	"github.com/joho/godotenv"
// )

// // Config holds all configuration values
// type Config struct {
// 	Runner RunnerConfig `json:"runner"`
// }

// // RunnerConfig holds runner-specific configuration
// type RunnerConfig struct {
// 	ID           string   `json:"id"`
// 	Name         string   `json:"name"`
// 	ServerURL    string   `json:"server_url"`
// 	MessageQueue string   `json:"message_queue"`
// 	ScriptsPath  string   `json:"scripts_path"`
// 	Tags         []string `json:"tags"`
// }

// // Load loads configuration from environment variables and .env file
// func Load() (*Config, error) {
// 	// Try to load .env file (don't fail if it doesn't exist)
// 	_ = godotenv.Load()
// 	// err := godotenv.Load()
// 	// if err != nil {
// 	// 	log.Println("⚠️ Không tìm thấy file .env, dùng env của hệ thống")
// 	// }

// 	config := &Config{}

// 	// Load runner configuration
// 	if err := config.loadRunnerConfig(); err != nil {
// 		return nil, fmt.Errorf("failed to load runner config: %w", err)
// 	}

// 	return config, nil
// }

// func (c *Config) loadRunnerConfig() error {
// 	c.Runner = RunnerConfig{
// 		ID:           getEnvWithDefault("RUNNER_ID", "default-runner"),
// 		Name:         getEnvWithDefault("RUNNER_NAME", getHostname()),
// 		ServerURL:    getEnvWithDefault("SERVER_URL", "http://localhost:8080"),
// 		MessageQueue: getEnvWithDefault("MESSAGE_QUEUE", "localhost"),
// 		ScriptsPath:  getEnvWithDefault("SCRIPTS_PATH", "./scripts"),
// 		Tags:         parseStringSlice(getEnvWithDefault("RUNNER_TAGS", "")),
// 	}

// 	return nil
// }

// // Helper functions
// func getEnvWithDefault(key, defaultValue string) string {
// 	if value := os.Getenv(key); value != "" {
// 		return value
// 	}
// 	return defaultValue
// }

// func parseStringSlice(value string) []string {
// 	if value == "" {
// 		return []string{}
// 	}
	
// 	parts := strings.Split(value, ",")
// 	result := make([]string, 0, len(parts))
	
// 	for _, part := range parts {
// 		if trimmed := strings.TrimSpace(part); trimmed != "" {
// 			result = append(result, trimmed)
// 		}
// 	}
	
// 	return result
// }

// func getHostname() string {
// 	hostname, err := os.Hostname()
// 	if err != nil {
// 		return "unknown"
// 	}
// 	return hostname
// }

// func getLocalIP() string {
// 	addrs, err := net.InterfaceAddrs()
// 	if err != nil {
// 		return "127.0.0.1"
// 	}

// 	for _, addr := range addrs {
// 		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
// 			if ipnet.IP.To4() != nil {
// 				return ipnet.IP.String()
// 			}
// 		}
// 	}

// 	return "127.0.0.1"
// }







package config

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration values
type Config struct {
	Runner RunnerConfig `json:"runner"`
}

// RunnerConfig holds runner-specific configuration
type RunnerConfig struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	ServerURL    string   `json:"server_url"`
	MessageQueue string   `json:"message_queue"`
	ScriptsPath  string   `json:"scripts_path"`
	Tags         []string `json:"tags"`
}

// Load loads configuration from environment variables and .env file
func Load() (*Config, error) {
	// Try to load .env file (don't fail if it doesn't exist)
	_ = godotenv.Load()

	config := &Config{}

	// Load runner configuration
	if err := config.loadRunnerConfig(); err != nil {
		return nil, fmt.Errorf("failed to load runner config: %w", err)
	}

	return config, nil
}

func (c *Config) loadRunnerConfig() error {
	c.Runner = RunnerConfig{
		ID:           getEnvWithDefault("KEY_ID", ""),
		Name:         getEnvWithDefault("RUNNER_NAME", getHostname()),
		ServerURL:    getEnvWithDefault("SERVER_URL", "http://localhost:5000"),
		MessageQueue: getEnvWithDefault("MESSAGE_QUEUE", "localhost"),
		ScriptsPath:  getEnvWithDefault("SCRIPTS_PATH", "../../scripts"),
		Tags:         parseStringSlice(getEnvWithDefault("RUNNER_TAGS", "")),
	}

	return nil
}

// Helper functions
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseStringSlice(value string) []string {
	if value == "" {
		return []string{}
	}
	
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	
	return result
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return "127.0.0.1"
}
