package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/tralireza/ReExp/reexp"
)

func init() {
	log.SetFlags(log.Ldate + log.Ltime)
	log.SetPrefix("")
}

func main() {
	var port, dbport int
	flag.IntVar(&port, "port", 8080, "HTTP Port")
	flag.IntVar(&dbport, "dbport", 3306, "MySQL Port")
	flag.Parse()

	dbParams := map[string]string{}
	dbParams["parseTime"] = "true"
	dbcfg := mysql.Config{
		User:   "reexp",
		Passwd: "mysql",
		Addr:   "127.0.0.1:" + strconv.Itoa(dbport),
		DBName: "reexp",
		Params: dbParams,
	}
	db, err := sql.Open("mysql", dbcfg.FormatDSN())
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(3 * time.Minute)
	if err := db.Ping(); err != nil {
		panic(err)
	}
	reexp.InitDBConn(db)
	log.Print("👍 connected to db server successfully.")

	if flag.NArg() > 0 {
		for i := range flag.NArg() {
			arg := os.Args[1+2*flag.NFlag()+i]
			if arg == "randClients" {
				log.Print("📀 generating random Clients ...")
				reexp.GRandomClients(20)
				log.Print("👍 done.")
			}
			if arg == "randFunds" {
				log.Print("📀 generating random Funds ...")
				reexp.GRandomFunds(10)
				log.Print("👍 done.")
			}
		}
	}

	h := http.DefaultServeMux
	h.HandleFunc("GET /alive", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "{%q: 0}", "status") })

	log.Printf("📀 starting on %d ...", port)
	if err := http.ListenAndServe(":"+strconv.Itoa(port), reexp.AppEPoints(h)); err != nil {
		panic(err)
	}
}
