package reexp

import (
	"database/sql"
	"log"
	"testing"

	"github.com/lib/pq"
)

func init() {
	dsn, err := pq.ParseURL("postgres://reexp:pg@127.0.0.1:5432/reexp?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	db, _ = sql.Open("postgres", dsn)
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func TestLoadObjects(t *testing.T) {
	rows, err := db.Query("SELECT * FROM client WHERE id>=$1 AND id<=$2", 1, 5)
	if err != nil {
		t.Fatal(err)
	}
	cls, _ := rows.Columns()
	tps, _ := rows.ColumnTypes()
	for i := range cls {
		log.Printf("%v -> %v %v %v", cls[i], tps[i].DatabaseTypeName(), tps[i].Name(), tps[i].ScanType())
	}

	os, err := LoadObjects[Client](rows)
	if err != nil {
		t.Fatal(err)
	}
	for _, o := range os {
		log.Printf("%T -> %+[1]v", o)
	}
}

func TestLoadObject(t *testing.T) {
	row := db.QueryRow("SELECT * FROM client WHERE id=$1", 1)
	client, err := LoadObject[Client](row)
	log.Print(client, err)
}
