package main

import (
	"encoding/json"
	"fmt"
	"gameapp/repository/mysql"
	"gameapp/service/userservice"
	"io"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/users/register", func(writer http.ResponseWriter, r *http.Request) {
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
		_, err = userservice.New(repo).Register(req)
		if err != nil {
			fmt.Fprintln(writer, err)
			return
		}
		fmt.Fprintln(writer, "user created")
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, "everything is good!")
	})

	http.ListenAndServe("localhost:8080", mux)
}
