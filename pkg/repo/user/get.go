package user

import "sum/pkg/models"

func (u *user) GetByPlatformID(userID, platform string) (models.User, error) {
	var user models.User
	return user, u.db.Where("user_id = ? AND platform = ?", userID, platform).First(&user).Error
}
