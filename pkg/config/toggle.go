package config

func IsDiscordEnabled(cfg Config) bool {
	return cfg.DiscordEnabled
}

func IsTelegramEnabled(cfg Config) bool {
	return cfg.TelegramEnabled
}
