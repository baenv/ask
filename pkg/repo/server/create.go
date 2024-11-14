package server

import (
	"sum/pkg/models"

	"gorm.io/gorm/clause"
)

func (s *server) Create(server models.Server) (models.Server, error) {
	return server, s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "server_id"}, {Name: "platform"}},
		DoNothing: true,
	}).Create(&server).Error
}
