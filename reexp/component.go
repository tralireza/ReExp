package reexp

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"
)

var (
	db *sql.DB
)

func InitDBConn(o *sql.DB) {
	db = o
}

type Client struct {
	Id   int       `json:"id"`
	DOB  time.Time `json:"dob"`
	Name string    `json:"name"`
	NI   string    `json:"ni"`
}

type Fund struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Sector string `json:"sector"`
	Type   string `json:"type"`
}

type Portfolio struct {
	Id       int       `json:"id"`
	Amount   int       `json:"amount"`
	Name     string    `json:"name"`
	OpenedAt time.Time `json:"opened_at"`
	State    int       `json:"state"`
	Fund     int       `json:"fund,omitempty"`
}

type CFP struct {
	Client    int `json:"client"`
	Fund      int `json:"fund"`
	Portfolio int `json:"portfolio"`
}

func queryClient(id int) (*Client, error) {
	var o Client
	row := db.QueryRow("SELECT * FROM client WHERE id=?", id)
	if err := row.Scan(&o.Id, &o.DOB, &o.Name, &o.NI); err != nil {
		return nil, err
	}
	return &o, nil
}

func queryFund(id int) (*Fund, error) {
	var o Fund
	row := db.QueryRow("SELECT * FROM fund WHERE id=?", id)
	if err := row.Scan(&o.Id, &o.Name, &o.Sector, &o.Type); err != nil {
		return nil, err
	}
	return &o, nil
}

func NewPortfolio(rNew RNewPortfolio) error {
	client, err := queryClient(rNew.Client)
	if err != nil {
		return err
	}
	fund, err := queryFund(rNew.Fund)
	if err != nil {
		return err
	}

	var o CFP
	row := db.QueryRow("SELECT * FROM cfp WHERE client=?", client.Id)
	if err := row.Scan(&o.Client, &o.Fund, &o.Portfolio); err != sql.ErrNoRows {
		if err != nil {
			return err
		}
		return errors.New("only one Portfolio is permitted for Client")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if rNew.Name == "" {
		rNew.Name = client.Name + " :: " + fund.Name
	}
	r, err := db.Exec("INSERT INTO portfolio(amount,state,name) VALUES(?,0,?)", rNew.Amount, rNew.Name)
	if err != nil {
		tx.Rollback()
		return err
	}
	portfolioId, err := r.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}
	if _, err := db.Exec("INSERT INTO cfp(client,fund,portfolio) VALUES(?,?,?)", client.Id, fund.Id, portfolioId); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func GRandomClients(N int) {
	t := []string{"Mr.", "Miss.", "Mrs.", "Dr.", "Sir", "Dame", "Prof."}
	f := []string{"Doe", "Smith", "Lee", "Kubrik", "Jupiter", "Spark", "Discovery"}

	rndR := func() string { return string('A' + byte(rand.Intn(26))) }
	rndN := func() int { return rand.Intn(10) }

	for range N {
		c := Client{
			Name: t[rand.Intn(len(t))] + " " + rndR() + ". " + f[rand.Intn(len(f))],
			DOB:  time.Now(),
			NI:   fmt.Sprintf("%s%s%d%d%d%d%d%d%s", rndR(), rndR(), rndN(), rndN(), rndN(), rndN(), rndN(), rndN(), rndR()),
		}
		r, _ := db.Exec("INSERT INTO client(name,dob,ni) VALUES(?,?,?)", c.Name, c.DOB, c.NI)
		Id, _ := r.LastInsertId()
		c.Id = int(Id)
		log.Printf("👍 -> %+v", c)
	}
}

func GRandomFunds(N int) {
	n := []string{"General", "Dynamic", "Stable"}
	t := []string{"Sustainable", "General"}
	s := []string{"Technology", "Stocks", "Government", "Construction", "Health"}

	for range N {
		f := Fund{
			Name:   n[rand.Intn(len(n))],
			Type:   t[rand.Intn(len(t))],
			Sector: s[rand.Intn(len(s))],
		}
		r, _ := db.Exec("INSERT INTO fund(name,sector,type) VALUES(?,?,?)", f.Name, f.Sector, f.Type)
		Id, _ := r.LastInsertId()
		f.Id = int(Id)
		log.Printf("👍 -> %+v", f)
	}
}

func lsClients(w io.Writer) error {
	rows, err := db.Query("SELECT * FROM client") // LIMIT?
	if err != nil {
		return err
	}
	defer rows.Close()

	w.Write([]byte{'['})
	i := 0
	for rows.Next() {
		if i > 0 {
			w.Write([]byte{','})
		}
		var o Client
		if err := rows.Scan(&o.Id, &o.DOB, &o.Name, &o.NI); err != nil {
			return err
		}
		bs, _ := json.MarshalIndent(&o, "", " ")
		w.Write(bs)
		i++
	}
	if err := rows.Err(); err != nil {
		return err
	}
	w.Write([]byte{']'})
	return nil
}

func writeJson(w io.Writer) *json.Encoder {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	return encoder
}

func lsFunds(w io.Writer) error {
	rows, _ := db.Query("SELECT id,name,sector,type FROM fund") // LIMIT?
	defer rows.Close()

	ls := []Fund{}
	for rows.Next() {
		var o Fund
		rows.Scan(&o.Id, &o.Name, &o.Sector, &o.Type)
		ls = append(ls, o)
	}

	writeJson(w).Encode(&ls)
	return nil
}

func findPortfolio(w io.Writer, clientId int) error {
	var o Portfolio
	row := db.QueryRow(`
      SELECT p.*,f.id FROM client c 
        JOIN cfp ON c.id=cfp.client 
        JOIN portfolio p ON p.id=cfp.portfolio 
        JOIN fund f ON f.id=cfp.fund 
      WHERE c.id=?`, clientId)
	if err := row.Scan(&o.Id, &o.Amount, &o.Name, &o.OpenedAt, &o.State, &o.Fund); err != nil {
		return err
	}

	writeJson(w).Encode(&o)
	return nil
}

func findClient(w io.Writer, id int) error {
	o, err := queryClient(id)
	if err != nil {
		return err
	}
	writeJson(w).Encode(o)
	return nil
}

func findFund(w io.Writer, id int) error {
	o, err := queryFund(id)
	if err != nil {
		return err
	}
	writeJson(w).Encode(o)
	return nil
}
