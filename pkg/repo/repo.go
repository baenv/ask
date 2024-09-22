package repo

import (
	"ask/pkg/repo/server"
	serverconfig "ask/pkg/repo/server_config"
	"ask/pkg/repo/user"
	userconfig "ask/pkg/repo/user_config"
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {

	// Add methods for other entities
	User() user.IUser
	Server() server.IServer
	ServerConfig() serverconfig.IServerConfig
	UserConfig() userconfig.IUserConfig
	WithTx(fn func(txRepo Repository) error) error
}

type repository struct {
	db           *gorm.DB
	user         user.IUser
	server       server.IServer
	serverConfig serverconfig.IServerConfig
	userConfig   userconfig.IUserConfig
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db:           db,
		user:         user.New(db),
		server:       server.New(db),
		serverConfig: serverconfig.New(db),
		userConfig:   userconfig.New(db),
	}
}

func (r *repository) WithTx(fn func(txRepo Repository) error) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	txRepo := &repository{
		db:           tx,
		user:         user.New(tx),
		server:       server.New(tx),
		serverConfig: serverconfig.New(tx),
		userConfig:   userconfig.New(tx),
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after rollback
		}
	}()

	if err := fn(txRepo); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *repository) User() user.IUser {
	return r.user
}

func (r *repository) Server() server.IServer {
	return r.server
}

func (r *repository) ServerConfig() serverconfig.IServerConfig {
	return r.serverConfig
}

func (r *repository) UserConfig() userconfig.IUserConfig {
	return r.userConfig
}
