package user

import "sum/pkg/models"

type IUser interface {
	Create(user models.User) (models.User, error)
	GetByPlatformID(userID, platform string) (models.User, error)
}
