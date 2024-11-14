package serverconfig

import "sum/pkg/models"

func (c serverConfig) SaveAPIKey(id string, apiKey string) error {
	return c.db.Model(&models.ServerAdminConfig{}).Where("id = ?", id).Update("api_key", apiKey).Error
}

func (c serverConfig) SaveEndpointURL(id string, endpointURL string) error {
	return c.db.Model(&models.ServerAdminConfig{}).Where("id = ?", id).Update("endpoint_url", endpointURL).Error
}

func (c serverConfig) SaveCommand(id string, command string) error {
	return c.db.Model(&models.ServerAdminConfig{}).Where("id = ?", id).Update("command", command).Error
}

func (c serverConfig) SaveDescription(id string, description string) error {
	return c.db.Model(&models.ServerAdminConfig{}).Where("id = ?", id).Update("description", description).Error
}
