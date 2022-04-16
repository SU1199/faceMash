package db

import (
	"database/sql"
	"facemash/algo"
	"facemash/models"
	"fmt"
	"html/template"
	"log"
	"strconv"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "go"
	password = "aezakmi"
	dbname   = "nalanda"
)

//global db pointer
var DBCon *sql.DB

func ConnectDB() {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}

	// check db
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	DBCon = db

	log.Println("database connected...")
}

func FetchRandoms() models.HomeData {

	homeRes := new(models.HomeData)

	type data struct {
		enum int
		url  string
	}

	dataFetch := "SELECT enum,profile FROM elo ORDER BY random() LIMIT 2;"
	rows, err := DBCon.Query(dataFetch)
	var res []data
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var s data
		var e int
		var p string
		err := rows.Scan(&e, &p)
		if err != nil {
			log.Println(err)
		}
		s.enum = e
		s.url = p
		res = append(res, s)
	}

	homeRes.EnumOne = res[0].enum
	homeRes.EnumTwo = res[1].enum
	homeRes.SrcOne = template.URL(res[0].url)
	homeRes.SrcTwo = template.URL(res[1].url)

	return *homeRes
}

func UpdateRanking(w, l string) error {
	winnerEnum, err := strconv.Atoi(w)
	if err != nil {
		log.Println(err)
		return err
	}
	loserEnum, err := strconv.Atoi(l)
	if err != nil {
		log.Println(err)
		return err
	}

	readWRating := "SELECT rating FROM elo WHERE enum=$1;"
	var winnerRating float64
	err = DBCon.QueryRow(readWRating, winnerEnum).Scan(&winnerRating)
	if err != nil {
		log.Println(err)
		return err
	}

	readLRating := "SELECT rating FROM elo WHERE enum=$1;"
	var loserRating float64
	err = DBCon.QueryRow(readLRating, loserEnum).Scan(&loserRating)
	if err != nil {
		log.Println(err)
		return err
	}

	wNew, lNew := algo.NewRatings(winnerRating, loserRating)

	log.Println(winnerEnum, wNew, loserEnum, lNew)

	update := "UPDATE elo SET rating=$1 WHERE enum=$2;"
	updateStmt, err := DBCon.Prepare(update)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = updateStmt.Exec(wNew, winnerEnum)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = updateStmt.Exec(lNew, loserEnum)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func GetRankings() []template.URL {
	rankingsStmt := "SELECT profile FROM elo ORDER BY rating DESC LIMIT 15;"
	rows, err := DBCon.Query(rankingsStmt)
	var res []template.URL
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var url string
		err := rows.Scan(&url)
		if err != nil {
			log.Println(err)
		}
		res = append(res, template.URL(url))
	}
	return res
}
