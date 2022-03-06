package main

import (
	"net/http"
	"os"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	t := Parse("home.html")
	t.Execute(w, struct {
		Error      bool
		CommitHash string
	}{
		r.URL.Query().Get("error") == "t",
		os.Getenv("HEROKU_SLUG_COMMIT"),
	})
}
