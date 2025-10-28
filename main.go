package main

import (
	"gameapp/config"
	"gameapp/delivery/httpserver"
	"gameapp/repository/migrator"
	"gameapp/repository/mysql"
	"gameapp/service/authservice"
	"gameapp/service/userservice"
	"time"
)

func main() {
	cfg := config.Config{
		HTTPServer: config.HTTPServer{
			Port: 8080,
		},
		Auth: authservice.Config{
			SignKey:               "go123",
			AccessExpirationTime:  time.Hour * 24,
			RefreshExpirationTime: time.Hour * 24 * 7,
			AccessSubject:         "at",
			RefreshSubject:        "rt",
		},
		DB: mysql.Config{
			Username: "gameapp",
			Password: "gameappt0lk2o20",
			Port:     3308,
			Host:     "localhost",
			DBName:   "gameapp_db",
		},
	}

	mgr := migrator.New(cfg.DB)
	mgr.Up()

	deps := setupDependencies(cfg)

	server := httpserver.New(cfg, deps.authSvc, deps.userSvc)

	server.Serve()
}

type dependencies struct {
	authSvc authservice.Service
	userSvc userservice.Service
}

func setupDependencies(cfg config.Config) dependencies {
	dbRepo := mysql.New(cfg.DB)
	authSvc := authservice.New(cfg.Auth)
	userSvc := userservice.New(dbRepo, authSvc)

	return dependencies{
		authSvc: *authSvc,
		userSvc: *userSvc,
	}
}
