package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/tralireza/ReExp/reexp"
)

func init() {
	log.SetFlags(log.Ldate + log.Ltime)
	log.SetPrefix("")
}

func main() {
	var port, dbport int
	var isPg, randClients, randFunds bool
	flag.IntVar(&port, "port", 8080, "HTTP Port")
	flag.IntVar(&dbport, "dbport", 0, "DB Port, 0 (auto) -> {MySQL: 3306, Postgres: 5432}")
	flag.BoolVar(&isPg, "pg", false, "Use Postgres as DB Server")
	flag.BoolVar(&randClients, "randClients", false, "Generate 20 random client records")
	flag.BoolVar(&randFunds, "randFunds", false, "Generate 10 random fund records")
	flag.Parse()

	var db *sql.DB
	var err error
	if isPg {
		if dbport == 0 {
			dbport = 5432
		}
		db, err = sql.Open("postgres", fmt.Sprintf("postgres://reexp:pg@127.0.0.1:%d/reexp?sslmode=disable", dbport))
	} else {
		if dbport == 0 {
			dbport = 3306
		}
		params := map[string]string{}
		params["parseTime"] = "true"
		myCfg := mysql.Config{Params: params, DBName: "reexp", User: "reexp", Passwd: "mysql", Addr: "127.0.0.1:" + strconv.Itoa(dbport)}
		db, err = sql.Open("mysql", myCfg.FormatDSN())
	}
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(3 * time.Minute)
	if err := db.Ping(); err != nil {
		panic(err)
	}

	reexp.InitDBConn(db)
	log.Print("👍 connected to db server successfully.")

	if randClients {
		log.Print("📀 generating random Clients ...")
		reexp.GRandomClients(20)
		log.Print("👍 done.")
	}
	if randFunds {
		log.Print("📀 generating random Funds ...")
		reexp.GRandomFunds(10)
		log.Print("👍 done.")
	}

	h := http.DefaultServeMux
	h.HandleFunc("GET /alive", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "{%q: 0}", "status") })

	log.Printf("📀 starting on %d ...", port)
	if err := http.ListenAndServe(":"+strconv.Itoa(port), reexp.AppEPoints(h)); err != nil {
		panic(err)
	}
}
