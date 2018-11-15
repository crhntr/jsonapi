package main

import (
	"log"
	"net/http"
	"os"

	"github.com/crhntr/jsonapi"
	"github.com/globalsign/mgo"
)

func main() {
	var mux jsonapi.Mux

	type Developer struct {
		ID   string `json:"-" bson:"_id"`
		Name string `json:"name" bson:"name"`

		ManagerID string `json:"-" bson:"manager_id"` // one-to-many recursive
		TeamID    string `json:"-" bson:"team_id"`    // one-to-many
	}

	type Team struct {
		ID   string `json:"-" bson:"_id"`
		Name string `json:"name" bson:"name"`
	}

	rootSession, err := mgo.Dial(":27017")
	if err != nil {
		log.Fatal(err)
	}

	const (
		developerResourceName = "developer"
		teamResourceName      = "team"
	)

	// Developer
	mux.HandleFetchOne(developerResourceName, func(res jsonapi.FetchOneResonder, req *http.Request, idStr string) {
		sess := rootSession.Clone()
		defer sess.Close()

		// Fetch Resources from DB
		var developer Developer
		if err := sess.DB("").C(developerResourceName).FindId(idStr).One(&developer); err != nil {
			res.AppendError(err)
			return
		}

		var manager Developer
		if err := sess.DB("").C(developerResourceName).FindId(developer.ManagerID).One(&manager); err != nil {
			res.AppendError(err)
			return
		}

		var team Team
		if err := sess.DB("").C(teamResourceName).FindId(developer.TeamID).One(&team); err != nil {
			res.AppendError(err)
			return
		}

		// Set Top Level Document
		var relationships jsonapi.Relationships
		relationships.SetToOne("manager", developerResourceName, developer.ManagerID, nil)
		relationships.SetToOne("team", teamResourceName, developer.TeamID, nil)

		if err := res.Include(developerResourceName, manager.ID, manager, nil, nil); err != nil {
			res.AppendError(err)
			return
		}
		if err := res.Include(teamResourceName, team.ID, team, nil, nil); err != nil {
			res.AppendError(err)
			return
		}

		if err := res.SetData(developer.ID, developer, relationships, nil, nil); err != nil {
			res.AppendError(err)
			return
		}
	})

	mux.HandleFetchMany(developerResourceName, func(res jsonapi.FetchManyResponder, req *http.Request) {
		sess := rootSession.Clone()
		defer sess.Close()

		// var developers []Developer
	})

	mux.HandleCreate(developerResourceName, func(res jsonapi.CreateResponder, req *http.Request) {
		// ...
	})

	mux.HandleUpdate(developerResourceName, func(res jsonapi.UpdateResponder, req *http.Request, idStr string) {
		// ...
	})

	mux.HandleDelete(developerResourceName, func(res jsonapi.DeleteResponder, req *http.Request, idStr string) {
		// ...
	})

	http.ListenAndServe(":"+os.Getenv("PORT"), mux)
}
