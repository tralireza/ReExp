package reexp

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type RNewPortfolio struct {
	Amount int `json:"amount"`
	Client int
	Fund   int    `json:"fund"`
	Name   string `json:"name"`
}

func AppEPoints(m *http.ServeMux) http.Handler {
	c := make(chan int)
	go func(i int) {
		for {
			c <- i
			i++
		}
	}(0)

	appServer := ReExp{
		Ver:     "0.0.1",
		XHeader: "X-RE-EXP",
		m:       m,
		chanId:  c,
	}

	type RWr = http.ResponseWriter
	type RQ = http.Request

	m.HandleFunc("POST /client/{id}/portfolio", func(w RWr, r *RQ) {
		clientId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		bfr := make([]byte, 512)
		defer r.Body.Close()

		n, err := r.Body.Read(bfr)
		if err != io.EOF {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		rNew := RNewPortfolio{}
		json.Unmarshal(bfr[:n], &rNew)
		rNew.Client = clientId // override

		if rNew.Client == 0 || rNew.Fund == 0 || rNew.Amount <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := NewPortfolio(rNew); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})

	m.HandleFunc("GET /fund", func(w RWr, r *RQ) { lsFunds(w) })

	m.HandleFunc("GET /client", func(w RWr, r *RQ) { lsClients(w) })

	m.HandleFunc("GET /client/{id}/portfolio", func(w RWr, r *RQ) {
		clientId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(404)
			return
		}
		if err := findPortfolio(w, clientId); err != nil {
			w.WriteHeader(404)
		}
	})

	m.HandleFunc("GET /client/{id}", func(w RWr, r *RQ) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(404)
			return
		}
		if err := findClient(w, id); err != nil {
			w.WriteHeader(404)
		}
	})

	m.HandleFunc("GET /fund/{id}", func(w RWr, r *RQ) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(404)
			return
		}
		if err := findFund(w, id); err != nil {
			w.WriteHeader(404)
		}
	})

	return &appServer
}

type ReExp struct {
	Ver     string
	XHeader string
	m       *http.ServeMux
	chanId  chan int
}

func (o *ReExp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rId := <-o.chanId

	ts := time.Now()
	w.Header().Add(o.XHeader, strconv.Itoa(rId))
	o.m.ServeHTTP(w, r)

	log.Printf("⚡️ [%v] {%d} %v -> %v %v [%v]", time.Since(ts), rId, r.RemoteAddr, r.Method, r.RequestURI, 0)
}
