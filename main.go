package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"log"

	"com/tarekkma/egyptmetro/station"

	"github.com/go-chi/render"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	dialogflow "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"

	"github.com/jmoiron/sqlx"
)

func main() {
	log.Println("Server Started")

	db, err := sqlx.Connect("sqlite3", "file:metro.db?cache=shared")
	if err != nil {
		log.Fatalln(err)
	}
	station.LoadDateFromDB(db)

	r := mux.NewRouter()

	r.HandleFunc("/from-to/{from}/{to}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		from := vars["from"]
		to := vars["to"]
		fmt.Println(from)
		fmt.Println(to)
		fromId := station.GetStationIdByName(from)
		toId := station.GetStationIdByName(to)
		if fromId == -1 || toId == -1 {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "NOT FOUND")
		}
		render.JSON(w, r, station.GoFromTo(fromId, toId))
	})

	r.HandleFunc("/stations", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, station.GetStations())
	})

	r.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		wr := dialogflow.WebhookRequest{}
		if err = json.NewDecoder(r.Body).Decode(&wr); err != nil {
			logrus.WithError(err).Error("Couldn't Unmarshal request to jsonpb")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "BAD REQUEST")
		}
		fmt.Println(wr.GetQueryResult().Action)
	})

	log.Fatal("Error happend while starting http server", http.ListenAndServe(":80", r))
}
