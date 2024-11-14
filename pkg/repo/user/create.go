package user

import (
	"sum/pkg/models"

	"gorm.io/gorm/clause"
)

func (u *user) Create(user models.User) (models.User, error) {
	return user, u.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "platform"}},
		DoNothing: true,
	}).Create(&user).Error
}
