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
		developerresourceType = "developer"
		teamresourceType      = "team"
	)

	// Developer
	mux.HandleFetchOne(developerresourceType, func(res jsonapi.FetchOneResonder, req *http.Request, idStr string) {
		sess := rootSession.Clone()
		defer sess.Close()

		// Fetch Resources from DB
		var developer Developer
		if err := sess.DB("").C(developerresourceType).FindId(idStr).One(&developer); err != nil {
			res.AppendError(err)
			return
		}

		var manager Developer
		if err := sess.DB("").C(developerresourceType).FindId(developer.ManagerID).One(&manager); err != nil {
			res.AppendError(err)
			return
		}

		var team Team
		if err := sess.DB("").C(teamresourceType).FindId(developer.TeamID).One(&team); err != nil {
			res.AppendError(err)
			return
		}

		// Set Top Level Document
		var relationships jsonapi.Relationships
		relationships.SetToOne("manager", developerresourceType, developer.ManagerID, nil)
		relationships.SetToOne("team", teamresourceType, developer.TeamID, nil)

		if err := res.Include(developerresourceType, manager.ID, manager, nil, nil); err != nil {
			res.AppendError(err)
			return
		}
		if err := res.Include(teamresourceType, team.ID, team, nil, nil); err != nil {
			res.AppendError(err)
			return
		}

		if err := res.SetData(developer.ID, developer, relationships, nil, nil); err != nil {
			res.AppendError(err)
			return
		}
	})

	mux.HandleFetchMany(developerresourceType, func(res jsonapi.FetchManyResponder, req *http.Request) {
		sess := rootSession.Clone()
		defer sess.Close()

		// var developers []Developer
	})

	mux.HandleCreate(developerresourceType, func(res jsonapi.CreateResponder, req *http.Request) {
		// ...
	})

	mux.HandleUpdate(developerresourceType, func(res jsonapi.UpdateResponder, req *http.Request, idStr string) {
		// ...
	})

	mux.HandleDelete(developerresourceType, func(res jsonapi.DeleteResponder, req *http.Request, idStr string) {
		// ...
	})

	http.ListenAndServe(":"+os.Getenv("PORT"), mux)
}
