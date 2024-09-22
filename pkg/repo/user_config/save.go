package userconfig

import (
	"ask/pkg/models"
)

func (c userConfig) SaveAPIKey(id string, apiKey string) error {
	return c.db.Model(&models.UserAgentConfig{}).Where("id = ?", id).Update("api_key", apiKey).Error
}

func (c userConfig) SaveEndpointURL(id string, endpointURL string) error {
	return c.db.Model(&models.UserAgentConfig{}).Where("id = ?", id).Update("endpoint_url", endpointURL).Error
}

func (c userConfig) SaveCommand(id string, command string) error {
	return c.db.Model(&models.UserAgentConfig{}).Where("id = ?", id).Update("command", command).Error
}

func (c userConfig) SaveDescription(id string, description string) error {
	return c.db.Model(&models.UserAgentConfig{}).Where("id = ?", id).Update("description", description).Error
}
