package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/unrolled/secure"
)

func listenAddress() string {
	if port := os.Getenv("PORT"); port != "" {
		return "0.0.0.0:" + port
	}

	return "127.0.0.1:8080"
}

func main() {
	rand.Seed(time.Now().UnixNano())

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/help", HelpHandler).Methods("GET")
	r.HandleFunc("/retrieve", RetrieveHandler).Methods("POST")

	srv := &http.Server{
		Handler: r,
		Addr:    listenAddress(),

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	r.Use(secure.New(secure.Options{
		BrowserXssFilter:     true,
		ContentTypeNosniff:   true,
		FrameDeny:            true,
		STSSeconds:           31536000,
		STSIncludeSubdomains: true,
		STSPreload:           true,
	}).Handler)

	log.Fatal(srv.ListenAndServe())
}
