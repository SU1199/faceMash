package main

import (
	"facemash/algo"
	"facemash/db"
	"facemash/models"
	"fmt"
	"html/template"
	"net/http"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("html/index.html")
	if err != nil {
		fmt.Println(err)
	}

	data := db.FetchRandoms()

	t.Execute(w, data)
}

func serveVote(w http.ResponseWriter, r *http.Request) {
	winner := r.Header.Get("w")
	loser := r.Header.Get("l")
	db.UpdateRanking(winner, loser)
}

func serveRanking(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("html/rankings.html")
	if err != nil {
		fmt.Println(err)
	}

	data := new(models.RankingData)

	data.Sources = db.GetRankings()

	t.Execute(w, data)
}

func main() {
	algo.SetParams(10, 400, 1)
	db.ConnectDB()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/vote", serveVote)
	http.HandleFunc("/ranking", serveRanking)
	http.ListenAndServe(":8080", nil)
}
