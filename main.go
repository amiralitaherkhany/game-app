package main

import (
	"gameapp/config"
	"gameapp/delivery/httpserver"
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

//func UserProfileHandler(writer http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodGet {
//		fmt.Fprintln(writer, "invalid method")
//		return
//	}
//
//	authHeader := r.Header.Get("Authorization")
//
//	claims, err := s.authSvc.ParseAccessToken(authHeader)
//	if err != nil {
//		fmt.Fprintln(writer, err)
//	}
//
//	resp, err := userSvc.GetProfile(
//		userservice.GetProfileRequest{
//			UserID: claims.UserID,
//		},
//	)
//	if err != nil {
//		fmt.Fprintln(writer, err)
//		return
//	}
//
//	json.NewEncoder(writer).Encode(resp)
//}
//
//func UserLoginHandler(writer http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodPost {
//		fmt.Fprintln(writer, "invalid method")
//		return
//	}
//
//	data, err := io.ReadAll(r.Body)
//	if err != nil {
//		fmt.Fprintln(writer, err)
//		return
//	}
//
//	var req userservice.LoginRequest
//	err = json.Unmarshal(data, &req)
//	if err != nil {
//		fmt.Fprintln(writer, err)
//		return
//	}
//
//	repo := mysql.New()
//	authSvc := authservice.New(
//		"go123",
//		"at",
//		"rt",
//		time.Hour*24,
//		time.Hour*24*7,
//	)
//
//	_, err = userservice.New(repo, authSvc).Login(req)
//	if err != nil {
//		fmt.Fprintln(writer, err)
//		return
//	}
//
//	fmt.Fprintln(writer, "user logged in")
//}
