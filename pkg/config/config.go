package config

import (
	"github.com/spf13/viper"
)

// Loader load config from reader into Viper
type Loader interface {
	Load(viper.Viper) (*viper.Viper, error)
}

// DBConfig holds the database configuration values
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// Config holds the configuration values for the application
type Config struct {
	DiscordBotToken  string   // Token for Discord bot
	TelegramBotToken string   // Token for Telegram bot
	DB               DBConfig // Database configuration
	DiscordEnabled   bool     // Flag to enable/disable Discord bot
	TelegramEnabled  bool     // Flag to enable/disable Telegram bot
	EncryptionKey    string   // Key for encryption/decryption operations
	AgentURL         string   // URL for the agent
	AgentToken       string   // Token for the agent
}

// ENV interface for environment variable retrieval
type ENV interface {
	GetString(string) string
	GetBool(string) bool
}

// Generate creates a Config struct from environment variables
func Generate(v ENV) Config {
	return Config{
		DiscordBotToken:  v.GetString("DISCORD_BOT_TOKEN"),
		TelegramBotToken: v.GetString("TELEGRAM_BOT_TOKEN"),
		DB: DBConfig{
			Host:     v.GetString("DB_HOST"),
			Port:     v.GetString("DB_PORT"),
			User:     v.GetString("DB_USER"),
			Password: v.GetString("DB_PASS"),
			Name:     v.GetString("DB_NAME"),
		},
		DiscordEnabled:  v.GetBool("DISCORD_ENABLED"),
		TelegramEnabled: v.GetBool("TELEGRAM_ENABLED"),
		EncryptionKey:   v.GetString("ENCRYPTION_KEY"),
		AgentURL:        v.GetString("AGENT_URL"),
		AgentToken:      v.GetString("AGENT_TOKEN"),
	}
}

// DefaultConfigLoaders returns a slice of default config loaders
func DefaultConfigLoaders() []Loader {
	loaders := []Loader{}
	fileLoader := NewFileLoader(".env", ".")
	loaders = append(loaders, fileLoader)
	loaders = append(loaders, NewENVLoader())

	return loaders
}

// LoadConfig loads configuration from a list of loaders
func LoadConfig(loaders []Loader) Config {
	v := viper.New()

	for idx := range loaders {
		newV, err := loaders[idx].Load(*v)

		if err == nil {
			v = newV
		}
	}
	return Generate(v)
}

// LoadTestConfig returns a Config with test values
func LoadTestConfig() Config {
	return Config{
		DiscordBotToken:  "test_discord_bot_token",
		TelegramBotToken: "test_telegram_bot_token",
		DB: DBConfig{
			Host:     "test_db_host",
			Port:     "test_db_port",
			User:     "test_db_user",
			Password: "test_db_password",
			Name:     "test_db_name",
		},
		DiscordEnabled:  false,
		TelegramEnabled: true,
		EncryptionKey:   "test_encryption_key",
		AgentURL:        "test_agent_url",
		AgentToken:      "test_agent_token",
	}
}
