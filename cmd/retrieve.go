package main

import (
	"net/http"

	"github.com/pnrsh/pnrsh/pkg/delta/pnr"
)

func RetrieveHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.Header().Add("Location", "/?error=t")
		w.WriteHeader(302)
		return
	}

	firstName := r.Form.Get("first_name")
	lastName := r.Form.Get("last_name")
	confirmationCode := r.Form.Get("confirmation_code")

	if len(confirmationCode) != 6 || len(firstName) == 0 || len(lastName) == 0 {
		w.Header().Add("Location", "/?error=t")
		w.WriteHeader(302)
		return
	}

	retrievedPNR, err := pnr.Retrieve(firstName, lastName, confirmationCode)
	if err != nil {
		w.Header().Add("Location", "/?error=t")
		w.WriteHeader(302)
		return
	}

	t := Parse("show.html")

	t.Execute(w, struct {
		PNR              pnr.PNR
		ConfirmationCode string
	}{
		retrievedPNR,
		confirmationCode,
	})
}
