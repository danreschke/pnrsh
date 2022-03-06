package main

import (
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	t := Parse("home.html")
	t.Execute(w, struct {
		Error bool
	}{
		r.URL.Query().Get("error") == "t",
	})
}
