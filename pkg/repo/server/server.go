package server

import "gorm.io/gorm"

type server struct {
	db *gorm.DB
}

func New(db *gorm.DB) IServer {
	return &server{db: db}
}
