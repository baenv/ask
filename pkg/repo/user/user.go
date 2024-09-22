package user

import "gorm.io/gorm"

type user struct {
	db *gorm.DB
}

func New(db *gorm.DB) IUser {
	return &user{db: db}
}
