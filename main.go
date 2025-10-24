package main

import (
	"encoding/json"
	"fmt"
	"gameapp/repository/mysql"
	"gameapp/service/authservice"
	"gameapp/service/userservice"
	"io"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/users/register", UserRegisterHandler)
	mux.HandleFunc("/users/login", UserLoginHandler)
	mux.HandleFunc("/users/profile", UserProfileHandler)

	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, "everything is good!")
	})

	http.ListenAndServe("localhost:8080", mux)
}

func UserProfileHandler(writer http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprintln(writer, "invalid method")
		return
	}

	authHeader := r.Header.Get("Authorization")
	authSvc := authservice.New(
		"go123",
		"at",
		"rt",
		time.Hour*24,
		time.Hour*24*7,
	)
	repo := mysql.New()
	userSvc := userservice.New(repo, authSvc)

	claims, err := authSvc.ParseAccessToken(authHeader)
	if err != nil {
		fmt.Fprintln(writer, err)
	}

	resp, err := userSvc.GetProfile(
		userservice.GetProfileRequest{
			UserID: claims.UserID,
		},
	)
	if err != nil {
		fmt.Fprintln(writer, err)
		return
	}

	json.NewEncoder(writer).Encode(resp)
}

func UserLoginHandler(writer http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprintln(writer, "invalid method")
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(writer, err)
		return
	}

	var req userservice.LoginRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		fmt.Fprintln(writer, err)
		return
	}

	repo := mysql.New()
	authSvc := authservice.New(
		"go123",
		"at",
		"rt",
		time.Hour*24,
		time.Hour*24*7,
	)

	_, err = userservice.New(repo, authSvc).Login(req)
	if err != nil {
		fmt.Fprintln(writer, err)
		return
	}

	fmt.Fprintln(writer, "user logged in")
}

func UserRegisterHandler(writer http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprintln(writer, "invalid method")
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(writer, err)
		return
	}

	var req userservice.RegisterRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		fmt.Fprintln(writer, err)
		return
	}

	repo := mysql.New()
	authSvc := authservice.New(
		"go123",
		"at",
		"rt",
		time.Hour*24,
		time.Hour*24*7,
	)

	resp, err := userservice.New(repo, authSvc).Register(req)
	if err != nil {
		fmt.Fprintln(writer, err)
		return
	}

	json.NewEncoder(writer).Encode(resp)
}
