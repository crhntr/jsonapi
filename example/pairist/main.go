package main

import (
	"net/http"
	"os"

	"github.com/crhntr/jsonapi"
)

func main() {
	var mux jsonapi.Mux

	// Develoepr is a strong entity
	type Developer struct {
		ID   string `json:"-"`
		Name string `json:"name"`

		ManagerID string `json:"-"` // one-to-many recursive
		TeamID    string `json:"-"` // one-to-many
	}

	// Developer
	mux.HandleFetchOne("developer", func(res jsonpai.FetchOneResonder, req *http.Request, idStr string) {
		// ...
	})

	mux.HandleFetchMany("developer", func(res jsonpai.FetchManyResponder, req *http.Request) {
		// ...
	})

	mux.HandleCreate("developer", func(res jsonpai.CreateResponder, req *http.Request) {
		// ...
	})

	mux.HandleUpdate("developer", func(res jsonpai.UpdateResponder, req *http.Request, idStr string) {
		// ...
	})

	mux.HandleDelete("developer", func(res jsonpai.DeleteResponder, req *http.Request, idStr string) {
		// ...
	})

	http.ListenAndServe(":"+os.Getenv("PORT"), mux)
}
