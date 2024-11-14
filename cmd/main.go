package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sum/pkg/config"
	"sum/pkg/listener"
	"sum/pkg/logger"

	"github.com/bwmarrin/discordgo"
	"github.com/go-telegram/bot"
)

func main() {
	cfg := config.LoadConfig(config.DefaultConfigLoaders())
	log := logger.NewLogrusLogger()

	// TODO: disable for now
	// if err := models.InitIDGenerators(); err != nil {
	// 	log.Error(err, "Failed to initialize ID generators")
	// 	return
	// }

	// TODO: disable for now
	// dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
	// 	cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)

	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
	// 	PrepareStmt: true,
	// 	NamingStrategy: schema.NamingStrategy{
	// 		SingularTable: false,
	// 	},
	// })
	// if err != nil {
	// 	log.Error(err, "Failed to connect to database")
	// 	return
	// }

	var (
		discordSession  *discordgo.Session
		telegramSession *bot.Bot
		err             error
	)

	if config.IsDiscordEnabled(cfg) {
		discordSession, err = discordgo.New("Bot " + cfg.DiscordBotToken)
		if err != nil {
			log.Error(err, "Failed to create Discord session")
			return
		}
	}

	if config.IsTelegramEnabled(cfg) {
		telegramSession, err = bot.New(cfg.TelegramBotToken)
		if err != nil {
			log.Error(err, "Failed to create Telegram session")
			return
		}
	}

	listener := listener.New(cfg, log, discordSession, telegramSession, nil)

	if config.IsDiscordEnabled(cfg) {
		if err := listener.Discord.Start(); err != nil {
			log.Error(err, "Failed to start Discord listener")
			return
		}
		listener.Discord.Register()
		defer listener.Discord.End()
	}

	if config.IsTelegramEnabled(cfg) {
		if err := listener.Telegram.Start(); err != nil {
			log.Error(err, "Failed to start Telegram listener")
			return
		}
		listener.Telegram.Register()
		defer listener.Telegram.End()
	}

	// Start server
	srv := startServer(log)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit
	log.Info("Shutting down server...")

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	if err := srv.Shutdown(ctx); err != nil {
		log.Error(err, "Server forced to shutdown")
	}

	log.Info("Server exiting")
}

func startServer(log logger.Logger) *http.Server {
	// Add healthz endpoint
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Create a new HTTP server
	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: nil, // Use the default ServeMux
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Info("Starting HTTP server on 0.0.0.0:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error(err, "Failed to start HTTP server")
		}
	}()

	return srv
}
