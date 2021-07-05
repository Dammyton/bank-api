package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_application_deposit(t *testing.T) {
	r, err := http.NewRequest("GET", "/deposit?number=1234567890&amount=50", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	app, db, err := setup()
	if err != nil {
		t.Fatal(err)
	}

	// close the db connection
	defer db.Close()

	handler := http.HandlerFunc(app.deposit)
	handler.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatal("request failed :", w.Body.String())
	}
}
